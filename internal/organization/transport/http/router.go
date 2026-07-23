package http

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	companies := router.Group("/companies")
	companies.Post("/", h.CreateCompany)
	companies.Get("/", h.ListCompanies)
	companies.Get("/:id", h.GetCompany)
	companies.Put("/:id", h.UpdateCompany)
	companies.Delete("/:id", h.DeleteCompany)

	companies.Post("/:companyId/branches", h.CreateBranch)
	companies.Get("/:companyId/branches", h.ListBranchesByCompany)

	branches := router.Group("/branches")
	branches.Get("/:id", h.GetBranch)
	branches.Put("/:id", h.UpdateBranch)
	branches.Delete("/:id", h.DeleteBranch)
}
