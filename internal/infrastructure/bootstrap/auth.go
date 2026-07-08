package bootstrap

import (
	applicationAuth "github.com/bagusyanuar/hris-backend/internal/application/auth"
	domainAuth "github.com/bagusyanuar/hris-backend/internal/domain/auth"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository"
	httpAuth "github.com/bagusyanuar/hris-backend/internal/interfaces/http/auth"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// InitAuthModule initializes the dependencies and routes for the Auth module.
func InitAuthModule(db *gorm.DB, api fiber.Router, tokenGen domainAuth.TokenGenerator) {
	// Initialize Repository
	userRepo := repository.NewUserRepository(db)

	// Initialize Application Service
	authService := applicationAuth.NewService(userRepo, tokenGen)

	// Initialize HTTP Handler
	authHandler := httpAuth.NewHandler(authService)

	// Register Routes
	httpAuth.RegisterRoutes(api, authHandler)
}
