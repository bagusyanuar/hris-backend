package employee

import (
	"github.com/bagusyanuar/hris-backend/internal/application/employee"
	domainEmployee "github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *employee.Service
}

func NewHandler(service *employee.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req employee.CreateEmployeeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	result, err := h.service.Create(c.Context(), req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "invalid input data" || err.Error() == "at least one primary bank account is required" {
			status = fiber.StatusBadRequest
		} else if err.Error() == "ktp number already exists" {
			status = fiber.StatusUnprocessableEntity
		}
		return response.Error(c, status, "Failed to create employee", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Employee created successfully", result)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "employee not found" {
			return response.Error(c, fiber.StatusNotFound, "Employee not found", nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve employee", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Successfully retrieved data", result)
}

func (h *Handler) FindAll(c fiber.Ctx) error {
	result, err := h.service.FindAll(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve employees", err.Error())
	}
	if result == nil {
		result = make([]*domainEmployee.Employee, 0)
	}
	return response.Success(c, fiber.StatusOK, "Successfully retrieved data", result)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var req employee.UpdateEmployeeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON format", err.Error())
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if err.Error() == "employee not found" {
			return response.Error(c, fiber.StatusNotFound, "Employee not found", nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update employee", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Employee updated successfully", nil)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete employee", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Employee deleted successfully", nil)
}
