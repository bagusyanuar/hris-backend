package application

import (
	"context"
	"errors"

	orgApplication "github.com/bagusyanuar/hris-backend/internal/organization/application"
	"github.com/bagusyanuar/hris-backend/internal/workforce/domain"
	"github.com/bagusyanuar/hris-backend/pkg/pagination"
)

// Service koordinasi transaksi & mapping DTO untuk 3 pilar Workforce Structure.
// orgService adalah Application Service modul Organization (bukan repository
// langsung) untuk validasi company_id — coding-convention.md §4.
type Service struct {
	departmentRepo  domain.DepartmentRepository
	jobTitleRepo    domain.JobTitleRepository
	jobPositionRepo domain.JobPositionRepository
	txManager       domain.TxManager
	orgService      *orgApplication.Service
}

func NewService(
	departmentRepo domain.DepartmentRepository,
	jobTitleRepo domain.JobTitleRepository,
	jobPositionRepo domain.JobPositionRepository,
	txManager domain.TxManager,
	orgService *orgApplication.Service,
) *Service {
	return &Service{
		departmentRepo:  departmentRepo,
		jobTitleRepo:    jobTitleRepo,
		jobPositionRepo: jobPositionRepo,
		txManager:       txManager,
		orgService:      orgService,
	}
}

func toDepartmentResponse(d *domain.Department) DepartmentResponse {
	return DepartmentResponse{
		ID:        d.ID,
		CompanyID: d.CompanyID,
		Code:      d.Code,
		Name:      d.Name,
		ParentID:  d.ParentID,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func toJobTitleResponse(jt *domain.JobTitle) JobTitleResponse {
	return JobTitleResponse{
		ID:         jt.ID,
		CompanyID:  jt.CompanyID,
		Code:       jt.Code,
		Name:       jt.Name,
		GradeLevel: jt.GradeLevel,
		IsActive:   jt.IsActive,
		CreatedAt:  jt.CreatedAt,
		UpdatedAt:  jt.UpdatedAt,
	}
}

func toJobPositionResponse(jp *domain.JobPosition, departmentName, jobTitleName string) JobPositionResponse {
	return JobPositionResponse{
		ID:             jp.ID,
		CompanyID:      jp.CompanyID,
		Department:     JobPositionRef{ID: jp.DepartmentID, Name: departmentName},
		JobTitle:       JobPositionRef{ID: jp.JobTitleID, Name: jobTitleName},
		Name:           jp.Name,
		ReportsToID:    jp.ReportsToID,
		HeadcountQuota: jp.HeadcountQuota,
		IsActive:       jp.IsActive,
		CreatedAt:      jp.CreatedAt,
		UpdatedAt:      jp.UpdatedAt,
	}
}

// jobPositionNamesBatch ambil department_name & job_title_name buat sekumpulan JobPosition
// sekali query batch (hindari N+1, decision-log.md ADR-006), dipakai List & Chart.
func (s *Service) jobPositionNamesBatch(ctx context.Context, items []*domain.JobPosition) ([]JobPositionResponse, error) {
	departmentIDs := make([]string, 0, len(items))
	jobTitleIDs := make([]string, 0, len(items))
	for _, jp := range items {
		departmentIDs = append(departmentIDs, jp.DepartmentID)
		jobTitleIDs = append(jobTitleIDs, jp.JobTitleID)
	}
	departmentNames, err := s.departmentRepo.FindNamesByIDs(ctx, departmentIDs)
	if err != nil {
		return nil, err
	}
	jobTitleNames, err := s.jobTitleRepo.FindNamesByIDs(ctx, jobTitleIDs)
	if err != nil {
		return nil, err
	}
	res := make([]JobPositionResponse, 0, len(items))
	for _, jp := range items {
		res = append(res, toJobPositionResponse(jp, departmentNames[jp.DepartmentID], jobTitleNames[jp.JobTitleID]))
	}
	return res, nil
}

// --- Department use cases ---

func (s *Service) CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*DepartmentResponse, error) {
	if _, err := s.orgService.GetCompany(ctx, req.CompanyID); err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		parent, err := s.departmentRepo.FindByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent.CompanyID != req.CompanyID {
			return nil, domain.ErrDepartmentCompanyMismatch
		}
	}

	department, err := domain.NewDepartment(req.CompanyID, req.Code, req.Name, req.ParentID)
	if err != nil {
		return nil, err
	}

	if _, err := s.departmentRepo.FindByCompanyAndCode(ctx, req.CompanyID, req.Code); err == nil {
		return nil, domain.ErrDepartmentCodeDuplicate
	} else if !errors.Is(err, domain.ErrDepartmentNotFound) {
		return nil, err
	}

	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return nil, err
	}
	res := toDepartmentResponse(department)
	return &res, nil
}

