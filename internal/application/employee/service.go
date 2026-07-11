package employee

import (
	"context"
	"time"

	domain "github.com/bagusyanuar/hris-backend/internal/domain/employee"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCore(ctx context.Context, req CreateEmployeeRequest) (*CreateEmployeeResponse, error) {
	parsedJoinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	emp, err := domain.NewEmployee(req.EmployeeCode, req.JobPositionID, req.EmploymentStatus, parsedJoinDate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveCore(ctx, emp); err != nil {
		return nil, err
	}

	return &CreateEmployeeResponse{
		ID:               emp.ID,
		EmployeeCode:     emp.EmployeeCode,
		EmploymentStatus: emp.EmploymentStatus,
		Status:           emp.Status,
		CreatedAt:        emp.CreatedAt,
	}, nil
}

func (s *Service) UpdatePersonalData(ctx context.Context, employeeID string, req UpdatePersonalDataRequest) error {
	// check if employee exists
	if _, err := s.repo.FindByID(ctx, employeeID); err != nil {
		return err
	}

	// check duplicate ktp
	existingKtp, err := s.repo.FindByKTP(ctx, req.KtpNumber)
	if err != nil {
		return err
	}
	if existingKtp != nil && existingKtp.EmployeeID != employeeID {
		return domain.ErrKTPDuplicate
	}

	personalData, err := domain.NewPersonalData(
		employeeID,
		req.FullName,
		req.KtpNumber,
		req.Gender,
		req.MaritalStatus,
		req.PtkpStatus,
		req.Religion,
	)
	if err != nil {
		return err
	}

	return s.repo.SavePersonalData(ctx, personalData)
}

func (s *Service) UpdateContact(ctx context.Context, employeeID string, req UpdateContactRequest) error {
	if _, err := s.repo.FindByID(ctx, employeeID); err != nil {
		return err
	}

	contact, err := domain.NewContact(
		employeeID,
		req.PersonalEmail,
		req.PhoneNumber,
		req.IdentityAddress,
		req.ResidentialAddress,
	)
	if err != nil {
		return err
	}

	return s.repo.SaveContact(ctx, contact)
}

func (s *Service) SaveBanks(ctx context.Context, employeeID string, req SaveBanksRequest) error {
	if _, err := s.repo.FindByID(ctx, employeeID); err != nil {
		return err
	}

	var hasPrimary bool
	var domainBanks []*domain.Bank

	for _, b := range req.Banks {
		if b.IsPrimary {
			hasPrimary = true
		}
		bank, err := domain.NewBank(employeeID, b.BankName, b.AccountNumber, b.AccountHolderName, b.IsPrimary)
		if err != nil {
			return err
		}
		domainBanks = append(domainBanks, bank)
	}

	if !hasPrimary {
		return domain.ErrPrimaryBankRequired
	}

	return s.repo.SaveBanks(ctx, employeeID, domainBanks)
}

func (s *Service) GetEmployeeDetail(ctx context.Context, employeeID string) (*GetEmployeeDetailResponse, error) {
	emp, err := s.repo.FindByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	res := &GetEmployeeDetailResponse{
		ID:               emp.ID,
		EmployeeCode:     emp.EmployeeCode,
		JobPositionID:    emp.JobPositionID,
		EmploymentStatus: emp.EmploymentStatus,
		JoinDate:         emp.JoinDate.Format("2006-01-02"),
		Status:           emp.Status,
	}

	if emp.PersonalData != nil {
		res.PersonalData = &PersonalDataResponse{
			FullName:      emp.PersonalData.FullName,
			KtpNumber:     emp.PersonalData.KtpNumber,
			Gender:        emp.PersonalData.Gender,
			MaritalStatus: emp.PersonalData.MaritalStatus,
			PtkpStatus:    emp.PersonalData.PtkpStatus,
			Religion:      emp.PersonalData.Religion,
		}
	}

	if emp.Contact != nil {
		res.Contact = &ContactResponse{
			PersonalEmail:      emp.Contact.PersonalEmail,
			PhoneNumber:        emp.Contact.PhoneNumber,
			IdentityAddress:    emp.Contact.IdentityAddress,
			ResidentialAddress: emp.Contact.ResidentialAddress,
		}
	}

	if len(emp.Banks) > 0 {
		var banks []BankResponse
		for _, b := range emp.Banks {
			banks = append(banks, BankResponse{
				BankName:          b.BankName,
				AccountNumber:     b.AccountNumber,
				AccountHolderName: b.AccountHolderName,
				IsPrimary:         b.IsPrimary,
			})
		}
		res.Banks = banks
	}

	return res, nil
}
