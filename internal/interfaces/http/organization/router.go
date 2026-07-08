package organization

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, handler *Handler) {
	orgGroup := router.Group("/organization")

	// Departments
	deptGroup := orgGroup.Group("/departments")
	deptGroup.Post("/", handler.CreateDepartment)
	deptGroup.Get("/", handler.GetAllDepartments)
	deptGroup.Get("/:id", handler.GetDepartmentByID)
	deptGroup.Put("/:id", handler.UpdateDepartment)
	deptGroup.Delete("/:id", handler.DeleteDepartment)

	// Job Titles
	titleGroup := orgGroup.Group("/job-titles")
	titleGroup.Post("/", handler.CreateJobTitle)
	titleGroup.Get("/", handler.GetAllJobTitles)
	titleGroup.Get("/:id", handler.GetJobTitleByID)
	titleGroup.Put("/:id", handler.UpdateJobTitle)
	titleGroup.Delete("/:id", handler.DeleteJobTitle)

	// Job Positions
	posGroup := orgGroup.Group("/job-positions")
	posGroup.Post("/", handler.CreateJobPosition)
	posGroup.Get("/", handler.GetAllJobPositions)
	posGroup.Get("/:id", handler.GetJobPositionByID)
	posGroup.Put("/:id", handler.UpdateJobPosition)
	posGroup.Delete("/:id", handler.DeleteJobPosition)
}
