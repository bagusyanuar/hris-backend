package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyNPWPDuplicate = errors.New("company npwp already registered")
)

// Company adalah legal root (badan hukum/PT) — tidak punya CompanyID,
// dia sendiri adalah scope (scoping-convention.md §1).
type Company struct {
	ID        string
	Code      string
	LegalName string
	Npwp      *string // nullable — belum dipakai (decision-log.md ADR-005)
	BpjsNo    *string // nullable — belum dipakai (decision-log.md ADR-005)
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCompany adalah satu-satunya tempat generate UUID (single source, uuid-generation.md).
func NewCompany(code, legalName string, npwp, bpjsNo *string) (*Company, error) {
	if code == "" || legalName == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &Company{
		ID:        uuid.NewString(),
		Code:      code,
		LegalName: legalName,
		Npwp:      npwp,
		BpjsNo:    bpjsNo,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
