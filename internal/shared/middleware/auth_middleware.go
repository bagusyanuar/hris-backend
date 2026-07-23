package middleware

import (
	"strings"

	authDomain "github.com/bagusyanuar/hris-backend/internal/auth/domain"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

func AuthProtected(tokenGenerator authDomain.TokenGenerator) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "missing authorization header", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "invalid authorization format", nil)
		}

		tokenString := parts[1]
		claims, err := tokenGenerator.ValidateToken(tokenString, "access")
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired token", nil)
		}

		// Store user info in context for next handlers
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
