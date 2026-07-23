package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/organization/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code      string         `gorm:"type:varchar(20);not null"`
	LegalName string         `gorm:"type:varchar(150);not null"`
	Npwp      *string        `gorm:"type:varchar(25)"`
	BpjsNo    *string        `gorm:"type:varchar(50)"`
	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CompanyModel) TableName() string { return "companies" }

// BeforeCreate = jaring pengaman UUID (hanya jalan pada Create(), bukan Save()).
func (m *CompanyModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// ToDomain merekonstruksi entity LANGSUNG (tidak lewat constructor,
// supaya tidak generate UUID baru & tidak menjalankan ulang validasi create).
func (m *CompanyModel) ToDomain() *domain.Company {
	return &domain.Company{
		ID:        m.ID.String(),
		Code:      m.Code,
		LegalName: m.LegalName,
		Npwp:      m.Npwp,
		BpjsNo:    m.BpjsNo,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func CompanyFromDomain(c *domain.Company) (*CompanyModel, error) {
	id, err := uuid.Parse(c.ID)
	if err != nil {
		return nil, err
	}
	return &CompanyModel{
		ID:        id,
		Code:      c.Code,
		LegalName: c.LegalName,
		Npwp:      c.Npwp,
		BpjsNo:    c.BpjsNo,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}

type BranchModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID      `gorm:"type:uuid;not null;index:idx_branches_company_id"`
	Code      string         `gorm:"type:varchar(20);not null"`
	Name      string         `gorm:"type:varchar(150);not null"`
	City      *string        `gorm:"type:varchar(100)"`
	IsMain    bool           `gorm:"not null;default:false"`
	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (BranchModel) TableName() string { return "branches" }

// BeforeCreate = jaring pengaman UUID (hanya jalan pada Create(), bukan Save()).
func (m *BranchModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// ToDomain merekonstruksi entity LANGSUNG (tidak lewat constructor,
// supaya tidak generate UUID baru & tidak menjalankan ulang validasi create).
func (m *BranchModel) ToDomain() *domain.Branch {
	return &domain.Branch{
		ID:        m.ID.String(),
		CompanyID: m.CompanyID.String(),
		Code:      m.Code,
		Name:      m.Name,
		City:      m.City,
		IsMain:    m.IsMain,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func BranchFromDomain(b *domain.Branch) (*BranchModel, error) {
	id, err := uuid.Parse(b.ID)
	if err != nil {
		return nil, err
	}
	companyID, err := uuid.Parse(b.CompanyID)
	if err != nil {
		return nil, err
	}
	return &BranchModel{
		ID:        id,
		CompanyID: companyID,
		Code:      b.Code,
		Name:      b.Name,
		City:      b.City,
		IsMain:    b.IsMain,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}
