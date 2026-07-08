package auth

import (
	"time"

	authApp "github.com/bagusyanuar/hris-backend/internal/application/auth"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Handler struct {
	service *authApp.Service
}

func NewHandler(service *authApp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrs := validator.ValidateStruct(req); validationErrs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", validationErrs)
	}

	tokenPair, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials", err.Error())
	}

	// Set HttpOnly Cookie for Refresh Token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		Domain:   "", // Adjust as needed
		MaxAge:   7 * 24 * 60 * 60, // 7 days in seconds
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Secure:   true, // true for HTTPS
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return response.Success(c, fiber.StatusOK, "Login successful", fiber.Map{
		"access_token": tokenPair.AccessToken,
		"expires_in":   tokenPair.ExpiresIn,
		"token_type":   "Bearer",
	})
}

func (h *Handler) Refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Refresh token missing", nil)
	}

	tokenPair, err := h.service.Refresh(c.Context(), refreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired refresh token", err.Error())
	}

	// Set new HttpOnly Cookie for Refresh Token (Rotation)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		Domain:   "",
		MaxAge:   7 * 24 * 60 * 60,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token": tokenPair.AccessToken,
		"expires_in":   tokenPair.ExpiresIn,
		"token_type":   "Bearer",
	})
}
