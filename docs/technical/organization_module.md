# Dokumentasi Teknis Modul Organisasi (Enterprise-Grade HRIS)

Dokumen ini mendeskripsikan implementasi teknis untuk modul `Organization` di dalam HRIS Backend, yang menangani struktur perusahaan menggunakan pendekatan 3-Entitas: `Department`, `JobTitle`, dan `JobPosition`.

## Struktur Entitas Database (GORM Models)

Modul ini diimplementasikan dalam *Bounded Context* bernama `organization`.

### 1. `Department` (Unit Kerja)
Menyimpan data unit kerja / departemen secara murni.
- `id` (UUID, Primary Key)
- `code` (String, Unique, Index) - misal "IT", "FIN"
- `name` (String) - misal "Information Technology"
- `parent_id` (UUID, Nullable) - untuk sub-departemen
- `is_active` (Boolean, Default: True)
- `created_at`, `updated_at`, `deleted_at`

### 2. `JobTitle` (Titel Pangkat / Grade)
Menyimpan standar pangkat dan jenjang karir (Salary Band).
- `id` (UUID, Primary Key)
- `code` (String, Unique, Index) - misal "MGR", "STF"
- `name` (String) - misal "Manager", "Staff"
- `grade_level` (Integer) - level hierarki (misal 5 untuk manager, 2 untuk staff)
- `is_active` (Boolean, Default: True)
- `created_at`, `updated_at`, `deleted_at`

### 3. `JobPosition` (Slot / Jabatan Aktif)
Tabel *mapping* yang membentuk struktur organisasi (Headcount & Reporting).
- `id` (UUID, Primary Key)
- `department_id` (UUID, Foreign Key ke `departments`)
- `job_title_id` (UUID, Foreign Key ke `job_titles`)
- `name` (String) - Nama spesifik posisi (misal: "IT Manager")
- `reports_to_id` (UUID, Nullable, Foreign Key ke `job_positions`) - Atasan dari posisi ini
- `headcount_quota` (Integer, Default: 1) - Batas maksimal orang di posisi ini
- `is_active` (Boolean, Default: True)
- `created_at`, `updated_at`, `deleted_at`

## API Endpoints (Interfaces)

Semua endpoint diawali dengan `/organization/` di dalam API Gateway/Router.

**Departments:**
- `GET /organization/departments` - Get all departments
- `POST /organization/departments` - Create a new department
- `GET /organization/departments/:id` - Get department by ID
- `PUT /organization/departments/:id` - Update department by ID
- `DELETE /organization/departments/:id` - Delete department by ID

**Job Titles:**
- `GET /organization/job-titles` - Get all job titles
- `POST /organization/job-titles` - Create a new job title
- `GET /organization/job-titles/:id` - Get job title by ID
- `PUT /organization/job-titles/:id` - Update job title by ID
- `DELETE /organization/job-titles/:id` - Delete job title by ID

**Job Positions:**
- `GET /organization/job-positions` - Get all job positions
- `POST /organization/job-positions` - Create a new job position
- `GET /organization/job-positions/:id` - Get job position by ID
- `PUT /organization/job-positions/:id` - Update job position by ID
- `DELETE /organization/job-positions/:id` - Delete job position by ID

## Struktur Folder dan Arsitektur (DDD)

Implementasi ini secara ketat mematuhi arsitektur **Domain-Driven Design (DDD)**:

### Domain Layer (Pusat Bisnis)
Berisi pure business logic, validasi entity, dan interface repository.
- `internal/domain/organization/entity.go`
- `internal/domain/organization/repository.go`

### Application Layer (Use Case)
Berisi service untuk mengkoordinasi transaksi, memanggil repo, serta DTO mapping.
- `internal/application/organization/service.go`
- `internal/application/organization/dto.go`

### Infrastructure Layer (Database & ORM)
Implementasi dari interface repository menggunakan Postgres & GORM.
- `internal/infrastructure/repository/organization_postgres.go`
- `internal/infrastructure/repository/models/organization_model.go`

### Interfaces Layer (HTTP & Router)
Menangani request/response HTTP menggunakan framework Fiber v3.
- `internal/interfaces/http/organization/handler.go`
- `internal/interfaces/http/organization/router.go`

## Dokumentasi API (Swagger & Bruno)

Dokumentasi lengkap mengenai request body, response, dan struktur JSON dapat ditemukan di:
- **Bruno Collection**: `/docs/api/bruno/Organization/...`
- **Swagger OpenAPI**: Belum ada YAML, dapat di generate menyusul sesuai aturan DDD project.
