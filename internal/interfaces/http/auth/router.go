package auth

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes registers the authentication routes.
func RegisterRoutes(router fiber.Router, handler *Handler) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.Refresh)
}
