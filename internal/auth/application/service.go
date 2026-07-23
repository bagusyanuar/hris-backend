package application

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/auth/domain"
	user "github.com/bagusyanuar/hris-backend/internal/user/domain"
	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type Service struct {
	userRepo       user.Repository
	tokenGenerator domain.TokenGenerator
}

func NewService(userRepo user.Repository, tokenGenerator domain.TokenGenerator) *Service {
	return &Service{
		userRepo:       userRepo,
		tokenGenerator: tokenGenerator,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		logger.FromContext(ctx).Error("failed to find user by email during login", zap.Error(err))
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password()), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Assuming static role "employee" for now until RBAC is implemented
	return s.tokenGenerator.GenerateTokenPair(u.ID(), "employee")
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	claims, err := s.tokenGenerator.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Optional: Check if user still exists and active
	_, err = s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		logger.FromContext(ctx).Error("failed to find user by id during token refresh", zap.Error(err))
		return nil, ErrInvalidToken
	}

	return s.tokenGenerator.GenerateTokenPair(claims.UserID, claims.Role)
}
