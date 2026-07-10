package employee

import (
	"context"
	"time"

	"github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/google/uuid"
)

type Service struct {
	repo employee.Repository
}

func NewService(repo employee.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateEmployeeRequest) (*EmployeeResponse, error) {
	// Parse dates
	joinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		return nil, employee.ErrInvalidInput
	}

	// 1. Mock Auth Module Integration (Generate User ID for now)
	mockUserID := uuid.NewString()

	// 2. Instantiate Aggregate Root
	emp, err := employee.NewEmployee(mockUserID, req.EmployeeCode, req.JobPositionID, req.EmploymentStatus, joinDate)
	if err != nil {
		return nil, err
	}

	// Handle optional pointers safely
	gender := ""
	if req.PersonalData.Gender != nil {
		gender = *req.PersonalData.Gender
	}
	maritalStatus := ""
	if req.PersonalData.MaritalStatus != nil {
		maritalStatus = *req.PersonalData.MaritalStatus
	}
	ptkpStatus := ""
	if req.PersonalData.PtkpStatus != nil {
		ptkpStatus = *req.PersonalData.PtkpStatus
	}
	religion := ""
	if req.PersonalData.Religion != nil {
		religion = *req.PersonalData.Religion
	}

	// 3. Attach Personal Data
	emp.SetPersonalData(
		req.PersonalData.FullName,
		req.PersonalData.KtpNumber,
		gender,
		maritalStatus,
		ptkpStatus,
		religion,
	)

	// 4. Attach Banks & Validate Primary
	hasPrimary := false
	for _, b := range req.Banks {
		emp.AddBank(b.BankName, b.AccountNumber, b.AccountHolderName, b.IsPrimary)
		if b.IsPrimary {
			hasPrimary = true
		}
	}

	if !hasPrimary {
		return nil, employee.ErrPrimaryBankRequired
	}

	// 5. Execute in Database Transaction
	err = s.repo.ExecuteInTx(ctx, func(txCtx context.Context) error {
		// Inside this callback, txCtx contains the *gorm.DB transaction object.
		return s.repo.Save(txCtx, emp)
	})

	if err != nil {
		return nil, err
	}

	return &EmployeeResponse{
		ID:               emp.ID,
		EmployeeCode:     emp.EmployeeCode,
		Status:           emp.Status,
		EmploymentStatus: emp.EmploymentStatus,
		CreatedAt:        emp.CreatedAt,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*employee.Employee, error) {
	return s.repo.FindByID(ctx, id)
}
