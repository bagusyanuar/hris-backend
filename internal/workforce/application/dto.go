package application

import (
	"time"

	"github.com/bagusyanuar/hris-backend/pkg/pagination"
)

// --- Department DTOs ---

type CreateDepartmentRequest struct {
	CompanyID string  `json:"company_id" validate:"required,uuid4"`
	Code      string  `json:"code" validate:"required,max=20"`
	Name      string  `json:"name" validate:"required,max=150"`
	ParentID  *string `json:"parent_id" validate:"omitempty,uuid4"`
}

// UpdateDepartmentRequest pakai pointer supaya partial update bisa membedakan
// "field tidak dikirim" dengan zero value.
type UpdateDepartmentRequest struct {
	Code     *string `json:"code" validate:"omitempty,max=20"`
	Name     *string `json:"name" validate:"omitempty,max=150"`
	ParentID *string `json:"parent_id" validate:"omitempty,uuid4"`
	IsActive *bool   `json:"is_active"`
}

type DepartmentResponse struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DepartmentListResponse struct {
	Items []DepartmentResponse `json:"items"`
	Meta  pagination.Meta      `json:"meta"`
}

// --- Job Title DTOs ---

type CreateJobTitleRequest struct {
	CompanyID  string `json:"company_id" validate:"required,uuid4"`
	Code       string `json:"code" validate:"required,max=20"`
	Name       string `json:"name" validate:"required,max=100"`
	GradeLevel int    `json:"grade_level" validate:"required"`
}

// UpdateJobTitleRequest pakai pointer supaya partial update bisa membedakan
// "field tidak dikirim" dengan zero value.
type UpdateJobTitleRequest struct {
	Code       *string `json:"code" validate:"omitempty,max=20"`
	Name       *string `json:"name" validate:"omitempty,max=100"`
	GradeLevel *int    `json:"grade_level"`
	IsActive   *bool   `json:"is_active"`
}

type JobTitleResponse struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	GradeLevel int       `json:"grade_level"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type JobTitleListResponse struct {
	Items []JobTitleResponse `json:"items"`
	Meta  pagination.Meta    `json:"meta"`
}

// --- Job Position DTOs ---

// CreateJobPositionRequest sengaja TIDAK punya field CompanyID — diturunkan
// otomatis dari DepartmentID di application service (decision-log.md ADR-001).
type CreateJobPositionRequest struct {
	DepartmentID   string  `json:"department_id" validate:"required,uuid4"`
	JobTitleID     string  `json:"job_title_id" validate:"required,uuid4"`
	Name           string  `json:"name" validate:"required,max=150"`
	ReportsToID    *string `json:"reports_to_id" validate:"omitempty,uuid4"`
	HeadcountQuota int     `json:"headcount_quota" validate:"omitempty,min=1"`
}

// UpdateJobPositionRequest pakai pointer supaya partial update bisa membedakan
// "field tidak dikirim" dengan zero value.
type UpdateJobPositionRequest struct {
	DepartmentID   *string `json:"department_id" validate:"omitempty,uuid4"`
	JobTitleID     *string `json:"job_title_id" validate:"omitempty,uuid4"`
	Name           *string `json:"name" validate:"omitempty,max=150"`
	ReportsToID    *string `json:"reports_to_id" validate:"omitempty,uuid4"`
	HeadcountQuota *int    `json:"headcount_quota" validate:"omitempty,min=1"`
	IsActive       *bool   `json:"is_active"`
}

// JobPositionRef adalah preload ringan {id, name} — dipakai FE tampilin nama
// Department/Job Title tanpa lookup terpisah (hindari raw UUID doang di tabel/form).
type JobPositionRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JobPositionResponse struct {
	ID             string         `json:"id"`
	CompanyID      string         `json:"company_id"`
	Department     JobPositionRef `json:"department"`
	JobTitle       JobPositionRef `json:"job_title"`
	Name           string         `json:"name"`
	ReportsToID    *string        `json:"reports_to_id"`
	HeadcountQuota int            `json:"headcount_quota"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type JobPositionListResponse struct {
	Items []JobPositionResponse `json:"items"`
	Meta  pagination.Meta       `json:"meta"`
}
