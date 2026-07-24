package di

import (
	authHTTP "github.com/bagusyanuar/hris-backend/internal/auth/transport/http"
	httpEmployee "github.com/bagusyanuar/hris-backend/internal/interfaces/http/employee"
	orgHTTP "github.com/bagusyanuar/hris-backend/internal/organization/transport/http"
	workforceHTTP "github.com/bagusyanuar/hris-backend/internal/workforce/transport/http"
	"github.com/gofiber/fiber/v3"
)

type APIHandlers struct {
	Auth         *authHTTP.Handler
	Employee     *httpEmployee.Handler
	Organization *orgHTTP.Handler
	Workforce    *workforceHTTP.Handler
}

// RegisterRoutes registers all modules' HTTP routes to the Fiber router
func (h *APIHandlers) RegisterRoutes(router fiber.Router) {
	authHTTP.RegisterRoutes(router, h.Auth)
	httpEmployee.RegisterRoutes(router, h.Employee)
	orgHTTP.RegisterRoutes(router, h.Organization)
	workforceHTTP.RegisterRoutes(router, h.Workforce)
}
