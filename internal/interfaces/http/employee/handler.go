package employee

import (
	"errors"

	appEmployee "github.com/bagusyanuar/hris-backend/internal/application/employee"
	domainEmployee "github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type Handler struct {
	service *appEmployee.Service
}

func NewHandler(service *appEmployee.Service) *Handler {
	return &Handler{service: service}
}

// serverError logs the real error server-side and returns a generic 500 to
// the client, so internal detail (SQL/driver/column names) never leaks.
func serverError(c fiber.Ctx, err error, message string) error {
	logger.FromContext(c.Context()).Error(message, zap.Error(err))
	return response.Error(c, 500, message, nil)
}

func (h *Handler) CreateCore(c fiber.Ctx) error {
	var req appEmployee.CreateEmployeeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, 400, "Invalid JSON payload", nil)
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, 422, "Validation failed", errs)
	}

	res, err := h.service.CreateCore(c.Context(), req)
	if err != nil {
		return serverError(c, err, "Failed to create employee")
	}

	return response.Success(c, 201, "Core employee created successfully", res)
}

func (h *Handler) UpdatePersonalData(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "Employee ID is required", nil)
	}

	var req appEmployee.UpdatePersonalDataRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, 400, "Invalid JSON payload", nil)
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, 422, "Validation failed", errs)
	}

	err := h.service.UpdatePersonalData(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, domainEmployee.ErrEmployeeNotFound) {
			return response.Error(c, 404, err.Error(), nil)
		}
		if errors.Is(err, domainEmployee.ErrKTPDuplicate) {
			return response.Error(c, 409, err.Error(), nil)
		}
		return serverError(c, err, "Failed to update personal data")
	}

	return response.Success(c, 200, "Personal data updated successfully", nil)
}

func (h *Handler) UpdateContact(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "Employee ID is required", nil)
	}

	var req appEmployee.UpdateContactRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, 400, "Invalid JSON payload", nil)
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, 422, "Validation failed", errs)
	}

	err := h.service.UpdateContact(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, domainEmployee.ErrEmployeeNotFound) {
			return response.Error(c, 404, err.Error(), nil)
		}
		return serverError(c, err, "Failed to update contact")
	}

	return response.Success(c, 200, "Contact updated successfully", nil)
}

func (h *Handler) SaveBanks(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "Employee ID is required", nil)
	}

	var req appEmployee.SaveBanksRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, 400, "Invalid JSON payload", nil)
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, 422, "Validation failed", errs)
	}

	err := h.service.SaveBanks(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, domainEmployee.ErrEmployeeNotFound) {
			return response.Error(c, 404, err.Error(), nil)
		}
		if errors.Is(err, domainEmployee.ErrPrimaryBankRequired) {
			return response.Error(c, 400, err.Error(), nil)
		}
		return serverError(c, err, "Failed to save bank accounts")
	}

	return response.Success(c, 200, "Bank accounts saved successfully", nil)
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "Employee ID is required", nil)
	}

	res, err := h.service.GetEmployeeDetail(c.Context(), id)
	if err != nil {
		if errors.Is(err, domainEmployee.ErrEmployeeNotFound) {
			return response.Error(c, 404, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch employee detail")
	}

	return response.Success(c, 200, "Employee detail fetched successfully", res)
}
