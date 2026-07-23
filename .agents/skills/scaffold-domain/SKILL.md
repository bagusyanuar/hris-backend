---
name: scaffold-domain
description: Scaffold a new domain module/bounded context following the HRIS DDD rules (without CQRS)
---

# Scaffolding a New DDD Domain Module (Domain-First)

Gunakan skill ini ketika user meminta untuk membuat modul domain baru (misalnya: `attendance`, `payroll`, `leave`).

Project ini memakai layout **domain-first**: satu bounded context = **satu folder utuh** di `internal/<domain_name>/` berisi keempat layer di dalamnya. Patuhi [architecture.md](../../rules/architecture.md), [uuid-generation.md](../../rules/uuid-generation.md), [persistence-convention.md](../../rules/persistence-convention.md), [scoping-convention.md](../../rules/scoping-convention.md), dan [pagination-convention.md](../../rules/pagination-convention.md) (kalau modul punya endpoint List/FindAll).

## Step 0 — Klasifikasi Scope Entity (WAJIB, SEBELUM nulis kode)

Sebelum generate layer manapun, **klasifikasikan tiap entity** sesuai [scoping-convention.md](../../rules/scoping-convention.md) §1 dan tentukan kolom scope-nya:

- **Company-owned** (default modul operasional) → entity/model/DTO/DBML bawa `CompanyID string` (`company_id` NOT NULL, FK `companies`).
- **Company + Location bound** (lokasi-spesifik, mis. Attendance) → tambah `BranchID string` (`branch_id` NOT NULL, FK `branches`).
- **Global master** (jarang, data identik lintas PT) → tanpa kolom scope, TAPI wajib justifikasi eksplisit.

Konsekuensi yang WAJIB ikut begitu entity company-owned:
1. Field `CompanyID` (+ `BranchID`) ada di **domain entity**, **GORM model**, **DTO**, dan **DBML**.
2. Constructor `New<Entity>` **menerima** `companyID` (dan `branchID` bila perlu) sebagai parameter wajib + validasi non-empty.
3. Repository `FindXxx`/`ListXxx` **scope-aware**: baca `scope.FromContext(ctx)` dan inject filter `WHERE company_id IN (...)`. JANGAN query baca tanpa filter scope.
4. Write memvalidasi `branch_id` se-`company_id` (mismatch → sentinel `ErrBranchCompanyMismatch`).

> **Staged (per scoping-convention.md §4):** modul Organization (`companies`/`branches`) & RBAC belum ada. Tetap tulis kolom + signature scope-aware sekarang; FK fisik & pengisian `scope.FromContext` nyusul saat modul itu landing. Jangan tunda kolomnya.

## Prosedur Pembuatan Modul

Untuk setiap modul domain baru `<domain_name>`, implementasikan *endpoints* sesuai `Tech Specs`.

> **CRITICAL AI RULE (LAYER CONSISTENCY)**: You MUST implement EXACTLY the operations defined in the Technical Specifications across ALL layers consistently. If a module has 3 operations (e.g., Create, FindAll, GetByID), you must write the FULL logic for those 3 operations in the Repository, Application Service, HTTP Handler, and Router. DO NOT use placeholders (e.g., `// ... tambahkan sisanya`) or drop operations to save tokens. Generate everything completely!

Struktur folder target untuk context baru:

```text
internal/<domain_name>/
├── domain/              # package domain
│   ├── entity.go
│   ├── repository.go
│   └── service.go       # opsional
├── application/         # package application
│   ├── service.go
│   └── dto.go
├── adapter/              # package adapter
│   ├── postgres.go
│   └── models/          # package models
│       └── <domain_name>_model.go
└── transport/
    └── http/            # package http
        ├── handler.go
        └── router.go
```

1. **Domain Layer** (`internal/<domain_name>/domain/`, package `domain`):
   - `entity.go`: entity utama, value objects, constructor `New<EntityName>`, validasi bisnis. Biasakan ada `CreatedAt`, `UpdatedAt`, `IsActive`. **UUID digenerate DI DALAM constructor** — constructor TIDAK menerima parameter `id` (single source of generation, lihat uuid-generation.md).
   - `repository.go`: interface `Repository` dengan `context.Context`. Not-found dikontrakkan sebagai sentinel error, BUKAN `(nil, nil)`. Kalau ada `FindAll`/List, signature **WAJIB** primitif `page, limit int, sort, order string` — **JANGAN** import `pkg/pagination` di domain (lihat pagination-convention.md §1, domain harus tetap pure Go tanpa dependency GORM).
   - `service.go` (opsional): logika bisnis murni domain.

2. **Application Layer** (`internal/<domain_name>/application/`, package `application`):
   - `service.go`: application service koordinasi transaksi. **DILARANG generate UUID di sini** — serahkan ke domain constructor. Untuk List: bungkus `page,limit,sort,order` jadi `pagination.Request{...}.Normalize()` sebelum diteruskan ke repo (pagination-convention.md §2).
   - `dto.go`: Request & Response DTOs. Untuk `Update...Request`, gunakan pointer (`*bool`, `*string`) agar bisa membedakan `null` dengan *zero value*. `List...Response` pakai field `Meta pagination.Meta` — JANGAN bikin struct meta lokal duplikat.

