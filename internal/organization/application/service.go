package application

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/organization/domain"
	"github.com/bagusyanuar/hris-backend/pkg/pagination"
)

type Service struct {
	companyRepo domain.CompanyRepository
	branchRepo  domain.BranchRepository
	txManager   domain.TxManager
}

func NewService(companyRepo domain.CompanyRepository, branchRepo domain.BranchRepository, txManager domain.TxManager) *Service {
	return &Service{
		companyRepo: companyRepo,
		branchRepo:  branchRepo,
		txManager:   txManager,
	}
}

func toCompanyResponse(c *domain.Company) CompanyResponse {
	return CompanyResponse{
		ID:        c.ID,
		Code:      c.Code,
		LegalName: c.LegalName,
		Npwp:      c.Npwp,
		BpjsNo:    c.BpjsNo,
		IsActive:  c.IsActive,
		Branches:  []BranchResponse{},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toBranchResponse(b *domain.Branch) BranchResponse {
	return BranchResponse{
		ID:        b.ID,
		CompanyID: b.CompanyID,
		Code:      b.Code,
		Name:      b.Name,
		City:      b.City,
		IsMain:    b.IsMain,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

// --- Company use cases ---

func (s *Service) CreateCompany(ctx context.Context, req CreateCompanyRequest) (*CompanyResponse, error) {
	company, err := domain.NewCompany(req.Code, req.LegalName, req.Npwp, req.BpjsNo)
	if err != nil {
		return nil, err
	}
	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}
	res := toCompanyResponse(company)
	return &res, nil
}

func (s *Service) GetCompany(ctx context.Context, id string) (*CompanyResponse, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toCompanyResponse(company)
	return &res, nil
}

func (s *Service) ListCompanies(ctx context.Context, page, limit int, sort, order, search string) (*CompanyListResponse, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
	companies, total, err := s.companyRepo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order, search)
	if err != nil {
		return nil, err
	}

	companyIDs := make([]string, 0, len(companies))
	for _, c := range companies {
		companyIDs = append(companyIDs, c.ID)
	}
	// Batch query, bukan N+1 (decision-log.md ADR-006) — branches nested tetap FULL LIST
	// milik company itu, tidak difilter oleh `search`.
	branches, err := s.branchRepo.FindAllByCompanyIDs(ctx, companyIDs)
	if err != nil {
		return nil, err
	}
	branchesByCompany := make(map[string][]BranchResponse, len(companyIDs))
	for _, b := range branches {
		branchesByCompany[b.CompanyID] = append(branchesByCompany[b.CompanyID], toBranchResponse(b))
	}

	items := make([]CompanyResponse, 0, len(companies))
	for _, c := range companies {
		item := toCompanyResponse(c)
		if bs, ok := branchesByCompany[c.ID]; ok {
			item.Branches = bs
		}
		items = append(items, item)
	}
	return &CompanyListResponse{
		Items: items,
		Meta:  pagination.NewMeta(req, total),
	}, nil
}

func (s *Service) UpdateCompany(ctx context.Context, id string, req UpdateCompanyRequest) (*CompanyResponse, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		company.Code = *req.Code
	}
	if req.LegalName != nil {
		company.LegalName = *req.LegalName
	}
	if req.Npwp != nil {
		company.Npwp = req.Npwp
	}
	if req.BpjsNo != nil {
		company.BpjsNo = req.BpjsNo
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}
	if company.Code == "" || company.LegalName == "" {
		return nil, domain.ErrInvalidInput
	}
	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, err
	}
	res := toCompanyResponse(company)
	return &res, nil
}

func (s *Service) DeleteCompany(ctx context.Context, id string) error {
	if _, err := s.companyRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.companyRepo.Delete(ctx, id)
}

// --- Branch use cases ---

func (s *Service) CreateBranch(ctx context.Context, companyID string, req CreateBranchRequest) (*BranchResponse, error) {
	if _, err := s.companyRepo.FindByID(ctx, companyID); err != nil {
		return nil, err
	}

	branch, err := domain.NewBranch(companyID, req.Code, req.Name, req.City, req.IsMain)
	if err != nil {
		return nil, err
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		if _, err := s.branchRepo.FindByCompanyAndCode(ctx, companyID, req.Code); err == nil {
			return domain.ErrBranchCodeDuplicate
		} else if !errors.Is(err, domain.ErrBranchNotFound) {
			return err
		}
		// Auto-demote main branch lama sebelum insert (decision-log.md ADR-004).
		if branch.IsMain {
			if err := s.branchRepo.DemoteMainBranch(ctx, companyID); err != nil {
				return err
			}
		}
		return s.branchRepo.Create(ctx, branch)
	})
	if err != nil {
		return nil, err
	}

	res := toBranchResponse(branch)
	return &res, nil
}

func (s *Service) GetBranch(ctx context.Context, id string) (*BranchResponse, error) {
	branch, err := s.branchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toBranchResponse(branch)
	return &res, nil
}

func (s *Service) ListBranchesByCompany(ctx context.Context, companyID string, page, limit int, sort, order string) (*BranchListResponse, error) {
	if _, err := s.companyRepo.FindByID(ctx, companyID); err != nil {
		return nil, err
	}
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
	branches, total, err := s.branchRepo.FindAllByCompany(ctx, companyID, req.Page, req.Limit, req.Sort, req.Order)
	if err != nil {
		return nil, err
	}
	items := make([]BranchResponse, 0, len(branches))
	for _, b := range branches {
		items = append(items, toBranchResponse(b))
	}
	return &BranchListResponse{
		Items: items,
		Meta:  pagination.NewMeta(req, total),
	}, nil
}

func (s *Service) UpdateBranch(ctx context.Context, id string, req UpdateBranchRequest) (*BranchResponse, error) {
	branch, err := s.branchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newCode := branch.Code
	if req.Code != nil {
		newCode = *req.Code
	}
	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.City != nil {
		branch.City = req.City
	}
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}
	wantMain := branch.IsMain
	if req.IsMain != nil {
		wantMain = *req.IsMain
	}
	if newCode == "" || branch.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		if newCode != branch.Code {
			if _, err := s.branchRepo.FindByCompanyAndCode(ctx, branch.CompanyID, newCode); err == nil {
				return domain.ErrBranchCodeDuplicate
			} else if !errors.Is(err, domain.ErrBranchNotFound) {
				return err
			}
		}
		// Auto-demote main branch lama sebelum promote branch ini (decision-log.md ADR-004).
		if wantMain && !branch.IsMain {
			if err := s.branchRepo.DemoteMainBranch(ctx, branch.CompanyID); err != nil {
				return err
			}
		}
		branch.Code = newCode
		branch.IsMain = wantMain
		return s.branchRepo.Update(ctx, branch)
	})
	if err != nil {
		return nil, err
	}

	res := toBranchResponse(branch)
	return &res, nil
}

func (s *Service) DeleteBranch(ctx context.Context, id string) error {
	if _, err := s.branchRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.branchRepo.Delete(ctx, id)
}
