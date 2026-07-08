package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/domain/user"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string     `gorm:"type:varchar(255);not null"`
	Password  string     `gorm:"type:varchar(255);not null"`
	Status    string     `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt time.Time  `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `gorm:"type:timestamp with time zone;default:null"`
}

// TableName menentukan nama tabel di database
func (UserModel) TableName() string {
	return "users"
}

// ToDomain mengonversi GORM model ke Domain Entity
func (m *UserModel) ToDomain() (*user.User, error) {
	return user.NewUser(m.ID.String(), m.Email, m.Password, m.Status)
}

// FromDomain mengonversi Domain Entity ke GORM model
func FromDomain(u *user.User) (*UserModel, error) {
	parsedID, err := uuid.Parse(u.ID())
	if err != nil {
		return nil, err
	}
	return &UserModel{
		ID:       parsedID,
		Email:    u.Email(),
		Password: u.Password(),
		Status:   u.Status(),
	}, nil
}