func (s *Service) GetDepartment(ctx context.Context, id string) (*DepartmentResponse, error) {
	department, err := s.departmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toDepartmentResponse(department)
	return &res, nil
}

func (s *Service) ListDepartments(ctx context.Context, page, limit int, sort, order, search string) (*DepartmentListResponse, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
	items, total, err := s.departmentRepo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order, search)
	if err != nil {
		return nil, err
	}
	res := make([]DepartmentResponse, 0, len(items))
	for _, d := range items {
		res = append(res, toDepartmentResponse(d))
	}
	return &DepartmentListResponse{Items: res, Meta: pagination.NewMeta(req, total)}, nil
}

func (s *Service) UpdateDepartment(ctx context.Context, id string, req UpdateDepartmentRequest) (*DepartmentResponse, error) {
	department, err := s.departmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newCode := department.Code
	if req.Code != nil {
		newCode = *req.Code
	}
	if req.Name != nil {
		department.Name = *req.Name
	}
	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}
	if newCode == "" || department.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	parentChanged := false
	newParentID := department.ParentID
	if req.ParentID != nil {
		newParentID = req.ParentID
		parentChanged = true
	}
	if parentChanged && newParentID != nil {
		parent, err := s.departmentRepo.FindByID(ctx, *newParentID)
		if err != nil {
			return nil, err
		}
		if parent.CompanyID != department.CompanyID {
			return nil, domain.ErrDepartmentCompanyMismatch
		}
		// Cycle check hanya relevan di Update — Create baru tidak mungkin cycle
		// (tech-spec.md §7.2).
		if err := domain.DetectCycle(ctx, department.ID, *newParentID, s.departmentRepo.FindParentID); err != nil {
			return nil, err
		}
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		if newCode != department.Code {
			if _, err := s.departmentRepo.FindByCompanyAndCode(ctx, department.CompanyID, newCode); err == nil {
				return domain.ErrDepartmentCodeDuplicate
			} else if !errors.Is(err, domain.ErrDepartmentNotFound) {
				return err
			}
		}
		department.Code = newCode
		department.ParentID = newParentID
		return s.departmentRepo.Update(ctx, department)
	})
	if err != nil {
		return nil, err
	}

	res := toDepartmentResponse(department)
	return &res, nil
}

// ListDepartmentsTree — TANPA pagination, dipakai FE Tabel (nested row, expand/collapse)
// & Bagan (tree diagram) render struktur Department utuh dari parent_id.
func (s *Service) ListDepartmentsTree(ctx context.Context) ([]DepartmentResponse, error) {
	items, err := s.departmentRepo.FindAllTree(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]DepartmentResponse, 0, len(items))
	for _, d := range items {
		res = append(res, toDepartmentResponse(d))
	}
	return res, nil
}

func (s *Service) DeleteDepartment(ctx context.Context, id string) error {
	if _, err := s.departmentRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.departmentRepo.Delete(ctx, id)
}

// --- Job Title use cases ---

func (s *Service) CreateJobTitle(ctx context.Context, req CreateJobTitleRequest) (*JobTitleResponse, error) {
	if _, err := s.orgService.GetCompany(ctx, req.CompanyID); err != nil {
		return nil, err
	}

	jobTitle, err := domain.NewJobTitle(req.CompanyID, req.Code, req.Name, req.GradeLevel)
	if err != nil {
		return nil, err
	}

	if _, err := s.jobTitleRepo.FindByCompanyAndCode(ctx, req.CompanyID, req.Code); err == nil {
		return nil, domain.ErrJobTitleCodeDuplicate
	} else if !errors.Is(err, domain.ErrJobTitleNotFound) {
		return nil, err
	}

	if err := s.jobTitleRepo.Create(ctx, jobTitle); err != nil {
		return nil, err
	}
	res := toJobTitleResponse(jobTitle)
	return &res, nil
}

