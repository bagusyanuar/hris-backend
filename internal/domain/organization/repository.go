package organization

import "context"

type Repository interface {
	// Department
	SaveDepartment(ctx context.Context, dept *Department) error
	FindDepartmentByID(ctx context.Context, id string) (*Department, error)
	FindAllDepartments(ctx context.Context) ([]*Department, error)
	UpdateDepartment(ctx context.Context, dept *Department) error
	DeleteDepartment(ctx context.Context, id string) error

	// JobTitle
	SaveJobTitle(ctx context.Context, title *JobTitle) error
	FindJobTitleByID(ctx context.Context, id string) (*JobTitle, error)
	FindAllJobTitles(ctx context.Context) ([]*JobTitle, error)
	UpdateJobTitle(ctx context.Context, title *JobTitle) error
	DeleteJobTitle(ctx context.Context, id string) error

	// JobPosition
	SaveJobPosition(ctx context.Context, pos *JobPosition) error
	FindJobPositionByID(ctx context.Context, id string) (*JobPosition, error)
	FindAllJobPositions(ctx context.Context) ([]*JobPosition, error)
	UpdateJobPosition(ctx context.Context, pos *JobPosition) error
	DeleteJobPosition(ctx context.Context, id string) error
}
