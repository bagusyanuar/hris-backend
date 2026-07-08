package organization

import (
	"context"

	"github.com/bagusyanuar/hris-backend/internal/domain/organization"
	"github.com/google/uuid"
)

type Service struct {
	repo organization.Repository
}

func NewService(repo organization.Repository) *Service {
	return &Service{repo: repo}
}

// Department
func (s *Service) CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*DepartmentResponse, error) {
	id := uuid.New().String()
	dept, err := organization.NewDepartment(id, req.Code, req.Name, req.ParentID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveDepartment(ctx, dept); err != nil {
		return nil, err
	}
	return mapDepartmentToResponse(dept), nil
}

func (s *Service) GetDepartmentByID(ctx context.Context, id string) (*DepartmentResponse, error) {
	dept, err := s.repo.FindDepartmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapDepartmentToResponse(dept), nil
}

func (s *Service) GetAllDepartments(ctx context.Context) ([]*DepartmentResponse, error) {
	depts, err := s.repo.FindAllDepartments(ctx)
	if err != nil {
		return nil, err
	}
	var res []*DepartmentResponse
	for _, d := range depts {
		res = append(res, mapDepartmentToResponse(d))
	}
	return res, nil
}

func (s *Service) UpdateDepartment(ctx context.Context, id string, req UpdateDepartmentRequest) (*DepartmentResponse, error) {
	dept, err := s.repo.FindDepartmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dept.Code = req.Code
	dept.Name = req.Name
	dept.ParentID = req.ParentID
	if req.IsActive != nil {
		dept.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateDepartment(ctx, dept); err != nil {
		return nil, err
	}
	return mapDepartmentToResponse(dept), nil
}

func (s *Service) DeleteDepartment(ctx context.Context, id string) error {
	return s.repo.DeleteDepartment(ctx, id)
}

// JobTitle
func (s *Service) CreateJobTitle(ctx context.Context, req CreateJobTitleRequest) (*JobTitleResponse, error) {
	id := uuid.New().String()
	title, err := organization.NewJobTitle(id, req.Code, req.Name, req.GradeLevel)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveJobTitle(ctx, title); err != nil {
		return nil, err
	}
	return mapJobTitleToResponse(title), nil
}

func (s *Service) GetJobTitleByID(ctx context.Context, id string) (*JobTitleResponse, error) {
	title, err := s.repo.FindJobTitleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapJobTitleToResponse(title), nil
}

func (s *Service) GetAllJobTitles(ctx context.Context) ([]*JobTitleResponse, error) {
	titles, err := s.repo.FindAllJobTitles(ctx)
	if err != nil {
		return nil, err
	}
	var res []*JobTitleResponse
	for _, t := range titles {
		res = append(res, mapJobTitleToResponse(t))
	}
	return res, nil
}

func (s *Service) UpdateJobTitle(ctx context.Context, id string, req UpdateJobTitleRequest) (*JobTitleResponse, error) {
	title, err := s.repo.FindJobTitleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	title.Code = req.Code
	title.Name = req.Name
	title.GradeLevel = req.GradeLevel
	if req.IsActive != nil {
		title.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateJobTitle(ctx, title); err != nil {
		return nil, err
	}
	return mapJobTitleToResponse(title), nil
}

func (s *Service) DeleteJobTitle(ctx context.Context, id string) error {
	return s.repo.DeleteJobTitle(ctx, id)
}

// JobPosition
func (s *Service) CreateJobPosition(ctx context.Context, req CreateJobPositionRequest) (*JobPositionResponse, error) {
	id := uuid.New().String()
	pos, err := organization.NewJobPosition(id, req.DepartmentID, req.JobTitleID, req.Name, req.ReportsToID, req.HeadcountQuota)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveJobPosition(ctx, pos); err != nil {
		return nil, err
	}
	return mapJobPositionToResponse(pos), nil
}

func (s *Service) GetJobPositionByID(ctx context.Context, id string) (*JobPositionResponse, error) {
	pos, err := s.repo.FindJobPositionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapJobPositionToResponse(pos), nil
}

func (s *Service) GetAllJobPositions(ctx context.Context) ([]*JobPositionResponse, error) {
	positions, err := s.repo.FindAllJobPositions(ctx)
	if err != nil {
		return nil, err
	}
	var res []*JobPositionResponse
	for _, p := range positions {
		res = append(res, mapJobPositionToResponse(p))
	}
	return res, nil
}

func (s *Service) UpdateJobPosition(ctx context.Context, id string, req UpdateJobPositionRequest) (*JobPositionResponse, error) {
	pos, err := s.repo.FindJobPositionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	pos.DepartmentID = req.DepartmentID
	pos.JobTitleID = req.JobTitleID
	pos.Name = req.Name
	pos.ReportsToID = req.ReportsToID
	pos.HeadcountQuota = req.HeadcountQuota
	if req.IsActive != nil {
		pos.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateJobPosition(ctx, pos); err != nil {
		return nil, err
	}
	return mapJobPositionToResponse(pos), nil
}

func (s *Service) DeleteJobPosition(ctx context.Context, id string) error {
	return s.repo.DeleteJobPosition(ctx, id)
}

// Mappers
func mapDepartmentToResponse(d *organization.Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:       d.ID,
		Code:     d.Code,
		Name:     d.Name,
		ParentID: d.ParentID,
		IsActive: d.IsActive,
	}
}

func mapJobTitleToResponse(t *organization.JobTitle) *JobTitleResponse {
	return &JobTitleResponse{
		ID:         t.ID,
		Code:       t.Code,
		Name:       t.Name,
		GradeLevel: t.GradeLevel,
		IsActive:   t.IsActive,
	}
}

func mapJobPositionToResponse(p *organization.JobPosition) *JobPositionResponse {
	return &JobPositionResponse{
		ID:             p.ID,
		DepartmentID:   p.DepartmentID,
		JobTitleID:     p.JobTitleID,
		Name:           p.Name,
		ReportsToID:    p.ReportsToID,
		HeadcountQuota: p.HeadcountQuota,
		IsActive:       p.IsActive,
	}
}
