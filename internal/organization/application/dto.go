package application

import (
	"time"

	"github.com/bagusyanuar/hris-backend/pkg/pagination"
)

// --- Company DTOs ---

type CreateCompanyRequest struct {
	Code      string  `json:"code" validate:"required,max=20"`
	LegalName string  `json:"legal_name" validate:"required,max=150"`
	Npwp      *string `json:"npwp" validate:"omitempty,max=25"`
	BpjsNo    *string `json:"bpjs_no" validate:"omitempty,max=50"`
}

// UpdateCompanyRequest pakai pointer supaya partial update bisa membedakan
// "field tidak dikirim" dengan zero value.
type UpdateCompanyRequest struct {
	Code      *string `json:"code" validate:"omitempty,max=20"`
	LegalName *string `json:"legal_name" validate:"omitempty,max=150"`
	Npwp      *string `json:"npwp" validate:"omitempty,max=25"`
	BpjsNo    *string `json:"bpjs_no" validate:"omitempty,max=50"`
	IsActive  *bool   `json:"is_active"`
}

type CompanyResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	LegalName string    `json:"legal_name"`
	Npwp      *string   `json:"npwp"`
	BpjsNo    *string   `json:"bpjs_no"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyListResponse struct {
	Items []CompanyResponse `json:"items"`
	Meta  pagination.Meta   `json:"meta"`
}

// --- Branch DTOs ---

type CreateBranchRequest struct {
	Code   string  `json:"code" validate:"required,max=20"`
	Name   string  `json:"name" validate:"required,max=150"`
	City   *string `json:"city" validate:"omitempty,max=100"`
	IsMain bool    `json:"is_main"`
}

// UpdateBranchRequest pakai pointer supaya partial update bisa membedakan
// "field tidak dikirim" dengan zero value.
type UpdateBranchRequest struct {
	Code     *string `json:"code" validate:"omitempty,max=20"`
	Name     *string `json:"name" validate:"omitempty,max=150"`
	City     *string `json:"city" validate:"omitempty,max=100"`
	IsMain   *bool   `json:"is_main"`
	IsActive *bool   `json:"is_active"`
}

type BranchResponse struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	City      *string   `json:"city"`
	IsMain    bool      `json:"is_main"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BranchListResponse struct {
	Items []BranchResponse `json:"items"`
	Meta  pagination.Meta  `json:"meta"`
}
