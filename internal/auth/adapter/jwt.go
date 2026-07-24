package adapter

import (
	"errors"
	"time"

	"github.com/bagusyanuar/hris-backend/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(secretKey string, accessExpiryHours int, refreshExpiryHours int) domain.TokenGenerator {
	if accessExpiryHours <= 0 {
		accessExpiryHours = 1
	}
	if refreshExpiryHours <= 0 {
		refreshExpiryHours = 168
	}

	return &jwtService{
		secretKey:     []byte(secretKey),
		accessExpiry:  time.Duration(accessExpiryHours) * time.Hour,
		refreshExpiry: time.Duration(refreshExpiryHours) * time.Hour,
	}
}

func (s *jwtService) GenerateTokenPair(userID string, role string) (*domain.TokenPair, error) {
	now := time.Now()

	// 1. Generate Access Token
	accessClaims := &jwtClaims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	// 2. Generate Refresh Token
	refreshClaims := &jwtClaims{
		UserID: userID,
		Role:   role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.accessExpiry.Seconds()),
	}, nil
}

func (s *jwtService) ValidateToken(tokenString string, expectedType string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		if claims.Type != expectedType {
			return nil, errors.New("invalid token type")
		}
		return &domain.TokenClaims{
			UserID: claims.UserID,
			Role:   claims.Role,
			Type:   claims.Type,
		}, nil
	}

	return nil, errors.New("invalid token")
}
