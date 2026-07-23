package di

import (
	authHTTP "github.com/bagusyanuar/hris-backend/internal/auth/transport/http"
	httpEmployee "github.com/bagusyanuar/hris-backend/internal/interfaces/http/employee"
	httpOrg "github.com/bagusyanuar/hris-backend/internal/interfaces/http/organization"
	orgHTTP "github.com/bagusyanuar/hris-backend/internal/organization/transport/http"
	"github.com/gofiber/fiber/v3"
)

type APIHandlers struct {
	Auth         *authHTTP.Handler
	Org          *httpOrg.Handler
	Employee     *httpEmployee.Handler
	Organization *orgHTTP.Handler
}

// RegisterRoutes registers all modules' HTTP routes to the Fiber router
func (h *APIHandlers) RegisterRoutes(router fiber.Router) {
	authHTTP.RegisterRoutes(router, h.Auth)
	httpOrg.RegisterRoutes(router, h.Org)
	httpEmployee.RegisterRoutes(router, h.Employee)
	orgHTTP.RegisterRoutes(router, h.Organization)
}
