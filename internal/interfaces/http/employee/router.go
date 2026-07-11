package employee

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	employeeGroup := router.Group("/employees")
	
	// Progressive Creation Endpoints
	employeeGroup.Post("/", h.CreateCore)
	employeeGroup.Get("/:id", h.GetByID)
	employeeGroup.Put("/:id/personal-data", h.UpdatePersonalData)
	employeeGroup.Put("/:id/contacts", h.UpdateContact)
	employeeGroup.Post("/:id/banks", h.SaveBanks)
}
