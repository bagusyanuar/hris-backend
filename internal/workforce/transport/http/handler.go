package http

import (
	"errors"
	"strconv"

	orgDomain "github.com/bagusyanuar/hris-backend/internal/organization/domain"
	"github.com/bagusyanuar/hris-backend/internal/workforce/application"
	"github.com/bagusyanuar/hris-backend/internal/workforce/domain"
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

// --- Department handlers ---

func (h *Handler) CreateDepartment(c fiber.Ctx) error {
	var req application.CreateDepartmentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.CreateDepartment(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, orgDomain.ErrCompanyNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrDepartmentNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrDepartmentCompanyMismatch):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrDepartmentCodeDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to create department")
		}
	}
	return response.Success(c, fiber.StatusCreated, "Department created successfully", res)
}

func (h *Handler) ListDepartments(c fiber.Ctx) error {
	page, limit, sort, order := parsePagination(c)
	search := c.Query("search")
	res, err := h.service.ListDepartments(c.Context(), page, limit, sort, order, search)
	if err != nil {
		return serverError(c, err, "Failed to fetch departments")
	}
	return response.SuccessList(c, fiber.StatusOK, "Departments fetched successfully", res.Items, res.Meta)
}

// TreeDepartments — GET /departments/tree, TANPA pagination. Dipakai FE Tabel
// (nested row, expand/collapse) & Bagan (tree diagram).
func (h *Handler) TreeDepartments(c fiber.Ctx) error {
	res, err := h.service.ListDepartmentsTree(c.Context())
	if err != nil {
		return serverError(c, err, "Failed to fetch department tree")
	}
	return response.Success(c, fiber.StatusOK, "Department tree fetched successfully", res)
}

func (h *Handler) GetDepartment(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetDepartment(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch department")
	}
	return response.Success(c, fiber.StatusOK, "Department fetched successfully", res)
}

func (h *Handler) UpdateDepartment(c fiber.Ctx) error {
	id := c.Params("id")
	var req application.UpdateDepartmentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.UpdateDepartment(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDepartmentNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrDepartmentCompanyMismatch):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrHierarchyCycle):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrDepartmentCodeDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to update department")
		}
	}
	return response.Success(c, fiber.StatusOK, "Department updated successfully", res)
}

func (h *Handler) DeleteDepartment(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteDepartment(c.Context(), id); err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to delete department")
	}
	return response.Success(c, fiber.StatusOK, "Department deleted successfully", nil)
}

// --- Job Title handlers ---

func (h *Handler) CreateJobTitle(c fiber.Ctx) error {
	var req application.CreateJobTitleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.CreateJobTitle(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, orgDomain.ErrCompanyNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrJobTitleCodeDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to create job title")
		}
	}
	return response.Success(c, fiber.StatusCreated, "Job title created successfully", res)
}

func (h *Handler) ListJobTitles(c fiber.Ctx) error {
	page, limit, sort, order := parsePagination(c)
	search := c.Query("search")
	res, err := h.service.ListJobTitles(c.Context(), page, limit, sort, order, search)
	if err != nil {
		return serverError(c, err, "Failed to fetch job titles")
	}
	return response.SuccessList(c, fiber.StatusOK, "Job titles fetched successfully", res.Items, res.Meta)
}

func (h *Handler) GetJobTitle(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetJobTitle(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrJobTitleNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch job title")
	}
	return response.Success(c, fiber.StatusOK, "Job title fetched successfully", res)
}

func (h *Handler) UpdateJobTitle(c fiber.Ctx) error {
	id := c.Params("id")
	var req application.UpdateJobTitleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.UpdateJobTitle(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrJobTitleNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrJobTitleCodeDuplicate):
			return response.Error(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to update job title")
		}
	}
	return response.Success(c, fiber.StatusOK, "Job title updated successfully", res)
}

func (h *Handler) DeleteJobTitle(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteJobTitle(c.Context(), id); err != nil {
		if errors.Is(err, domain.ErrJobTitleNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to delete job title")
	}
	return response.Success(c, fiber.StatusOK, "Job title deleted successfully", nil)
}

// --- Job Position handlers ---

func (h *Handler) CreateJobPosition(c fiber.Ctx) error {
	var req application.CreateJobPositionRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.CreateJobPosition(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDepartmentNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrJobTitleNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrJobPositionNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrJobPositionCompanyMismatch):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrReportingCompanyMismatch):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrHierarchyCycle):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to create job position")
		}
	}
	return response.Success(c, fiber.StatusCreated, "Job position created successfully", res)
}

func (h *Handler) ListJobPositions(c fiber.Ctx) error {
	page, limit, sort, order := parsePagination(c)
	search := c.Query("search")
	res, err := h.service.ListJobPositions(c.Context(), page, limit, sort, order, search)
	if err != nil {
		return serverError(c, err, "Failed to fetch job positions")
	}
	return response.SuccessList(c, fiber.StatusOK, "Job positions fetched successfully", res.Items, res.Meta)
}

// ChartJobPositions — GET /job-positions/chart, TANPA pagination (decision-log.md ADR-004).
func (h *Handler) ChartJobPositions(c fiber.Ctx) error {
	res, err := h.service.ListJobPositionsChart(c.Context())
	if err != nil {
		return serverError(c, err, "Failed to fetch job position chart")
	}
	return response.Success(c, fiber.StatusOK, "Job position chart fetched successfully", res)
}

func (h *Handler) GetJobPosition(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetJobPosition(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrJobPositionNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to fetch job position")
	}
	return response.Success(c, fiber.StatusOK, "Job position fetched successfully", res)
}

func (h *Handler) UpdateJobPosition(c fiber.Ctx) error {
	id := c.Params("id")
	var req application.UpdateJobPositionRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	res, err := h.service.UpdateJobPosition(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrJobPositionNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrDepartmentNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrJobTitleNotFound):
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrJobPositionCompanyMismatch):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrReportingCompanyMismatch):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		case errors.Is(err, domain.ErrHierarchyCycle):
			return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		default:
			return serverError(c, err, "Failed to update job position")
		}
	}
	return response.Success(c, fiber.StatusOK, "Job position updated successfully", res)
}

func (h *Handler) DeleteJobPosition(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteJobPosition(c.Context(), id); err != nil {
		if errors.Is(err, domain.ErrJobPositionNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return serverError(c, err, "Failed to delete job position")
	}
	return response.Success(c, fiber.StatusOK, "Job position deleted successfully", nil)
}