3. **Adapter Layer** (`internal/<domain_name>/adapter/`, package `adapter`):
   - `postgres.go`: implementasi interface repository. **Insert pakai `Create()`, JANGAN `Save()`** (lihat persistence-convention.md §1). Not-found → map `gorm.ErrRecordNotFound` ke sentinel error domain. `FindAll` pakai `pagination.Query[T]` + **WAJIB** whitelist `pagination.SortMap`/`OrderClause` sebelum `.Order(...)` — jangan pernah teruskan `sort` mentah dari client ke GORM (SQL injection, pagination-convention.md §3).
   - `models/<domain_name>_model.go` (package `models`): Model GORM + mapper `ToDomain()` / `FromDomain()`. `ToDomain()` **merekonstruksi struct langsung** (tidak lewat constructor, agar tidak generate UUID baru).

4. **Transport Layer** (`internal/<domain_name>/transport/http/`, package `http`):
   - `handler.go`: HTTP handler (Fiber v3). 5xx JANGAN bocorkan `err.Error()` mentah.
   - `router.go`: register routes.

5. **Dependency Injection (Google Wire)** (`internal/di/`):
   - Tambahkan struct Handler ke `APIHandlers` di `internal/di/api.go` dan panggil `RegisterRoutes`.
   - Tambahkan constructor Repository, Service, Handler ke `ProviderSet` di `internal/di/wire.go`. **Gunakan import alias deskriptif** (mis. `<domain>Domain`, `<domain>App`, `<domain>Adapter`, `<domain>HTTP`).
   - Jalankan `go run github.com/google/wire/cmd/wire@latest ./internal/di` untuk regen `wire_gen.go`.

## Contoh Template

### `internal/<domain_name>/domain/entity.go`
```go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	Err<EntityName>NotFound = errors.New("<domain_name> not found")
)

type <EntityName> struct {
	ID        string
	CompanyID string // scope: kelas Company-owned (scoping-convention.md §1). Tambah BranchID kalau lokasi-spesifik.
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New<EntityName> adalah satu-satunya tempat generate UUID (single source).
// Constructor TIDAK menerima id — id selalu digenerate di sini.
// companyID WAJIB (scope) — entity company-owned tak boleh lahir tanpa scope.
func New<EntityName>(companyID, name string) (*<EntityName>, error) {
	if companyID == "" || name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now()
	return &<EntityName>{
		ID:        uuid.NewString(),
		CompanyID: companyID,
		Name:      name,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
```

### `internal/<domain_name>/domain/repository.go`
```go
package domain

import "context"

// Semua FindXxx/FindAll scope-aware: baca scope.FromContext(ctx) lalu inject
// filter WHERE company_id IN (...) (scoping-convention.md §3). Signature terima ctx —
// scope diambil dari context (diisi RBAC middleware), bukan parameter eksplisit tiap call.
type Repository interface {
	Create(ctx context.Context, item *<EntityName>) error
	FindByID(ctx context.Context, id string) (*<EntityName>, error) // not-found => Err<EntityName>NotFound; wajib dalam scope ctx
	FindAll(ctx context.Context) ([]*<EntityName>, error)           // wajib filter company_id dari scope.FromContext(ctx)
	Update(ctx context.Context, item *<EntityName>) error
	Delete(ctx context.Context, id string) error
}
```

> **Kalau modul punya endpoint List yang perlu paginated** (mayoritas kasus — jangan `Find()` tanpa limit), ganti signature `FindAll` di atas jadi versi paginated berikut (pagination-convention.md §1 — tetap primitif, TANPA import `pkg/pagination` di domain):
> ```go
> FindAll(ctx context.Context, page, limit int, sort, order string) ([]*<EntityName>, int64, error)
> ```

### `internal/<domain_name>/application/dto.go` (List Response, kalau ada pagination)
```go
package application

import "github.com/bagusyanuar/hris-backend/pkg/pagination"

type <EntityName>ListResponse struct {
	Items []<EntityName>Response `json:"items"`
	Meta  pagination.Meta        `json:"meta"` // JANGAN bikin struct meta lokal duplikat
}
```

### `internal/<domain_name>/application/service.go` (List use case, kalau ada pagination)
```go
func (s *Service) List<EntityName>s(ctx context.Context, page, limit int, sort, order string) (*<EntityName>ListResponse, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
	items, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order)
	if err != nil {
		return nil, err
	}
	res := make([]<EntityName>Response, 0, len(items))
	for _, it := range items {
		res = append(res, to<EntityName>Response(it))
	}
	return &<EntityName>ListResponse{Items: res, Meta: pagination.NewMeta(req, total)}, nil
}
```

