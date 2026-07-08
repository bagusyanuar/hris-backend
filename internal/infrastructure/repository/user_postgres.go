package repository

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/domain/user"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return model.ToDomain()
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return model.ToDomain()
}
