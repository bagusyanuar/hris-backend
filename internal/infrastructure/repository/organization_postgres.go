package repository

import (
	"context"

	"github.com/bagusyanuar/hris-backend/internal/domain/organization"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository/models"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) organization.Repository {
	return &organizationRepository{db: db}
}

// Department
func (r *organizationRepository) SaveDepartment(ctx context.Context, dept *organization.Department) error {
	model := models.DepartmentFromDomain(dept)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *organizationRepository) FindDepartmentByID(ctx context.Context, id string) (*organization.Department, error) {
	var model models.DepartmentModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, organization.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *organizationRepository) FindAllDepartments(ctx context.Context) ([]*organization.Department, error) {
	var dbModels []models.DepartmentModel
	if err := r.db.WithContext(ctx).Find(&dbModels).Error; err != nil {
		return nil, err
	}
	var domains []*organization.Department
	for _, m := range dbModels {
		domains = append(domains, m.ToDomain())
	}
	return domains, nil
}

func (r *organizationRepository) UpdateDepartment(ctx context.Context, dept *organization.Department) error {
	model := models.DepartmentFromDomain(dept)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *organizationRepository) DeleteDepartment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.DepartmentModel{}, "id = ?", id).Error
}

// JobTitle
func (r *organizationRepository) SaveJobTitle(ctx context.Context, title *organization.JobTitle) error {
	model := models.JobTitleFromDomain(title)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *organizationRepository) FindJobTitleByID(ctx context.Context, id string) (*organization.JobTitle, error) {
	var model models.JobTitleModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, organization.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *organizationRepository) FindAllJobTitles(ctx context.Context) ([]*organization.JobTitle, error) {
	var dbModels []models.JobTitleModel
	if err := r.db.WithContext(ctx).Find(&dbModels).Error; err != nil {
		return nil, err
	}
	var domains []*organization.JobTitle
	for _, m := range dbModels {
		domains = append(domains, m.ToDomain())
	}
	return domains, nil
}

func (r *organizationRepository) UpdateJobTitle(ctx context.Context, title *organization.JobTitle) error {
	model := models.JobTitleFromDomain(title)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *organizationRepository) DeleteJobTitle(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.JobTitleModel{}, "id = ?", id).Error
}

// JobPosition
func (r *organizationRepository) SaveJobPosition(ctx context.Context, pos *organization.JobPosition) error {
	model := models.JobPositionFromDomain(pos)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *organizationRepository) FindJobPositionByID(ctx context.Context, id string) (*organization.JobPosition, error) {
	var model models.JobPositionModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, organization.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *organizationRepository) FindAllJobPositions(ctx context.Context) ([]*organization.JobPosition, error) {
	var dbModels []models.JobPositionModel
	if err := r.db.WithContext(ctx).Find(&dbModels).Error; err != nil {
		return nil, err
	}
	var domains []*organization.JobPosition
	for _, m := range dbModels {
		domains = append(domains, m.ToDomain())
	}
	return domains, nil
}

func (r *organizationRepository) UpdateJobPosition(ctx context.Context, pos *organization.JobPosition) error {
	model := models.JobPositionFromDomain(pos)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *organizationRepository) DeleteJobPosition(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.JobPositionModel{}, "id = ?", id).Error
}
