package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput              = errors.New("invalid input")
	ErrDepartmentNotFound        = errors.New("department not found")
	ErrDepartmentCodeDuplicate   = errors.New("department code already registered in this company")
	ErrDepartmentCompanyMismatch = errors.New("department does not belong to the given company")
)

// Department adalah unit kerja/divisi — Company-owned (scoping-convention.md §1),
// self-referencing hierarki via ParentID.
type Department struct {
	ID        string
	CompanyID string
	Code      string
	Name      string
	ParentID  *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDepartment adalah satu-satunya tempat generate UUID (single source, uuid-generation.md).
// companyID WAJIB — Department tak boleh lahir tanpa scope Company yang valid.
func NewDepartment(companyID, code, name string, parentID *string) (*Department, error) {
	if companyID == "" || code == "" || name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Department{
		ID:        uuid.NewString(),
		CompanyID: companyID,
		Code:      code,
		Name:      name,
		ParentID:  parentID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
