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
	
	// Relations for Progressive Save / Detail
	PersonalData *PersonalData
	Contact      *Contact
	Banks        []*Bank
	Educations   []*Education
	Documents    []*Document
}

type PersonalData struct {
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

type Contact struct {
	ID                 string
	EmployeeID         string
	PersonalEmail      string
	PhoneNumber        string
	IdentityAddress    string
	ResidentialAddress string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Bank struct {
	ID                string
	EmployeeID        string
	BankName          string
	AccountNumber     string
	AccountHolderName string
	IsPrimary         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Education struct {
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

type Document struct {
	ID           string
	EmployeeID   string
	DocumentType string
	DocumentURL  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewEmployee is the factory function enforcing invariants on creation
func NewEmployee(employeeCode, jobPositionID, employmentStatus string, joinDate time.Time) (*Employee, error) {
	if employeeCode == "" || jobPositionID == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Employee{
		ID:               uuid.NewString(),
		EmployeeCode:     employeeCode,
		JobPositionID:    jobPositionID,
		EmploymentStatus: employmentStatus,
		JoinDate:         joinDate,
		Status:           "ACTIVE",
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func NewPersonalData(employeeID, fullName, ktpNumber, gender, maritalStatus, ptkpStatus, religion string) (*PersonalData, error) {
	if employeeID == "" || fullName == "" || ktpNumber == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &PersonalData{
		ID:            uuid.NewString(),
		EmployeeID:    employeeID,
		FullName:      fullName,
		KtpNumber:     ktpNumber,
		Gender:        gender,
		MaritalStatus: maritalStatus,
		PtkpStatus:    ptkpStatus,
		Religion:      religion,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func NewContact(employeeID, personalEmail, phoneNumber, identityAddress, residentialAddress string) (*Contact, error) {
	if employeeID == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Contact{
		ID:                 uuid.NewString(),
		EmployeeID:         employeeID,
		PersonalEmail:      personalEmail,
		PhoneNumber:        phoneNumber,
		IdentityAddress:    identityAddress,
		ResidentialAddress: residentialAddress,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func NewBank(employeeID, bankName, accountNumber, accountHolderName string, isPrimary bool) (*Bank, error) {
	if employeeID == "" || bankName == "" || accountNumber == "" || accountHolderName == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Bank{
		ID:                uuid.NewString(),
		EmployeeID:        employeeID,
		BankName:          bankName,
		AccountNumber:     accountNumber,
		AccountHolderName: accountHolderName,
		IsPrimary:         isPrimary,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func NewEducation(employeeID, level, institutionName, major string, startYear, endYear int, score float64) (*Education, error) {
	if employeeID == "" || level == "" || institutionName == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Education{
		ID:              uuid.NewString(),
		EmployeeID:      employeeID,
		Level:           level,
		InstitutionName: institutionName,
		Major:           major,
		StartYear:       startYear,
		EndYear:         endYear,
		Score:           score,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func NewDocument(employeeID, documentType, documentURL string) (*Document, error) {
	if employeeID == "" || documentURL == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Document{
		ID:           uuid.NewString(),
		EmployeeID:   employeeID,
		DocumentType: documentType,
		DocumentURL:  documentURL,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