func (s *Service) GetJobTitle(ctx context.Context, id string) (*JobTitleResponse, error) {
	jobTitle, err := s.jobTitleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toJobTitleResponse(jobTitle)
	return &res, nil
}

func (s *Service) ListJobTitles(ctx context.Context, page, limit int, sort, order, search string) (*JobTitleListResponse, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
	items, total, err := s.jobTitleRepo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order, search)
	if err != nil {
		return nil, err
	}
	res := make([]JobTitleResponse, 0, len(items))
	for _, jt := range items {
		res = append(res, toJobTitleResponse(jt))
	}
	return &JobTitleListResponse{Items: res, Meta: pagination.NewMeta(req, total)}, nil
}

func (s *Service) UpdateJobTitle(ctx context.Context, id string, req UpdateJobTitleRequest) (*JobTitleResponse, error) {
	jobTitle, err := s.jobTitleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newCode := jobTitle.Code
	if req.Code != nil {
		newCode = *req.Code
	}
	if req.Name != nil {
		jobTitle.Name = *req.Name
	}
	if req.GradeLevel != nil {
		jobTitle.GradeLevel = *req.GradeLevel
	}
	if req.IsActive != nil {
		jobTitle.IsActive = *req.IsActive
	}
	if newCode == "" || jobTitle.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		if newCode != jobTitle.Code {
			if _, err := s.jobTitleRepo.FindByCompanyAndCode(ctx, jobTitle.CompanyID, newCode); err == nil {
				return domain.ErrJobTitleCodeDuplicate
			} else if !errors.Is(err, domain.ErrJobTitleNotFound) {
				return err
			}
		}
		jobTitle.Code = newCode
		return s.jobTitleRepo.Update(ctx, jobTitle)
	})
	if err != nil {
		return nil, err
	}

	res := toJobTitleResponse(jobTitle)
	return &res, nil
}

func (s *Service) DeleteJobTitle(ctx context.Context, id string) error {
	if _, err := s.jobTitleRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.jobTitleRepo.Delete(ctx, id)
}

// --- Job Position use cases ---

func (s *Service) CreateJobPosition(ctx context.Context, req CreateJobPositionRequest) (*JobPositionResponse, error) {
	department, err := s.departmentRepo.FindByID(ctx, req.DepartmentID)
	if err != nil {
		return nil, err
	}
	jobTitle, err := s.jobTitleRepo.FindByID(ctx, req.JobTitleID)
	if err != nil {
		return nil, err
	}
	if department.CompanyID != jobTitle.CompanyID {
		return nil, domain.ErrJobPositionCompanyMismatch
	}

	if req.ReportsToID != nil {
		reportsTo, err := s.jobPositionRepo.FindByID(ctx, *req.ReportsToID)
		if err != nil {
			return nil, err
		}
		if reportsTo.CompanyID != department.CompanyID {
			return nil, domain.ErrReportingCompanyMismatch
		}
		// Cycle check hanya relevan di Update — Create baru tidak mungkin cycle
		// (tech-spec.md §7.2).
	}

	jobPosition, err := domain.NewJobPosition(department.CompanyID, req.DepartmentID, req.JobTitleID, req.Name, req.ReportsToID, req.HeadcountQuota)
	if err != nil {
		return nil, err
	}

	if err := s.jobPositionRepo.Create(ctx, jobPosition); err != nil {
		return nil, err
	}
	res := toJobPositionResponse(jobPosition, department.Name, jobTitle.Name)
	return &res, nil
}

func (s *Service) GetJobPosition(ctx context.Context, id string) (*JobPositionResponse, error) {
	jobPosition, err := s.jobPositionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	department, err := s.departmentRepo.FindByID(ctx, jobPosition.DepartmentID)
	if err != nil {
		return nil, err
	}
	jobTitle, err := s.jobTitleRepo.FindByID(ctx, jobPosition.JobTitleID)
	if err != nil {
		return nil, err
	}
	res := toJobPositionResponse(jobPosition, department.Name, jobTitle.Name)
	return &res, nil
}

