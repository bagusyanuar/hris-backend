package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/workforce/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID      `gorm:"type:uuid;not null;index:idx_departments_company_id"`
	Code      string         `gorm:"type:varchar(20);not null"`
	Name      string         `gorm:"type:varchar(150);not null"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index:idx_departments_parent_id"`
	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (DepartmentModel) TableName() string { return "departments" }

// BeforeCreate = jaring pengaman UUID (hanya jalan pada Create(), bukan Save()).
func (m *DepartmentModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// ToDomain merekonstruksi entity LANGSUNG (tidak lewat constructor,
// supaya tidak generate UUID baru & tidak menjalankan ulang validasi create).
func (m *DepartmentModel) ToDomain() *domain.Department {
	return &domain.Department{
		ID:        m.ID.String(),
		CompanyID: m.CompanyID.String(),
		Code:      m.Code,
		Name:      m.Name,
		ParentID:  uuidPtrToStringPtr(m.ParentID),
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func DepartmentFromDomain(d *domain.Department) (*DepartmentModel, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, err
	}
	companyID, err := uuid.Parse(d.CompanyID)
	if err != nil {
		return nil, err
	}
	parentID, err := stringPtrToUUIDPtr(d.ParentID)
	if err != nil {
		return nil, err
	}
	return &DepartmentModel{
		ID:        id,
		CompanyID: companyID,
		Code:      d.Code,
		Name:      d.Name,
		ParentID:  parentID,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}

type JobTitleModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID  uuid.UUID      `gorm:"type:uuid;not null;index:idx_job_titles_company_id"`
	Code       string         `gorm:"type:varchar(20);not null"`
	Name       string         `gorm:"type:varchar(100);not null"`
	GradeLevel int            `gorm:"not null"`
	IsActive   bool           `gorm:"not null;default:true"`
	CreatedAt  time.Time      `gorm:"not null;default:now()"`
	UpdatedAt  time.Time      `gorm:"not null;default:now()"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (JobTitleModel) TableName() string { return "job_titles" }

// BeforeCreate = jaring pengaman UUID (hanya jalan pada Create(), bukan Save()).
func (m *JobTitleModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// ToDomain merekonstruksi entity LANGSUNG (tidak lewat constructor,
// supaya tidak generate UUID baru & tidak menjalankan ulang validasi create).
func (m *JobTitleModel) ToDomain() *domain.JobTitle {
	return &domain.JobTitle{
		ID:         m.ID.String(),
		CompanyID:  m.CompanyID.String(),
		Code:       m.Code,
		Name:       m.Name,
		GradeLevel: m.GradeLevel,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func JobTitleFromDomain(jt *domain.JobTitle) (*JobTitleModel, error) {
	id, err := uuid.Parse(jt.ID)
	if err != nil {
		return nil, err
	}
	companyID, err := uuid.Parse(jt.CompanyID)
	if err != nil {
		return nil, err
	}
	return &JobTitleModel{
		ID:         id,
		CompanyID:  companyID,
		Code:       jt.Code,
		Name:       jt.Name,
		GradeLevel: jt.GradeLevel,
		IsActive:   jt.IsActive,
		CreatedAt:  jt.CreatedAt,
		UpdatedAt:  jt.UpdatedAt,
	}, nil
}

type JobPositionModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_job_positions_company_id"`
	DepartmentID   uuid.UUID      `gorm:"type:uuid;not null;index:idx_job_positions_department_id"`
	JobTitleID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_job_positions_job_title_id"`
	Name           string         `gorm:"type:varchar(150);not null"`
	ReportsToID    *uuid.UUID     `gorm:"type:uuid;index:idx_job_positions_reports_to_id"`
	HeadcountQuota int            `gorm:"not null;default:1"`
	IsActive       bool           `gorm:"not null;default:true"`
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (JobPositionModel) TableName() string { return "job_positions" }

// BeforeCreate = jaring pengaman UUID (hanya jalan pada Create(), bukan Save()).
func (m *JobPositionModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// ToDomain merekonstruksi entity LANGSUNG (tidak lewat constructor,
// supaya tidak generate UUID baru & tidak menjalankan ulang validasi create).
func (m *JobPositionModel) ToDomain() *domain.JobPosition {
	return &domain.JobPosition{
		ID:             m.ID.String(),
		CompanyID:      m.CompanyID.String(),
		DepartmentID:   m.DepartmentID.String(),
		JobTitleID:     m.JobTitleID.String(),
		Name:           m.Name,
		ReportsToID:    uuidPtrToStringPtr(m.ReportsToID),
		HeadcountQuota: m.HeadcountQuota,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func JobPositionFromDomain(jp *domain.JobPosition) (*JobPositionModel, error) {
	id, err := uuid.Parse(jp.ID)
	if err != nil {
		return nil, err
	}
	companyID, err := uuid.Parse(jp.CompanyID)
	if err != nil {
		return nil, err
	}
	departmentID, err := uuid.Parse(jp.DepartmentID)
	if err != nil {
		return nil, err
	}
	jobTitleID, err := uuid.Parse(jp.JobTitleID)
	if err != nil {
		return nil, err
	}
	reportsToID, err := stringPtrToUUIDPtr(jp.ReportsToID)
	if err != nil {
		return nil, err
	}
	return &JobPositionModel{
		ID:             id,
		CompanyID:      companyID,
		DepartmentID:   departmentID,
		JobTitleID:     jobTitleID,
		Name:           jp.Name,
		ReportsToID:    reportsToID,
		HeadcountQuota: jp.HeadcountQuota,
		IsActive:       jp.IsActive,
		CreatedAt:      jp.CreatedAt,
		UpdatedAt:      jp.UpdatedAt,
	}, nil
}

// uuidPtrToStringPtr & stringPtrToUUIDPtr adalah helper konversi kolom self-referencing
// nullable (ParentID/ReportsToID) antara uuid.UUID (GORM) dan string (domain).
func uuidPtrToStringPtr(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}

func stringPtrToUUIDPtr(s *string) (*uuid.UUID, error) {
	if s == nil {
		return nil, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
