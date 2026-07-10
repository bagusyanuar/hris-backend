package repository

import (
	"context"

	"github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository/models"
	"gorm.io/gorm"
)

type txKey struct{}

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) employee.Repository {
	return &EmployeeRepository{db: db}
}

// getDB returns the transaction DB if it exists in context, otherwise returns the standard DB
func (r *EmployeeRepository) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *EmployeeRepository) ExecuteInTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func (r *EmployeeRepository) Save(ctx context.Context, emp *employee.Employee) error {
	model := models.EmployeeFromDomain(emp)
	// Full association save
	return r.getDB(ctx).Save(model).Error
}

func (r *EmployeeRepository) FindByID(ctx context.Context, id string) (*employee.Employee, error) {
	var model models.EmployeeModel
	err := r.getDB(ctx).
		Preload("PersonalData").
		Preload("Banks").
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, employee.ErrEmployeeNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *EmployeeRepository) FindAll(ctx context.Context) ([]*employee.Employee, error) {
	var dbModels []models.EmployeeModel
	err := r.getDB(ctx).
		Preload("PersonalData").
		Find(&dbModels).Error

	if err != nil {
		return nil, err
	}

	var domains []*employee.Employee
	for _, m := range dbModels {
		domains = append(domains, m.ToDomain())
	}
	return domains, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, emp *employee.Employee) error {
	model := models.EmployeeFromDomain(emp)
	// Omit relationships on standard update to avoid unintended overwrites, handle separately if needed
	return r.getDB(ctx).Omit("PersonalData", "Banks").Updates(model).Error
}

func (r *EmployeeRepository) Delete(ctx context.Context, id string) error {
	return r.getDB(ctx).Delete(&models.EmployeeModel{}, "id = ?", id).Error
}
