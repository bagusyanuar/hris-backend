package employee

import (
	"github.com/bagusyanuar/hris-backend/internal/application/employee"
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
