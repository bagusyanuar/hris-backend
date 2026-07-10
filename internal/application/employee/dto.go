package employee

import "time"

type CreateEmployeeRequest struct {
	EmployeeCode     string `json:"employee_code" validate:"required"`
	JobPositionID    string `json:"job_position_id" validate:"required,uuid"`
	JoinDate         string `json:"join_date" validate:"required"` // Format: YYYY-MM-DD
	EmploymentStatus string `json:"employment_status" validate:"required,oneof=PERMANENT CONTRACT PROBATION INTERNSHIP"`

	PersonalData PersonalDataRequest `json:"personal_data" validate:"required"`
	Banks        []BankRequest       `json:"banks" validate:"required,min=1"`
}

type PersonalDataRequest struct {
	FullName      string  `json:"full_name" validate:"required,max=150"`
	KtpNumber     string  `json:"ktp_number" validate:"required,len=16,numeric"`
	Gender        *string `json:"gender" validate:"omitempty,oneof=MALE FEMALE"`
	MaritalStatus *string `json:"marital_status" validate:"omitempty"`
	PtkpStatus    *string `json:"ptkp_status" validate:"omitempty"`
	Religion      *string `json:"religion" validate:"omitempty"`
}

type BankRequest struct {
	BankName          string `json:"bank_name" validate:"required"`
	AccountNumber     string `json:"account_number" validate:"required"`
	AccountHolderName string `json:"account_holder_name" validate:"required"`
	IsPrimary         bool   `json:"is_primary"`
}

type EmployeeResponse struct {
	ID               string    `json:"id"`
	EmployeeCode     string    `json:"employee_code"`
	Status           string    `json:"status"`
	EmploymentStatus string    `json:"employment_status"`
	CreatedAt        time.Time `json:"created_at"`
}
