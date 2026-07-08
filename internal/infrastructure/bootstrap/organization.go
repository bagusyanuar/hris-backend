package bootstrap

import (
	applicationOrg "github.com/bagusyanuar/hris-backend/internal/application/organization"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository"
	httpOrg "github.com/bagusyanuar/hris-backend/internal/interfaces/http/organization"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// InitOrganizationModule initializes the dependencies and routes for the Organization module.
func InitOrganizationModule(db *gorm.DB, api fiber.Router) {
	// Initialize Repository
	orgRepo := repository.NewOrganizationRepository(db)

	// Initialize Application Service
	orgService := applicationOrg.NewService(orgRepo)

	// Initialize HTTP Handler
	orgHandler := httpOrg.NewHandler(orgService)

	// Register Routes
	httpOrg.RegisterRoutes(api, orgHandler)
}
