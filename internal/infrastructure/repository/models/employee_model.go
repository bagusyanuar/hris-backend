package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeModel struct {
	ID               string     `gorm:"primaryKey;type:uuid"`
	UserID           string     `gorm:"type:uuid;not null"`
	EmployeeCode     string     `gorm:"type:varchar(50);unique"`
	JobPositionID    string     `gorm:"type:uuid;not null"`
	EmploymentStatus string     `gorm:"type:varchar(20)"`
	JoinDate         time.Time  `gorm:"type:date;not null"`
	ResignDate       *time.Time `gorm:"type:date"`
	Status           string     `gorm:"type:varchar(20);default:'ACTIVE'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Associations
	PersonalData *EmployeePersonalDataModel `gorm:"foreignKey:EmployeeID"`
	Banks        []EmployeeBankModel        `gorm:"foreignKey:EmployeeID"`
}

func (EmployeeModel) TableName() string { return "employees" }

type EmployeePersonalDataModel struct {
	ID            string `gorm:"primaryKey;type:uuid"`
	EmployeeID    string `gorm:"type:uuid;uniqueIndex"`
	FullName      string `gorm:"type:varchar(150);not null"`
	KtpNumber     string `gorm:"type:varchar(16);uniqueIndex"`
	Gender        string `gorm:"type:varchar(10)"`
	MaritalStatus string `gorm:"type:varchar(20)"`
	PtkpStatus    string `gorm:"type:varchar(10)"`
	Religion      string `gorm:"type:varchar(30)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (EmployeePersonalDataModel) TableName() string { return "employee_personal_data" }

type EmployeeBankModel struct {
	ID                string `gorm:"primaryKey;type:uuid"`
	EmployeeID        string `gorm:"type:uuid"`
	BankName          string `gorm:"type:varchar(50);not null"`
	AccountNumber     string `gorm:"type:varchar(50);not null"`
	AccountHolderName string `gorm:"type:varchar(100);not null"`
	IsPrimary         bool   `gorm:"default:true"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (EmployeeBankModel) TableName() string { return "employee_banks" }

// BeforeCreate hooks for UUIDs
func (m *EmployeeModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}
func (m *EmployeePersonalDataModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}
func (m *EmployeeBankModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

// Mappers
func EmployeeFromDomain(e *employee.Employee) *EmployeeModel {
	if e == nil {
		return nil
	}
	model := &EmployeeModel{
		ID:               e.ID,
		UserID:           e.UserID,
		EmployeeCode:     e.EmployeeCode,
		JobPositionID:    e.JobPositionID,
		EmploymentStatus: e.EmploymentStatus,
		JoinDate:         e.JoinDate,
		ResignDate:       e.ResignDate,
		Status:           e.Status,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}

	if e.PersonalData != nil {
		model.PersonalData = &EmployeePersonalDataModel{
			ID:            e.PersonalData.ID,
			EmployeeID:    e.PersonalData.EmployeeID,
			FullName:      e.PersonalData.FullName,
			KtpNumber:     e.PersonalData.KtpNumber,
			Gender:        e.PersonalData.Gender,
			MaritalStatus: e.PersonalData.MaritalStatus,
			PtkpStatus:    e.PersonalData.PtkpStatus,
			Religion:      e.PersonalData.Religion,
			CreatedAt:     e.PersonalData.CreatedAt,
			UpdatedAt:     e.PersonalData.UpdatedAt,
		}
	}

	for _, b := range e.Banks {
		model.Banks = append(model.Banks, EmployeeBankModel{
			ID:                b.ID,
			EmployeeID:        b.EmployeeID,
			BankName:          b.BankName,
			AccountNumber:     b.AccountNumber,
			AccountHolderName: b.AccountHolderName,
			IsPrimary:         b.IsPrimary,
			CreatedAt:         b.CreatedAt,
			UpdatedAt:         b.UpdatedAt,
		})
	}

	return model
}

func (m *EmployeeModel) ToDomain() *employee.Employee {
	e := &employee.Employee{
		ID:               m.ID,
		UserID:           m.UserID,
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
		e.PersonalData = &employee.EmployeePersonalData{
			ID:            m.PersonalData.ID,
			EmployeeID:    m.PersonalData.EmployeeID,
			FullName:      m.PersonalData.FullName,
			KtpNumber:     m.PersonalData.KtpNumber,
			Gender:        m.PersonalData.Gender,
			MaritalStatus: m.PersonalData.MaritalStatus,
			PtkpStatus:    m.PersonalData.PtkpStatus,
			Religion:      m.PersonalData.Religion,
			CreatedAt:     m.PersonalData.CreatedAt,
			UpdatedAt:     m.PersonalData.UpdatedAt,
		}
	}

	for _, b := range m.Banks {
		e.Banks = append(e.Banks, &employee.EmployeeBank{
			ID:                b.ID,
			EmployeeID:        b.EmployeeID,
			BankName:          b.BankName,
			AccountNumber:     b.AccountNumber,
			AccountHolderName: b.AccountHolderName,
			IsPrimary:         b.IsPrimary,
			CreatedAt:         b.CreatedAt,
			UpdatedAt:         b.UpdatedAt,
		})
	}
	return e
}
