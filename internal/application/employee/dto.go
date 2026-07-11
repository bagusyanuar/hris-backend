package employee

import "time"

type CreateEmployeeRequest struct {
	EmployeeCode     string `json:"employee_code" validate:"required"`
	JobPositionID    string `json:"job_position_id" validate:"required,uuid"`
	EmploymentStatus string `json:"employment_status" validate:"required,oneof=PERMANENT CONTRACT PROBATION INTERNSHIP"`
	JoinDate         string `json:"join_date" validate:"required,datetime=2006-01-02"`
}

type CreateEmployeeResponse struct {
	ID               string    `json:"id"`
	EmployeeCode     string    `json:"employee_code"`
	EmploymentStatus string    `json:"employment_status"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type UpdatePersonalDataRequest struct {
	FullName      string `json:"full_name" validate:"required"`
	KtpNumber     string `json:"ktp_number" validate:"required,len=16"`
	Gender        string `json:"gender" validate:"omitempty,oneof=MALE FEMALE"`
	MaritalStatus string `json:"marital_status" validate:"omitempty"`
	PtkpStatus    string `json:"ptkp_status" validate:"omitempty"`
	Religion      string `json:"religion" validate:"omitempty"`
}

type UpdateContactRequest struct {
	PersonalEmail      string `json:"personal_email" validate:"omitempty,email"`
	PhoneNumber        string `json:"phone_number" validate:"omitempty"`
	IdentityAddress    string `json:"identity_address" validate:"omitempty"`
	ResidentialAddress string `json:"residential_address" validate:"omitempty"`
}

type BankDTO struct {
	BankName          string `json:"bank_name" validate:"required"`
	AccountNumber     string `json:"account_number" validate:"required"`
	AccountHolderName string `json:"account_holder_name" validate:"required"`
	IsPrimary         bool   `json:"is_primary"`
}

type SaveBanksRequest struct {
	Banks []BankDTO `json:"banks" validate:"required,min=1,dive"`
}

type PersonalDataResponse struct {
	FullName      string `json:"full_name"`
	KtpNumber     string `json:"ktp_number"`
	Gender        string `json:"gender"`
	MaritalStatus string `json:"marital_status"`
	PtkpStatus    string `json:"ptkp_status"`
	Religion      string `json:"religion"`
}

type ContactResponse struct {
	PersonalEmail      string `json:"personal_email"`
	PhoneNumber        string `json:"phone_number"`
	IdentityAddress    string `json:"identity_address"`
	ResidentialAddress string `json:"residential_address"`
}

type BankResponse struct {
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
	IsPrimary         bool   `json:"is_primary"`
}

type GetEmployeeDetailResponse struct {
	ID               string                `json:"id"`
	EmployeeCode     string                `json:"employee_code"`
	JobPositionID    string                `json:"job_position_id"`
	EmploymentStatus string                `json:"employment_status"`
	JoinDate         string                `json:"join_date"`
	Status           string                `json:"status"`
	PersonalData     *PersonalDataResponse `json:"personal_data"`
	Contact          *ContactResponse      `json:"contact"`
	Banks            []BankResponse        `json:"banks"`
}
