package http

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	departments := router.Group("/departments")
	departments.Post("/", h.CreateDepartment)
	departments.Get("/", h.ListDepartments)
	departments.Get("/tree", h.TreeDepartments) // WAJIB sebelum "/:id" biar "tree" gak ketangkep param id
	departments.Get("/:id", h.GetDepartment)
	departments.Put("/:id", h.UpdateDepartment)
	departments.Delete("/:id", h.DeleteDepartment)

	jobTitles := router.Group("/job-titles")
	jobTitles.Post("/", h.CreateJobTitle)
	jobTitles.Get("/", h.ListJobTitles)
	jobTitles.Get("/:id", h.GetJobTitle)
	jobTitles.Put("/:id", h.UpdateJobTitle)
	jobTitles.Delete("/:id", h.DeleteJobTitle)

	jobPositions := router.Group("/job-positions")
	jobPositions.Post("/", h.CreateJobPosition)
	jobPositions.Get("/", h.ListJobPositions)
	jobPositions.Get("/chart", h.ChartJobPositions) // WAJIB sebelum "/:id" biar "chart" gak ketangkep param id
	jobPositions.Get("/:id", h.GetJobPosition)
	jobPositions.Put("/:id", h.UpdateJobPosition)
	jobPositions.Delete("/:id", h.DeleteJobPosition)
}