### `internal/<domain_name>/adapter/models/<domain_name>_model.go`
```go
package models

import (
	"time"

	"github.com/bagusyanuar/hris-backend/internal/<domain_name>/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type <EntityName>Model struct {
	ID        string `gorm:"primaryKey;type:varchar(50)"`
	Name      string `gorm:"type:varchar(100);not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (<EntityName>Model) TableName() string { return "<domain_name>s" }

// BeforeCreate = jaring pengaman UUID (hanya jalan pada Create(), bukan Save()).
func (m *<EntityName>Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

// ToDomain merekonstruksi entity LANGSUNG (tidak lewat constructor,
// supaya tidak generate UUID baru & tidak menjalankan ulang validasi create).
func (m *<EntityName>Model) ToDomain() *domain.<EntityName> {
	return &domain.<EntityName>{
		ID:        m.ID,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func <EntityName>FromDomain(e *domain.<EntityName>) *<EntityName>Model {
	return &<EntityName>Model{
		ID:        e.ID,
		Name:      e.Name,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
```

### `internal/<domain_name>/adapter/postgres.go`
```go
package adapter

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/<domain_name>/domain"
	"github.com/bagusyanuar/hris-backend/internal/<domain_name>/adapter/models"
	"gorm.io/gorm"
)

type <EntityName>Repository struct {
	db *gorm.DB
}

func New<EntityName>Repository(db *gorm.DB) domain.Repository {
	return &<EntityName>Repository{db: db}
}

// Create = INSERT baru => WAJIB Create(), BUKAN Save() (lihat persistence-convention.md §1).
func (r *<EntityName>Repository) Create(ctx context.Context, item *domain.<EntityName>) error {
	model := models.<EntityName>FromDomain(item)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *<EntityName>Repository) FindByID(ctx context.Context, id string) (*domain.<EntityName>, error) {
	var model models.<EntityName>Model
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.Err<EntityName>NotFound // sentinel, JANGAN return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *<EntityName>Repository) FindAll(ctx context.Context) ([]*domain.<EntityName>, error) {
	var rows []models.<EntityName>Model
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.<EntityName>, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, nil
}

// Update = record PK sudah pasti ada => Save()/Updates() aman untuk semantik update.
func (r *<EntityName>Repository) Update(ctx context.Context, item *domain.<EntityName>) error {
	model := models.<EntityName>FromDomain(item)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *<EntityName>Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.<EntityName>Model{}, "id = ?", id).Error
}
```

> **Versi paginated** (kalau `FindAll` di domain pakai signature `page, limit int, sort, order string`, pagination-convention.md §3) — ganti method `FindAll` di atas jadi ini. WAJIB whitelist `SortMap`, JANGAN pass `sort` mentah ke `.Order()`:
> ```go
> var <domainName>SortMap = pagination.SortMap{
> 	"name":       "name",
> 	"created_at": "created_at",
> }
>
> func (r *<EntityName>Repository) FindAll(ctx context.Context, page, limit int, sort, order string) ([]*domain.<EntityName>, int64, error) {
> 	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}
> 	db := r.db.WithContext(ctx).Order(req.OrderClause(<domainName>SortMap, "created_at"))
> 	rows, meta, err := pagination.Query[models.<EntityName>Model](db, req)
> 	if err != nil {
> 		return nil, 0, err
> 	}
> 	result := make([]*domain.<EntityName>, 0, len(rows))
> 	for i := range rows {
> 		result = append(result, rows[i].ToDomain())
> 	}
> 	return result, meta.Total, nil
> }
> ```

### `internal/<domain_name>/transport/http/handler.go`
```go
package http

import (
	"github.com/bagusyanuar/hris-backend/internal/<domain_name>/application"
	"github.com/bagusyanuar/hris-backend/pkg/response"
	"github.com/bagusyanuar/hris-backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// PERHATIAN: Untuk validasi request, JANGAN inject validator ke struct Handler.
// Gunakan fungsi global pkg/validator langsung: `errs := validator.ValidateStruct(req)`.
// Untuk 5xx, JANGAN kirim err.Error() mentah ke client — log lalu balikan pesan generik.

// WAJIB KONSISTEN! Implementasikan HANYA DAN SEMUA operasi yang disetujui di Tech Specs.
// DILARANG MEMOTONG KODE ATAU MENGGUNAKAN PLACEHOLDER. TULIS SEMUA METHOD LENGKAP.
func (h *Handler) Create(c fiber.Ctx) error { /* full impl */ }
func (h *Handler) Get(c fiber.Ctx) error    { /* full impl */ }
// ... sertakan method lain (FindAll, Update, Delete) HANYA jika diwajibkan spesifikasi.
```

### Import alias di `internal/di/wire.go`
```go
import (
	<domain>Domain  "github.com/bagusyanuar/hris-backend/internal/<domain_name>/domain"
	<domain>App     "github.com/bagusyanuar/hris-backend/internal/<domain_name>/application"
	<domain>Adapter "github.com/bagusyanuar/hris-backend/internal/<domain_name>/adapter"
	<domain>HTTP    "github.com/bagusyanuar/hris-backend/internal/<domain_name>/transport/http"
)
```
