package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJobPositionNotFound        = errors.New("job position not found")
	ErrJobPositionCompanyMismatch = errors.New("department and job title belong to different companies")
	ErrReportingCompanyMismatch   = errors.New("reports-to position belongs to a different company")
)

// JobPosition adalah kursi jabatan aktual = Department x JobTitle. CompanyID
// didenormalisasi dari Department, bukan input client langsung (decision-log.md
// ADR-001) — tetap divalidasi non-empty di constructor seperti field wajib lain.
type JobPosition struct {
	ID             string
	CompanyID      string
	DepartmentID   string
	JobTitleID     string
	Name           string
	ReportsToID    *string
	HeadcountQuota int
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewJobPosition adalah satu-satunya tempat generate UUID (single source, uuid-generation.md).
// headcountQuota default 1 kalau < 1 (decision-log.md ADR-003) — bukan reject.
func NewJobPosition(companyID, departmentID, jobTitleID, name string, reportsToID *string, headcountQuota int) (*JobPosition, error) {
	if companyID == "" || departmentID == "" || jobTitleID == "" || name == "" {
		return nil, ErrInvalidInput
	}
	if headcountQuota < 1 {
		headcountQuota = 1
	}
	now := time.Now()
	return &JobPosition{
		ID:             uuid.NewString(),
		CompanyID:      companyID,
		DepartmentID:   departmentID,
		JobTitleID:     jobTitleID,
		Name:           name,
		ReportsToID:    reportsToID,
		HeadcountQuota: headcountQuota,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
