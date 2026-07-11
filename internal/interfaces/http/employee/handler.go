package employee

import (
	"errors"

	appEmployee "github.com/bagusyanuar/hris-backend/internal/application/employee"
	domainEmployee "github.com/bagusyanuar/hris-backend/internal/domain/employee"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *appEmployee.Service
}

func NewHandler(service *appEmployee.Service) *Handler {
	return &Handler{service: service}
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
		return response.Error(c, 500, err.Error(), nil)
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
		return response.Error(c, 500, err.Error(), nil)
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
		return response.Error(c, 500, err.Error(), nil)
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
		return response.Error(c, 500, err.Error(), nil)
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
		return response.Error(c, 500, err.Error(), nil)
	}

	return response.Success(c, 200, "Employee detail fetched successfully", res)
}
