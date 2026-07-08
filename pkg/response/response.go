package response

import (
	"github.com/gofiber/fiber/v3"
)

// Response merepresentasikan format balikan API yang sukses.
type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse merepresentasikan format balikan API ketika terjadi error.
type ErrorResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func getStatusFromCode(code int) string {
	if code >= 200 && code < 300 {
		return "success"
	}
	return "error"
}

// Success mengirimkan HTTP response sukses dengan format standar.
func Success(c fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Status:  getStatusFromCode(code),
		Message: message,
		Data:    data,
	})
}

// Error mengirimkan HTTP response error dengan format standar.
func Error(c fiber.Ctx, code int, message string, errors interface{}) error {
	return c.Status(code).JSON(ErrorResponse{
		Code:    code,
		Status:  getStatusFromCode(code),
		Message: message,
		Errors:  errors,
	})
}
