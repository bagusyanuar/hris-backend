package http

import (
	"errors"
	"strconv"

	"github.com/bagusyanuar/hris-backend/internal/organization/application"
	"github.com/bagusyanuar/hris-backend/internal/organization/domain"
	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// serverError logs the real error server-side and returns a generic 500 to
// the client, so internal detail (SQL/driver/column names) never leaks.
func serverError(c fiber.Ctx, err error, message string) error {
	logger.FromContext(c.Context()).Error(message, zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, message, nil)
}

func parsePagination(c fiber.Ctx) (page, limit int, sort, order string) {
	page, _ = strconv.Atoi(c.Query("page"))
	limit, _ = strconv.Atoi(c.Query("limit"))
	sort = c.Query("sort")
	order = c.Query("order")
	return
}

// --- Company handlers ---

func (h *Handler) CreateCompany(c fiber.Ctx) error {
	var req application.CreateCompanyRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.CreateCompany(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrCompanyNPWPDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to create company")
		}
	}
	return response.Success(c, fiber.StatusCreated, "Company created successfully", res)
}

func (h *Handler) ListCompanies(c fiber.Ctx) error {
	page, limit, sort, order := parsePagination(c)
	res, err := h.service.ListCompanies(c.Context(), page, limit, sort, order)
	if err != nil {
		return serverError(c, err, "Failed to fetch companies")
	}
	return response.Success(c, fiber.StatusOK, "Companies fetched successfully", res)
}

func (h *Handler) GetCompany(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetCompany(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch company")
	}
	return response.Success(c, fiber.StatusOK, "Company fetched successfully", res)
}

func (h *Handler) UpdateCompany(c fiber.Ctx) error {
	id := c.Params("id")
	var req application.UpdateCompanyRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.UpdateCompany(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCompanyNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrCompanyNPWPDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to update company")
		}
	}
	return response.Success(c, fiber.StatusOK, "Company updated successfully", res)
}

func (h *Handler) DeleteCompany(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteCompany(c.Context(), id); err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to delete company")
	}
	return response.Success(c, fiber.StatusOK, "Company deleted successfully", nil)
}

// --- Branch handlers ---

func (h *Handler) CreateBranch(c fiber.Ctx) error {
	companyID := c.Params("companyId")
	var req application.CreateBranchRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.CreateBranch(c.Context(), companyID, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCompanyNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrBranchCodeDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to create branch")
		}
	}
	return response.Success(c, fiber.StatusCreated, "Branch created successfully", res)
}

func (h *Handler) ListBranchesByCompany(c fiber.Ctx) error {
	companyID := c.Params("companyId")
	page, limit, sort, order := parsePagination(c)
	res, err := h.service.ListBranchesByCompany(c.Context(), companyID, page, limit, sort, order)
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch branches")
	}
	return response.Success(c, fiber.StatusOK, "Branches fetched successfully", res)
}

func (h *Handler) GetBranch(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetBranch(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBranchNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch branch")
	}
	return response.Success(c, fiber.StatusOK, "Branch fetched successfully", res)
}

func (h *Handler) UpdateBranch(c fiber.Ctx) error {
	id := c.Params("id")
	var req application.UpdateBranchRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.UpdateBranch(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBranchNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrBranchCodeDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to update branch")
		}
	}
	return response.Success(c, fiber.StatusOK, "Branch updated successfully", res)
}

func (h *Handler) DeleteBranch(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteBranch(c.Context(), id); err != nil {
		if errors.Is(err, domain.ErrBranchNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to delete branch")
	}
	return response.Success(c, fiber.StatusOK, "Branch deleted successfully", nil)
}
