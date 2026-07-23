package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBranchNotFound        = errors.New("branch not found")
	ErrBranchCodeDuplicate   = errors.New("branch code already registered in this company")
	ErrBranchCompanyMismatch = errors.New("branch does not belong to the given company")
)

// Branch adalah location root — Company-owned, scope wajib CompanyID
// (scoping-convention.md §1). Aggregate root sendiri, bukan child entity
// Company (decision-log.md ADR-003).
type Branch struct {
	ID        string
	CompanyID string
	Code      string
	Name      string
	City      *string
	IsMain    bool
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBranch adalah satu-satunya tempat generate UUID (single source, uuid-generation.md).
// companyID WAJIB — Branch tak boleh lahir tanpa scope Company yang valid.
func NewBranch(companyID, code, name string, city *string, isMain bool) (*Branch, error) {
	if companyID == "" || code == "" || name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Branch{
		ID:        uuid.NewString(),
		CompanyID: companyID,
		Code:      code,
		Name:      name,
		City:      city,
		IsMain:    isMain,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
