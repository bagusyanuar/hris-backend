---
name: scaffold-domain
description: Scaffold a new domain module/bounded context following the HRIS DDD rules (without CQRS)
---

# Scaffolding a New DDD Domain Module

Gunakan skill ini ketika user meminta untuk membuat modul domain baru (misalnya: `attendance`, `payroll`, `leave`).

## Prosedur Pembuatan Modul

Untuk setiap modul domain baru `<domain_name>`, pastikan mengimplementasikan *endpoints* sesuai dengan yang dirancang di `Tech Specs`. 

> **CRITICAL AI RULE (LAYER CONSISTENCY)**: You MUST implement EXACTLY the operations defined in the Technical Specifications across ALL layers consistently. If a module has 3 operations (e.g., Create, FindAll, GetByID), you must write the FULL logic for those 3 operations in the Repository, Application Service, HTTP Handler, and Router. DO NOT use placeholders (e.g., `// ... tambahkan sisanya`) or drop operations to save tokens. Generate everything completely!

1. **Domain Layer** (`internal/domain/<domain_name>/`):
   - `entity.go`: Definisikan entity utama, value objects, constructor `New<EntityName>`, dan validasi bisnis. Biasakan ada `CreatedAt`, `UpdatedAt`, dan `IsActive`.
   - `repository.go`: Definisikan interface `Repository` yang menggunakan `context.Context` untuk operasi yang dibutuhkan.
   - `service.go` (Opsional): Buat jika ada logika bisnis murni domain.

2. **Application Layer** (`internal/application/<domain_name>/`):
   - `service.go`: Definisikan application service untuk mengkoordinasikan transaksi.
   - `dto.go`: Request & Response DTOs. Untuk `Update...Request`, gunakan pointer (e.g., `*bool`, `*string`) untuk field opsional agar bisa membedakan `null` dengan *zero value*.

3. **Infrastructure Layer** (`internal/infrastructure/repository/`):
   - `<domain_name>_postgres.go` (atau DB target lainnya): Implementasikan interface repository dari domain.
   - `models/<domain_name>_model.go`: Model GORM dengan fungsi konversi `ToDomain()` dan `FromDomain()`.

4. **Interfaces Layer** (`internal/interfaces/http/`):
   - `<domain_name>/handler.go`: HTTP handler (Gin/Fiber) untuk operasi yang diizinkan.
   - `<domain_name>/router.go`: Register routes.

5. **Dependency Injection (Google Wire)** (`internal/di/`):
   - Tambahkan struct Handler ke dalam `APIHandlers` di `internal/di/api.go` dan panggil fungsi `RegisterRoutes` di dalamnya.
   - Tambahkan fungsi constructor Repository, Service, dan Handler ke masing-masing `ProviderSet` di dalam `internal/di/wire.go`.
   - Jalankan perintah `go run github.com/google/wire/cmd/wire@latest ./internal/di` di terminal untuk men-generate `wire_gen.go`.

## Contoh Template Entity & Repository

### `domain/<domain_name>/entity.go`
```go
package <domain_name>

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid input")

type <EntityName> struct {
	ID        string
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New<EntityName>(id string, name string) (*<EntityName>, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}
	if id == "" {
		id = uuid.New().String()
	}
	now := time.Now()
	return &<EntityName>{
		ID:        id,
		Name:      name,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
```

### `domain/<domain_name>/repository.go`
```go
package <domain_name>

import "context"

type Repository interface {
	Save(ctx context.Context, item *<EntityName>) error
	FindByID(ctx context.Context, id string) (*<EntityName>, error)
	FindAll(ctx context.Context) ([]*<EntityName>, error)
	Update(ctx context.Context, item *<EntityName>) error
	Delete(ctx context.Context, id string) error
}
```

### `internal/infrastructure/repository/models/<domain_name>_model.go`
```go
package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/domain/<domain_name>"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type <EntityName>Model struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	Name      string    `gorm:"type:varchar(100);not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName menentukan nama tabel di database
func (<EntityName>Model) TableName() string {
	return "<domain_name>s"
}

// BeforeCreate untuk memastikan UUID selalu digenerate jika kosong
func (m *<EntityName>Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// ToDomain mengonversi GORM model ke Domain Entity
func (m *<EntityName>Model) ToDomain() (*<domain_name>.<EntityName>, error) {
	entity, _ := <domain_name>.New<EntityName>(m.ID, m.Name)
	entity.IsActive = m.IsActive
	entity.CreatedAt = m.CreatedAt
	entity.UpdatedAt = m.UpdatedAt
	return entity, nil
}

// <EntityName>FromDomain mengonversi Domain Entity ke GORM model
func <EntityName>FromDomain(entity *<domain_name>.<EntityName>) *<EntityName>Model {
	return &<EntityName>Model{
		ID:        entity.ID,
		Name:      entity.Name,
		IsActive:  entity.IsActive,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
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
	model := models.<EntityName>FromDomain(item)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *<EntityName>Repository) FindByID(ctx context.Context, id string) (*<domain_name>.<EntityName>, error) {
	var model models.<EntityName>Model
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain()
}

func (r *<EntityName>Repository) FindAll(ctx context.Context) ([]*<domain_name>.<EntityName>, error) {
	var dbModels []models.<EntityName>Model
	if err := r.db.WithContext(ctx).Find(&dbModels).Error; err != nil {
		return nil, err
	}
	var domains []*<domain_name>.<EntityName>
	for _, m := range dbModels {
		d, _ := m.ToDomain()
		domains = append(domains, d)
	}
	return domains, nil
}

func (r *<EntityName>Repository) Update(ctx context.Context, item *<domain_name>.<EntityName>) error {
	model := models.<EntityName>FromDomain(item)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *<EntityName>Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.<EntityName>Model{}, "id = ?", id).Error
}
```

### `internal/interfaces/http/<domain_name>/handler.go`
```go
package <domain_name>

import (
	"github.com/gofiber/fiber/v3"
	app<EntityName> "github.com/bagusyanuar/hris-backend/internal/application/<domain_name>"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
)

type Handler struct {
	service *app<EntityName>.Service
}

func NewHandler(service *app<EntityName>.Service) *Handler {
	return &Handler{service: service}
}

// PERHATIAN: Untuk validasi request, JANGAN inject validator ke struct Handler.
// Gunakan fungsi global pkg/validator langsung, contoh: `errs := validator.ValidateStruct(req)`

// WAJIB KONSISTEN! Implementasikan HANYA DAN SEMUA operasi yang disetujui di Tech Specs.
// DILARANG MEMOTONG KODE ATAU MENGGUNAKAN PLACEHOLDER.
// TULIS SEMUA METHOD HANDLER SECARA LENGKAP.
func (h *Handler) Create(c fiber.Ctx) error { /* full impl */ }
func (h *Handler) Get(c fiber.Ctx) error { /* full impl */ }
// ... sertakan method lain (FindAll, Update, Delete) HANYA jika memang diwajibkan oleh spesifikasi.
```
