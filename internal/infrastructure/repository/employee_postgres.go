package repository

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository/models"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) employee.Repository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) SaveCore(ctx context.Context, emp *employee.Employee) error {
	model := models.EmployeeFromDomain(emp)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *EmployeeRepository) FindByID(ctx context.Context, id string) (*employee.Employee, error) {
	var model models.EmployeeModel
	if err := r.db.WithContext(ctx).
		Preload("PersonalData").
		Preload("Contact").
		Preload("Banks").
		Preload("Educations").
		Preload("Documents").
		First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, employee.ErrEmployeeNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *EmployeeRepository) SavePersonalData(ctx context.Context, data *employee.PersonalData) error {
	model := models.PersonalDataFromDomain(data)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *EmployeeRepository) FindByKTP(ctx context.Context, ktpNumber string) (*employee.PersonalData, error) {
	var model models.EmployeePersonalDataModel
	if err := r.db.WithContext(ctx).First(&model, "ktp_number = ?", ktpNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not found is fine, means no duplicate
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *EmployeeRepository) SaveContact(ctx context.Context, contact *employee.Contact) error {
	model := models.ContactFromDomain(contact)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *EmployeeRepository) SaveBanks(ctx context.Context, employeeID string, banks []*employee.Bank) error {
	// Start a transaction to replace banks
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("employee_id = ?", employeeID).Delete(&models.EmployeeBankModel{}).Error; err != nil {
			return err
		}
		for _, b := range banks {
			model := models.BankFromDomain(b)
			if err := tx.Create(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EmployeeRepository) SaveEducations(ctx context.Context, employeeID string, educations []*employee.Education) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("employee_id = ?", employeeID).Delete(&models.EmployeeEducationModel{}).Error; err != nil {
			return err
		}
		for _, ed := range educations {
			model := models.EducationFromDomain(ed)
			if err := tx.Create(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EmployeeRepository) SaveDocument(ctx context.Context, doc *employee.Document) error {
	model := models.DocumentFromDomain(doc)
	return r.db.WithContext(ctx).Save(model).Error
}
