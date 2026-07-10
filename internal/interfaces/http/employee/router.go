package employee

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	api := router.Group("/api/v1/employees")
	api.Post("/", h.Create)
	api.Get("/", h.FindAll)
	api.Get("/:id", h.Get)
	api.Put("/:id", h.Update)
	api.Delete("/:id", h.Delete)
}
