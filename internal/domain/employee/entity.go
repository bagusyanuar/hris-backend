package employee

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput        = errors.New("invalid input data")
	ErrEmployeeNotFound    = errors.New("employee not found")
	ErrKTPDuplicate        = errors.New("ktp number already exists")
	ErrPrimaryBankRequired = errors.New("at least one primary bank account is required")
)

type Employee struct {
	ID               string
	UserID           string
	EmployeeCode     string
	JobPositionID    string
	EmploymentStatus string
	JoinDate         time.Time
	ResignDate       *time.Time
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relations (Value Objects / Child Entities)
	PersonalData *EmployeePersonalData
	Banks        []*EmployeeBank
	Educations   []*EmployeeEducation
	Documents    []*EmployeeDocument
}

type EmployeePersonalData struct {
	ID            string
	EmployeeID    string
	FullName      string
	KtpNumber     string
	Gender        string
	MaritalStatus string
	PtkpStatus    string
	Religion      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EmployeeBank struct {
	ID                string
	EmployeeID        string
	BankName          string
	AccountNumber     string
	AccountHolderName string
	IsPrimary         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type EmployeeEducation struct {
	ID              string
	EmployeeID      string
	Level           string
	InstitutionName string
	Major           string
	StartYear       int
	EndYear         int
	Score           float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type EmployeeDocument struct {
	ID           string
	EmployeeID   string
	DocumentType string
	DocumentURL  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewEmployee is the factory function enforcing invariants on creation
func NewEmployee(userID, employeeCode, jobPositionID, employmentStatus string, joinDate time.Time) (*Employee, error) {
	if userID == "" || employeeCode == "" || jobPositionID == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Employee{
		ID:               uuid.NewString(),
		UserID:           userID,
		EmployeeCode:     employeeCode,
		JobPositionID:    jobPositionID,
		EmploymentStatus: employmentStatus,
		JoinDate:         joinDate,
		Status:           "ACTIVE",
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// SetPersonalData attaches personal data to the employee
func (e *Employee) SetPersonalData(fullName, ktpNumber, gender, maritalStatus, ptkpStatus, religion string) {
	now := time.Now()
	e.PersonalData = &EmployeePersonalData{
		ID:            uuid.NewString(),
		EmployeeID:    e.ID,
		FullName:      fullName,
		KtpNumber:     ktpNumber,
		Gender:        gender,
		MaritalStatus: maritalStatus,
		PtkpStatus:    ptkpStatus,
		Religion:      religion,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// AddBank attaches a bank account to the employee
func (e *Employee) AddBank(bankName, accountNumber, accountHolderName string, isPrimary bool) {
	now := time.Now()
	e.Banks = append(e.Banks, &EmployeeBank{
		ID:                uuid.NewString(),
		EmployeeID:        e.ID,
		BankName:          bankName,
		AccountNumber:     accountNumber,
		AccountHolderName: accountHolderName,
		IsPrimary:         isPrimary,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
}
