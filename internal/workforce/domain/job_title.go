package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJobTitleNotFound      = errors.New("job title not found")
	ErrJobTitleCodeDuplicate = errors.New("job title code already registered in this company")
)

// JobTitle adalah master grade/pangkat — Company-owned (scoping-convention.md §1),
// tiap PT punya grade sendiri (bukan Global master, PRD §7.2).
type JobTitle struct {
	ID         string
	CompanyID  string
	Code       string
	Name       string
	GradeLevel int
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewJobTitle adalah satu-satunya tempat generate UUID (single source, uuid-generation.md).
// companyID WAJIB — JobTitle tak boleh lahir tanpa scope Company yang valid.
func NewJobTitle(companyID, code, name string, gradeLevel int) (*JobTitle, error) {
	if companyID == "" || code == "" || name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &JobTitle{
		ID:         uuid.NewString(),
		CompanyID:  companyID,
		Code:       code,
		Name:       name,
		GradeLevel: gradeLevel,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
