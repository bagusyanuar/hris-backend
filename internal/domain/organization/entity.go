package organization

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("organization data not found")
)

type Department struct {
	ID        string
	Code      string
	Name      string
	ParentID  *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDepartment(id, code, name string, parentID *string) (*Department, error) {
	if code == "" || name == "" {
		return nil, ErrInvalidInput
	}
	if id == "" {
		id = uuid.New().String()
	}
	now := time.Now()
	return &Department{
		ID:        id,
		Code:      code,
		Name:      name,
		ParentID:  parentID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

type JobTitle struct {
	ID         string
	Code       string
	Name       string
	GradeLevel int
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewJobTitle(id, code, name string, gradeLevel int) (*JobTitle, error) {
	if code == "" || name == "" {
		return nil, ErrInvalidInput
	}
	if id == "" {
		id = uuid.New().String()
	}
	now := time.Now()
	return &JobTitle{
		ID:         id,
		Code:       code,
		Name:       name,
		GradeLevel: gradeLevel,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

type JobPosition struct {
	ID             string
	DepartmentID   string
	JobTitleID     string
	Name           string
	ReportsToID    *string
	HeadcountQuota int
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewJobPosition(id, departmentID, jobTitleID, name string, reportsToID *string, headcountQuota int) (*JobPosition, error) {
	if departmentID == "" || jobTitleID == "" || name == "" {
		return nil, ErrInvalidInput
	}
	if id == "" {
		id = uuid.New().String()
	}
	if headcountQuota < 1 {
		headcountQuota = 1
	}
	now := time.Now()
	return &JobPosition{
		ID:             id,
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