func (s *Service) ListJobPositions(ctx context.Context, page, limit int, sort, order, search string) (*JobPositionListResponse, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
	items, total, err := s.jobPositionRepo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order, search)
	if err != nil {
		return nil, err
	}
	res, err := s.jobPositionNamesBatch(ctx, items)
	if err != nil {
		return nil, err
	}
	return &JobPositionListResponse{Items: res, Meta: pagination.NewMeta(req, total)}, nil
}

// ListJobPositionsChart — TANPA pagination (decision-log.md ADR-004), dipakai
// FE render Organization Chart utuh dari reports_to_id.
func (s *Service) ListJobPositionsChart(ctx context.Context) ([]JobPositionResponse, error) {
	items, err := s.jobPositionRepo.FindAllChart(ctx)
	if err != nil {
		return nil, err
	}
	return s.jobPositionNamesBatch(ctx, items)
}

func (s *Service) UpdateJobPosition(ctx context.Context, id string, req UpdateJobPositionRequest) (*JobPositionResponse, error) {
	jobPosition, err := s.jobPositionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	companyID := jobPosition.CompanyID
	departmentID := jobPosition.DepartmentID
	jobTitleID := jobPosition.JobTitleID
	structureChanged := false
	if req.DepartmentID != nil {
		departmentID = *req.DepartmentID
		structureChanged = true
	}
	if req.JobTitleID != nil {
		jobTitleID = *req.JobTitleID
		structureChanged = true
	}
	if structureChanged {
		department, err := s.departmentRepo.FindByID(ctx, departmentID)
		if err != nil {
			return nil, err
		}
		jobTitle, err := s.jobTitleRepo.FindByID(ctx, jobTitleID)
		if err != nil {
			return nil, err
		}
		if department.CompanyID != jobTitle.CompanyID {
			return nil, domain.ErrJobPositionCompanyMismatch
		}
		companyID = department.CompanyID
	}

	if req.Name != nil {
		jobPosition.Name = *req.Name
	}
	if req.HeadcountQuota != nil {
		jobPosition.HeadcountQuota = *req.HeadcountQuota
	}
	if jobPosition.HeadcountQuota < 1 {
		jobPosition.HeadcountQuota = 1
	}
	if req.IsActive != nil {
		jobPosition.IsActive = *req.IsActive
	}
	if jobPosition.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	reportsToChanged := false
	newReportsToID := jobPosition.ReportsToID
	if req.ReportsToID != nil {
		newReportsToID = req.ReportsToID
		reportsToChanged = true
	}
	if reportsToChanged && newReportsToID != nil {
		reportsTo, err := s.jobPositionRepo.FindByID(ctx, *newReportsToID)
		if err != nil {
			return nil, err
		}
		if reportsTo.CompanyID != companyID {
			return nil, domain.ErrReportingCompanyMismatch
		}
		if err := domain.DetectCycle(ctx, jobPosition.ID, *newReportsToID, s.jobPositionRepo.FindParentID); err != nil {
			return nil, err
		}
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		jobPosition.CompanyID = companyID
		jobPosition.DepartmentID = departmentID
		jobPosition.JobTitleID = jobTitleID
		jobPosition.ReportsToID = newReportsToID
		return s.jobPositionRepo.Update(ctx, jobPosition)
	})
	if err != nil {
		return nil, err
	}

	department, err := s.departmentRepo.FindByID(ctx, jobPosition.DepartmentID)
	if err != nil {
		return nil, err
	}
	jobTitle, err := s.jobTitleRepo.FindByID(ctx, jobPosition.JobTitleID)
	if err != nil {
		return nil, err
	}
	res := toJobPositionResponse(jobPosition, department.Name, jobTitle.Name)
	return &res, nil
}

func (s *Service) DeleteJobPosition(ctx context.Context, id string) error {
	if _, err := s.jobPositionRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.jobPositionRepo.Delete(ctx, id)
}
