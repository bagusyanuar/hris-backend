package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/domain/organization"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentModel struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Code      string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	ParentID  *string   `gorm:"type:uuid"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time
}

func (DepartmentModel) TableName() string {
	return "departments"
}

func (m *DepartmentModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

func (m *DepartmentModel) ToDomain() *organization.Department {
	return &organization.Department{
		ID:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		ParentID:  m.ParentID,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func DepartmentFromDomain(d *organization.Department) *DepartmentModel {
	return &DepartmentModel{
		ID:        d.ID,
		Code:      d.Code,
		Name:      d.Name,
		ParentID:  d.ParentID,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type JobTitleModel struct {
	ID         string    `gorm:"primaryKey;type:uuid"`
	Code       string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name       string    `gorm:"type:varchar(255);not null"`
	GradeLevel int       `gorm:"not null;index"`
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	DeletedAt  *time.Time
}

func (JobTitleModel) TableName() string {
	return "job_titles"
}

func (m *JobTitleModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

func (m *JobTitleModel) ToDomain() *organization.JobTitle {
	return &organization.JobTitle{
		ID:         m.ID,
		Code:       m.Code,
		Name:       m.Name,
		GradeLevel: m.GradeLevel,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func JobTitleFromDomain(t *organization.JobTitle) *JobTitleModel {
	return &JobTitleModel{
		ID:         t.ID,
		Code:       t.Code,
		Name:       t.Name,
		GradeLevel: t.GradeLevel,
		IsActive:   t.IsActive,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

type JobPositionModel struct {
	ID             string    `gorm:"primaryKey;type:uuid"`
	DepartmentID   string    `gorm:"type:uuid;not null;index"`
	JobTitleID     string    `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"type:varchar(255);not null"`
	ReportsToID    *string   `gorm:"type:uuid"`
	HeadcountQuota int       `gorm:"default:1"`
	IsActive       bool      `gorm:"default:true"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time
}

func (JobPositionModel) TableName() string {
	return "job_positions"
}

func (m *JobPositionModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

func (m *JobPositionModel) ToDomain() *organization.JobPosition {
	return &organization.JobPosition{
		ID:             m.ID,
		DepartmentID:   m.DepartmentID,
		JobTitleID:     m.JobTitleID,
		Name:           m.Name,
		ReportsToID:    m.ReportsToID,
		HeadcountQuota: m.HeadcountQuota,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func JobPositionFromDomain(p *organization.JobPosition) *JobPositionModel {
	return &JobPositionModel{
		ID:             p.ID,
		DepartmentID:   p.DepartmentID,
		JobTitleID:     p.JobTitleID,
		Name:           p.Name,
		ReportsToID:    p.ReportsToID,
		HeadcountQuota: p.HeadcountQuota,
		IsActive:       p.IsActive,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
