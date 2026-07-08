package organization

import (
	appOrg "github.com/bagusyanuar/hris-backend/internal/application/organization"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *appOrg.Service
}

func NewHandler(service *appOrg.Service) *Handler {
	return &Handler{service: service}
}

// Department Handlers
func (h *Handler) CreateDepartment(c fiber.Ctx) error {
	ctx := c.Context()
	var req appOrg.CreateDepartmentRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.CreateDepartment(ctx, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create department", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Department created successfully", res)
}

func (h *Handler) GetDepartmentByID(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	res, err := h.service.GetDepartmentByID(ctx, id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Department not found", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved department", res)
}

func (h *Handler) GetAllDepartments(c fiber.Ctx) error {
	ctx := c.Context()
	res, err := h.service.GetAllDepartments(ctx)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve departments", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved all departments", res)
}

func (h *Handler) UpdateDepartment(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	var req appOrg.UpdateDepartmentRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.UpdateDepartment(ctx, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update department", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Department updated successfully", res)
}

func (h *Handler) DeleteDepartment(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if err := h.service.DeleteDepartment(ctx, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete department", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Department deleted successfully", []interface{}{})
}

// JobTitle Handlers
func (h *Handler) CreateJobTitle(c fiber.Ctx) error {
	ctx := c.Context()
	var req appOrg.CreateJobTitleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.CreateJobTitle(ctx, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create job title", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Job title created successfully", res)
}

func (h *Handler) GetJobTitleByID(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	res, err := h.service.GetJobTitleByID(ctx, id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Job title not found", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved job title", res)
}

func (h *Handler) GetAllJobTitles(c fiber.Ctx) error {
	ctx := c.Context()
	res, err := h.service.GetAllJobTitles(ctx)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve job titles", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved all job titles", res)
}

func (h *Handler) UpdateJobTitle(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	var req appOrg.UpdateJobTitleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.UpdateJobTitle(ctx, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update job title", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Job title updated successfully", res)
}

func (h *Handler) DeleteJobTitle(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if err := h.service.DeleteJobTitle(ctx, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete job title", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Job title deleted successfully", []interface{}{})
}

// JobPosition Handlers
func (h *Handler) CreateJobPosition(c fiber.Ctx) error {
	ctx := c.Context()
	var req appOrg.CreateJobPositionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.CreateJobPosition(ctx, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create job position", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Job position created successfully", res)
}

func (h *Handler) GetJobPositionByID(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	res, err := h.service.GetJobPositionByID(ctx, id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Job position not found", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved job position", res)
}

func (h *Handler) GetAllJobPositions(c fiber.Ctx) error {
	ctx := c.Context()
	res, err := h.service.GetAllJobPositions(ctx)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve job positions", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved all job positions", res)
}

func (h *Handler) UpdateJobPosition(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	var req appOrg.UpdateJobPositionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.UpdateJobPosition(ctx, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update job position", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Job position updated successfully", res)
}

func (h *Handler) DeleteJobPosition(c fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if err := h.service.DeleteJobPosition(ctx, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete job position", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Job position deleted successfully", []interface{}{})
}
