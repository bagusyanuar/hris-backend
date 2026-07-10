package employee

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	api := router.Group("/api/v1/employees")
	api.Post("/", h.Create)
	api.Get("/:id", h.Get)
}
