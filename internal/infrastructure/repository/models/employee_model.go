package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeModel struct {
	ID               string         `gorm:"primaryKey;type:varchar(50)"`
	UserID           *string        `gorm:"type:varchar(50)"`
	EmployeeCode     string         `gorm:"type:varchar(50);unique"`
	JobPositionID    string         `gorm:"type:varchar(50);not null"`
	EmploymentStatus string         `gorm:"type:varchar(20)"`
	JoinDate         time.Time      `gorm:"type:date;not null"`
	EndDate          *time.Time     `gorm:"type:date"`
	ResignDate       *time.Time     `gorm:"type:date"`
	Status           string         `gorm:"type:varchar(20);default:'ACTIVE'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Relations
	PersonalData *EmployeePersonalDataModel `gorm:"foreignKey:EmployeeID"`
	Contact      *EmployeeContactModel      `gorm:"foreignKey:EmployeeID"`
	Banks        []EmployeeBankModel        `gorm:"foreignKey:EmployeeID"`
	Educations   []EmployeeEducationModel   `gorm:"foreignKey:EmployeeID"`
	Documents    []EmployeeDocumentModel    `gorm:"foreignKey:EmployeeID"`
}

func (EmployeeModel) TableName() string { return "employees" }

func (m *EmployeeModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (m *EmployeeModel) ToDomain() *employee.Employee {
	var uid string
	if m.UserID != nil {
		uid = *m.UserID
	}
	emp := &employee.Employee{
		ID:               m.ID,
		UserID:           uid,
		EmployeeCode:     m.EmployeeCode,
		JobPositionID:    m.JobPositionID,
		EmploymentStatus: m.EmploymentStatus,
		JoinDate:         m.JoinDate,
		ResignDate:       m.ResignDate,
		Status:           m.Status,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}

	if m.PersonalData != nil {
		emp.PersonalData = m.PersonalData.ToDomain()
	}
	if m.Contact != nil {
		emp.Contact = m.Contact.ToDomain()
	}
	if len(m.Banks) > 0 {
		var banks []*employee.Bank
		for _, b := range m.Banks {
			banks = append(banks, b.ToDomain())
		}
		emp.Banks = banks
	}
	if len(m.Educations) > 0 {
		var educations []*employee.Education
		for _, e := range m.Educations {
			educations = append(educations, e.ToDomain())
		}
		emp.Educations = educations
	}
	if len(m.Documents) > 0 {
		var documents []*employee.Document
		for _, d := range m.Documents {
			documents = append(documents, d.ToDomain())
		}
		emp.Documents = documents
	}

	return emp
}

func EmployeeFromDomain(e *employee.Employee) *EmployeeModel {
	var uid *string
	if e.UserID != "" {
		uid = &e.UserID
	}
	return &EmployeeModel{
		ID:               e.ID,
		UserID:           uid,
		EmployeeCode:     e.EmployeeCode,
		JobPositionID:    e.JobPositionID,
		EmploymentStatus: e.EmploymentStatus,
		JoinDate:         e.JoinDate,
		ResignDate:       e.ResignDate,
		Status:           e.Status,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

// Personal Data
type EmployeePersonalDataModel struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	EmployeeID    string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	FullName      string    `gorm:"type:varchar(150);not null"`
	KtpNumber     string    `gorm:"type:varchar(16);unique;not null"`
	Gender        string    `gorm:"type:varchar(10)"`
	MaritalStatus string    `gorm:"type:varchar(20)"`
	PtkpStatus    string    `gorm:"type:varchar(10)"`
	Religion      string    `gorm:"type:varchar(30)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (EmployeePersonalDataModel) TableName() string { return "employee_personal_data" }

func (m *EmployeePersonalDataModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (m *EmployeePersonalDataModel) ToDomain() *employee.PersonalData {
	return &employee.PersonalData{
		ID:            m.ID,
		EmployeeID:    m.EmployeeID,
		FullName:      m.FullName,
		KtpNumber:     m.KtpNumber,
		Gender:        m.Gender,
		MaritalStatus: m.MaritalStatus,
		PtkpStatus:    m.PtkpStatus,
		Religion:      m.Religion,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func PersonalDataFromDomain(d *employee.PersonalData) *EmployeePersonalDataModel {
	return &EmployeePersonalDataModel{
		ID:            d.ID,
		EmployeeID:    d.EmployeeID,
		FullName:      d.FullName,
		KtpNumber:     d.KtpNumber,
		Gender:        d.Gender,
		MaritalStatus: d.MaritalStatus,
		PtkpStatus:    d.PtkpStatus,
		Religion:      d.Religion,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// Contacts
type EmployeeContactModel struct {
	ID                 string    `gorm:"primaryKey;type:varchar(50)"`
	EmployeeID         string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	PersonalEmail      string    `gorm:"type:varchar(100)"`
	PhoneNumber        string    `gorm:"type:varchar(20)"`
	IdentityAddress    string    `gorm:"type:text"`
	ResidentialAddress string    `gorm:"type:text"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (EmployeeContactModel) TableName() string { return "employee_contacts" }

func (m *EmployeeContactModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (m *EmployeeContactModel) ToDomain() *employee.Contact {
	return &employee.Contact{
		ID:                 m.ID,
		EmployeeID:         m.EmployeeID,
		PersonalEmail:      m.PersonalEmail,
		PhoneNumber:        m.PhoneNumber,
		IdentityAddress:    m.IdentityAddress,
		ResidentialAddress: m.ResidentialAddress,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func ContactFromDomain(c *employee.Contact) *EmployeeContactModel {
	return &EmployeeContactModel{
		ID:                 c.ID,
		EmployeeID:         c.EmployeeID,
		PersonalEmail:      c.PersonalEmail,
		PhoneNumber:        c.PhoneNumber,
		IdentityAddress:    c.IdentityAddress,
		ResidentialAddress: c.ResidentialAddress,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

// Banks
type EmployeeBankModel struct {
	ID                string    `gorm:"primaryKey;type:varchar(50)"`
	EmployeeID        string    `gorm:"type:varchar(50);not null"`
	BankName          string    `gorm:"type:varchar(50);not null"`
	AccountNumber     string    `gorm:"type:varchar(50);not null"`
	AccountHolderName string    `gorm:"type:varchar(100);not null"`
	IsPrimary         bool      `gorm:"default:true"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (EmployeeBankModel) TableName() string { return "employee_banks" }

func (m *EmployeeBankModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (m EmployeeBankModel) ToDomain() *employee.Bank {
	return &employee.Bank{
		ID:                m.ID,
		EmployeeID:        m.EmployeeID,
		BankName:          m.BankName,
		AccountNumber:     m.AccountNumber,
		AccountHolderName: m.AccountHolderName,
		IsPrimary:         m.IsPrimary,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func BankFromDomain(b *employee.Bank) *EmployeeBankModel {
	return &EmployeeBankModel{
		ID:                b.ID,
		EmployeeID:        b.EmployeeID,
		BankName:          b.BankName,
		AccountNumber:     b.AccountNumber,
		AccountHolderName: b.AccountHolderName,
		IsPrimary:         b.IsPrimary,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
}

// Education
type EmployeeEducationModel struct {
	ID              string    `gorm:"primaryKey;type:varchar(50)"`
	EmployeeID      string    `gorm:"type:varchar(50);not null"`
	Level           string    `gorm:"type:varchar(20)"`
	InstitutionName string    `gorm:"type:varchar(150)"`
	Major           string    `gorm:"type:varchar(100)"`
	StartYear       int       `gorm:"type:int"`
	EndYear         int       `gorm:"type:int"`
	Score           float64   `gorm:"type:decimal(5,2)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (EmployeeEducationModel) TableName() string { return "employee_educations" }

func (m *EmployeeEducationModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (m EmployeeEducationModel) ToDomain() *employee.Education {
	return &employee.Education{
		ID:              m.ID,
		EmployeeID:      m.EmployeeID,
		Level:           m.Level,
		InstitutionName: m.InstitutionName,
		Major:           m.Major,
		StartYear:       m.StartYear,
		EndYear:         m.EndYear,
		Score:           m.Score,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func EducationFromDomain(e *employee.Education) *EmployeeEducationModel {
	return &EmployeeEducationModel{
		ID:              e.ID,
		EmployeeID:      e.EmployeeID,
		Level:           e.Level,
		InstitutionName: e.InstitutionName,
		Major:           e.Major,
		StartYear:       e.StartYear,
		EndYear:         e.EndYear,
		Score:           e.Score,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// Documents
type EmployeeDocumentModel struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)"`
	EmployeeID   string    `gorm:"type:varchar(50);not null"`
	DocumentType string    `gorm:"type:varchar(50)"`
	DocumentURL  string    `gorm:"type:text;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (EmployeeDocumentModel) TableName() string { return "employee_documents" }

func (m *EmployeeDocumentModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (m EmployeeDocumentModel) ToDomain() *employee.Document {
	return &employee.Document{
		ID:           m.ID,
		EmployeeID:   m.EmployeeID,
		DocumentType: m.DocumentType,
		DocumentURL:  m.DocumentURL,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func DocumentFromDomain(d *employee.Document) *EmployeeDocumentModel {
	return &EmployeeDocumentModel{
		ID:           d.ID,
		EmployeeID:   d.EmployeeID,
		DocumentType: d.DocumentType,
		DocumentURL:  d.DocumentURL,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
