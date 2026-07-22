---
name: api-validation
description: Guide for adding API request validation and DTOs in HRIS Backend
---

# API Validation & DTOs

Gunakan skill ini ketika user meminta untuk membuat Data Transfer Object (DTO), menambahkan validasi pada endpoint API, atau membuat request struct.

## 1. Aturan Pembuatan DTO (Request Struct)
Setiap struct yang digunakan untuk menerima *payload* JSON HTTP wajib memenuhi kriteria berikut:
- Tentukan `json` tag.
- Tentukan `validate` tag menggunakan format `go-playground/validator/v10`.
- Jangan menggunakan package validasi lain.

Contoh:
```go
type CreateDepartmentRequest struct {
	Code string `json:"code" validate:"required,min=2,max=10"`
	Name string `json:"name" validate:"required,max=100"`
}
```

## 2. Cara Eksekusi Validasi di Handler
Panggil *wrapper* yang ada di `pkg/validator` untuk memvalidasi struct, setelah JSON di-*bind*. Jika gagal, WAJIB mengembalikan HTTP Status `422 Unprocessable Entity` menggunakan `response.Error`.

Contoh implementasi di Fiber handler (`internal/<domain>/transport/http/handler.go`):
```go
import (
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateDepartmentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Lakukan validasi
	if validationErrs := validator.ValidateStruct(req); validationErrs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validation failed", validationErrs)
	}

	// Lanjutkan proses bisnis
	// ...
}
```

## 3. Dokumentasi (Bruno & Swagger)
Setiap menambahkan atau memodifikasi endpoint dengan validasi, JANGAN LUPA untuk mengupdate dokumentasi API (baik `.bru` maupun `swagger.yaml`):

### Swagger YAML:
Wajib menambahkan blok `422` di dalam `responses`:
```yaml
        '422':
          description: Unprocessable Entity - Validation failed
          content:
            application/json:
              example:
                code: 422
                status: "error"
                message: "Validation failed"
                errors:
                  code: ["Kolom code wajib diisi.", "Kolom code minimal harus berisi 2 karakter."]
```

### Bruno `.bru`:
Wajib menambahkan blok `422` di dalam blok `docs`:
```markdown
  **422 Unprocessable Entity**
  Gagal validasi input.
  ```json
  {
    "code": 422,
    "status": "error",
    "message": "Validation failed",
    "errors": {
      "code": [
        "Kolom code wajib diisi."
      ]
    }
  }
  ```
```
