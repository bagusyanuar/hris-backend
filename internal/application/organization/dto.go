package organization

// DTOs for Department
type CreateDepartmentRequest struct {
	Code     string  `json:"code" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	ParentID *string `json:"parent_id,omitempty"`
}

type UpdateDepartmentRequest struct {
	Code     string  `json:"code" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	ParentID *string `json:"parent_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type DepartmentResponse struct {
	ID       string  `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
	IsActive bool    `json:"is_active"`
}

// DTOs for JobTitle
type CreateJobTitleRequest struct {
	Code       string `json:"code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	GradeLevel int    `json:"grade_level" validate:"required"`
}

type UpdateJobTitleRequest struct {
	Code       string `json:"code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	GradeLevel int    `json:"grade_level" validate:"required"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type JobTitleResponse struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	GradeLevel int    `json:"grade_level"`
	IsActive   bool   `json:"is_active"`
}

// DTOs for JobPosition
type CreateJobPositionRequest struct {
	DepartmentID   string  `json:"department_id" validate:"required"`
	JobTitleID     string  `json:"job_title_id" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	ReportsToID    *string `json:"reports_to_id,omitempty"`
	HeadcountQuota int     `json:"headcount_quota"`
}

type UpdateJobPositionRequest struct {
	DepartmentID   string  `json:"department_id" validate:"required"`
	JobTitleID     string  `json:"job_title_id" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	ReportsToID    *string `json:"reports_to_id,omitempty"`
	HeadcountQuota int     `json:"headcount_quota"`
	IsActive       *bool   `json:"is_active,omitempty"`
}

type JobPositionResponse struct {
	ID             string  `json:"id"`
	DepartmentID   string  `json:"department_id"`
	JobTitleID     string  `json:"job_title_id"`
	Name           string  `json:"name"`
	ReportsToID    *string `json:"reports_to_id,omitempty"`
	HeadcountQuota int     `json:"headcount_quota"`
	IsActive       bool    `json:"is_active"`
}
