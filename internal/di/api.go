package di

import (
	authHTTP "github.com/bagusyanuar/hris-backend/internal/auth/transport/http"
	httpEmployee "github.com/bagusyanuar/hris-backend/internal/interfaces/http/employee"
	httpOrg "github.com/bagusyanuar/hris-backend/internal/interfaces/http/organization"
	"github.com/gofiber/fiber/v3"
)

type APIHandlers struct {
	Auth     *authHTTP.Handler
	Org      *httpOrg.Handler
	Employee *httpEmployee.Handler
}

// RegisterRoutes registers all modules' HTTP routes to the Fiber router
func (h *APIHandlers) RegisterRoutes(router fiber.Router) {
	authHTTP.RegisterRoutes(router, h.Auth)
	httpOrg.RegisterRoutes(router, h.Org)
	httpEmployee.RegisterRoutes(router, h.Employee)
}
