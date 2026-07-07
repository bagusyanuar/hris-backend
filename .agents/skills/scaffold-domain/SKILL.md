---
name: scaffold-domain
description: Scaffold a new domain module/bounded context following the HRIS DDD rules (without CQRS)
---

# Scaffolding a New DDD Domain Module

Gunakan skill ini ketika user meminta untuk membuat modul domain baru (misalnya: `attendance`, `payroll`, `leave`).

## Prosedur Pembuatan Modul

Untuk setiap modul domain baru `<domain_name>`, buat file-file berikut:

1. **Domain Layer** (`internal/domain/<domain_name>/`):
   - `entity.go`: Definisikan entity utama, value objects, constructor `New<EntityName>`, dan validasi bisnis.
   - `repository.go`: Definisikan interface `Repository` yang menggunakan `context.Context`.
   - `service.go` (Opsional): Buat jika ada logika bisnis yang mengoordinasikan beberapa entity.

2. **Application Layer** (`internal/application/<domain_name>/`):
   - `service.go`: Definisikan application service untuk mengkoordinasikan transaksi dan use cases (read/write terpadu, tidak memakai CQRS).

3. **Infrastructure Layer** (`internal/infrastructure/repository/`):
   - `<domain_name>_postgres.go` (atau DB target lainnya): Implementasikan interface repository dari domain.

4. **Interfaces Layer** (`internal/interfaces/http/`):
   - `<domain_name>_handler.go`: HTTP handler (Gin/Fiber) untuk memetakan request, validasi input dasar, memanggil application service, dan mengembalikan response JSON.

## Contoh Template Entity & Repository

### `domain/<domain_name>/entity.go`
```go
package <domain_name>

import (
	"context"
	"errors"
)

var ErrInvalidInput = errors.New("invalid input")

type <EntityName> struct {
	id   string
	name string
}

func New<EntityName>(id string, name string) (*<EntityName>, error) {
	if id == "" || name == "" {
		return nil, ErrInvalidInput
	}
	return &<EntityName>{
		id:   id,
		name: name,
	}, nil
}

func (e *<EntityName>) ID() string {
	return e.id
}

func (e *<EntityName>) Name() string {
	return e.name
}
```

### `domain/<domain_name>/repository.go`
```go
package <domain_name>

import "context"

type Repository interface {
	Save(ctx context.Context, item *<EntityName>) error
	FindByID(ctx context.Context, id string) (*<EntityName>, error)
}
```

### `internal/infrastructure/repository/models/<domain_name>_model.go`
```go
package models

import "github.com/bagusyanuar/hris-backend/internal/domain/<domain_name>"

type <EntityName>Model struct {
	ID   string `gorm:"primaryKey;type:varchar(50)"`
	Name string `gorm:"type:varchar(100);not null"`
}

// TableName menentukan nama tabel di database
func (<EntityName>Model) TableName() string {
	return "<domain_name>s"
}

// ToDomain mengonversi GORM model ke Domain Entity
func (m *<EntityName>Model) ToDomain() (*<domain_name>.<EntityName>, error) {
	return <domain_name>.New<EntityName>(m.ID, m.Name)
}

// FromDomain mengonversi Domain Entity ke GORM model
func FromDomain(entity *<domain_name>.<EntityName>) *<EntityName>Model {
	return &<EntityName>Model{
		ID:   entity.ID(),
		Name: entity.Name(),
	}
}
```

### `internal/infrastructure/repository/<domain_name>_postgres.go`
```go
package repository

import (
	"context"
	"github.com/bagusyanuar/hris-backend/internal/domain/<domain_name>"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository/models"
	"gorm.io/gorm"
)

type <EntityName>Repository struct {
	db *gorm.DB
}

func New<EntityName>Repository(db *gorm.DB) <domain_name>.Repository {
	return &<EntityName>Repository{db: db}
}

func (r *<EntityName>Repository) Save(ctx context.Context, item *<domain_name>.<EntityName>) error {
	model := models.FromDomain(item)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *<EntityName>Repository) FindByID(ctx context.Context, id string) (*<domain_name>.<EntityName>, error) {
	var model models.<EntityName>Model
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain()
}
```

### `internal/interfaces/http/<domain_name>_handler.go`
```go
package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/bagusyanuar/hris-backend/internal/application/<domain_name>"
)

type <EntityName>Handler struct {
	service *<domain_name>.Service
}

func New<EntityName>Handler(service *<domain_name>.Service) *<EntityName>Handler {
	return &<EntityName>Handler{service: service}
}

func (h *<EntityName>Handler) Get(c fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	
	result, err := h.service.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(result)
}
```
