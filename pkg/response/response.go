package response

import (
	"github.com/bagusyanuar/hris-backend/pkg/pagination"
	"github.com/gofiber/fiber/v3"
)

// Response merepresentasikan format balikan API yang sukses.
type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Meta menampung metadata response di luar payload Data, dikelompokkan per
// concern (mis. Pagination) supaya nambah metadata baru (request tracing,
// summary, dst) di masa depan gak nabrak/refactor field yang udah ada.
type Meta struct {
	Pagination *pagination.Meta `json:"pagination,omitempty"`
}

// ErrorResponse merepresentasikan format balikan API ketika terjadi error.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

func getStatusFromCode(code int) string {
	if code >= 200 && code < 300 {
		return "success"
	}
	return "error"
}

// Success mengirimkan HTTP response sukses dengan format standar.
func Success(c fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Status:  getStatusFromCode(code),
		Message: message,
		Data:    data,
	})
}

// SuccessList mengirimkan HTTP response sukses untuk endpoint List/FindAll.
// Data WAJIB payload murni (array item), bukan wrapper {items, meta} —
// pagination masuk Meta.Pagination (selaras pagination-convention.md).
func SuccessList(c fiber.Ctx, code int, message string, items any, pageMeta pagination.Meta) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Status:  getStatusFromCode(code),
		Message: message,
		Data:    items,
		Meta:    &Meta{Pagination: &pageMeta},
	})
}

// Error mengirimkan HTTP response error dengan format standar.
func Error(c fiber.Ctx, code int, message string, errors any) error {
	return c.Status(code).JSON(ErrorResponse{
		Code:    code,
		Status:  getStatusFromCode(code),
		Message: message,
		Errors:  errors,
	})
}
