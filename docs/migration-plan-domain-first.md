# Migration Plan — Layer-First → Domain-First

**Tanggal:** 2026-07-22
**Status:** PLAN (belum dieksekusi)
**Tujuan:** Restrukturisasi `internal/` dari *package-by-layer* ke *package-by-bounded-context* (domain-first), agar tiap bounded context berdiri sebagai satu unit utuh dan siap diekstrak menjadi service terpisah. Dependency rule Clean Architecture **tidak berubah** — hanya pengelompokan folder.

> Refactor ini **murni mekanis** (pindah file + rewrite import + regen wire). Tidak ada perubahan logika bisnis.

---

## 1. Keputusan yang Sudah Disepakati

- Pola target: **domain-first** dengan sub-folder layer di dalam tiap context.
- Cross-cutting (`config`, `database`, `middleware`) → `internal/shared/`.
- `pkg/response` + `pkg/validator` → **TETAP di `pkg/`** (konvensi Go untuk util reusable).
- `interfaces/` di-rename → `transport/` (hindari bentrok makna dengan keyword `interface` Go).
- `infrastructure/security/jwt.go` → milik domain `auth` → `internal/auth/infrastructure/`, **bukan** shared.
- `internal/di/` tetap top-level (wiring perlu tahu semua context).

---

## 2. Target Layout

```
internal/
  employee/
    domain/          entity.go  repository.go
    application/     service.go  dto.go
    infrastructure/  postgres.go  models/…
    transport/http/  handler.go  router.go
  organization/
    domain/ application/ infrastructure/ transport/http/
  auth/
    domain/          token_generator.go
    application/     service.go
    infrastructure/  jwt.go
    transport/http/  handler.go  router.go
  user/                              # domain internal, dikonsumsi auth
    domain/          entity.go  repository.go
    infrastructure/  postgres.go  models/…
  shared/
    config/          config.go
    database/        postgres.go
    middleware/      auth_middleware.go
  di/                wire.go  wire_gen.go  api.go
pkg/                 response/  validator/     # TIDAK dipindah
cmd/                 api/  seed/               # path tidak berubah
```

---

## 3. Mapping File (Lama → Baru)

### Employee
| Lama | Baru |
|------|------|
| `internal/domain/employee/entity.go` | `internal/employee/domain/entity.go` |
| `internal/domain/employee/repository.go` | `internal/employee/domain/repository.go` |
| `internal/application/employee/service.go` | `internal/employee/application/service.go` |
| `internal/application/employee/dto.go` | `internal/employee/application/dto.go` |
| `internal/infrastructure/repository/employee_postgres.go` | `internal/employee/infrastructure/postgres.go` |
| `internal/infrastructure/repository/models/employee_model.go` | `internal/employee/infrastructure/models/employee_model.go` |
| `internal/interfaces/http/employee/handler.go` | `internal/employee/transport/http/handler.go` |
| `internal/interfaces/http/employee/router.go` | `internal/employee/transport/http/router.go` |

### Organization
| Lama | Baru |
|------|------|
| `internal/domain/organization/entity.go` | `internal/organization/domain/entity.go` |
| `internal/domain/organization/repository.go` | `internal/organization/domain/repository.go` |
| `internal/application/organization/service.go` | `internal/organization/application/service.go` |
| `internal/application/organization/dto.go` | `internal/organization/application/dto.go` |
| `internal/infrastructure/repository/organization_postgres.go` | `internal/organization/infrastructure/postgres.go` |
| `internal/infrastructure/repository/models/organization_model.go` | `internal/organization/infrastructure/models/organization_model.go` |
| `internal/interfaces/http/organization/handler.go` | `internal/organization/transport/http/handler.go` |
| `internal/interfaces/http/organization/router.go` | `internal/organization/transport/http/router.go` |

### Auth
| Lama | Baru |
|------|------|
| `internal/domain/auth/token_generator.go` | `internal/auth/domain/token_generator.go` |
| `internal/application/auth/service.go` | `internal/auth/application/service.go` |
| `internal/infrastructure/security/jwt.go` | `internal/auth/infrastructure/jwt.go` |
| `internal/interfaces/http/auth/handler.go` | `internal/auth/transport/http/handler.go` |
| `internal/interfaces/http/auth/router.go` | `internal/auth/transport/http/router.go` |

### User (domain internal — hanya domain + infrastructure)
| Lama | Baru |
|------|------|
| `internal/domain/user/entity.go` | `internal/user/domain/entity.go` |
| `internal/domain/user/repository.go` | `internal/user/domain/repository.go` |
| `internal/infrastructure/repository/user_postgres.go` | `internal/user/infrastructure/postgres.go` |
| `internal/infrastructure/repository/models/user_model.go` | `internal/user/infrastructure/models/user_model.go` |

### Shared
| Lama | Baru |
|------|------|
| `internal/infrastructure/config/config.go` | `internal/shared/config/config.go` |
| `internal/infrastructure/database/postgres.go` | `internal/shared/database/postgres.go` |
| `internal/interfaces/http/middleware/auth_middleware.go` | `internal/shared/middleware/auth_middleware.go` |

### Tidak dipindah
- `internal/di/*` (hanya update import path di dalamnya)
- `pkg/response/*`, `pkg/validator/*`
- `cmd/api/*`, `cmd/seed/*` (hanya update import path)

---

## 4. Konvensi Package Name (Penting)

Target memakai **package name = nama layer** agar konsisten & jelas di dalam satu context:

| Folder | Package name |
|--------|-------------|
| `internal/<ctx>/domain` | `domain` |
| `internal/<ctx>/application` | `application` |
| `internal/<ctx>/infrastructure` | `infrastructure` |
| `internal/<ctx>/infrastructure/models` | `models` |
| `internal/<ctx>/transport/http` | `http` |

Karena banyak context berbagi nama package sama (`domain`, `application`, dst.), **di layer `di/` dan `cmd/` wajib pakai import alias** yang deskriptif. Contoh:

```go
import (
    empDomain "github.com/bagusyanuar/hris-backend/internal/employee/domain"
    empApp    "github.com/bagusyanuar/hris-backend/internal/employee/application"
    empInfra  "github.com/bagusyanuar/hris-backend/internal/employee/infrastructure"
    empHTTP   "github.com/bagusyanuar/hris-backend/internal/employee/transport/http"

    orgApp    "github.com/bagusyanuar/hris-backend/internal/organization/application"
    // dst.
)
```

> Aturan alias ini menggantikan pola auto-generated Wire (`employee2`, `auth3`) yang sekarang membingungkan. Setelah migrasi, `wire_gen.go` akan memakai alias yang lebih rapi (regen otomatis dari `wire.go`).

Di dalam satu context (mis. `employee/transport/http/handler.go` meng-import `employee/application`), referensinya cukup `application.Service` tanpa perlu alias `appEmployee` yang panjang — lebih bersih dari sekarang.

---

## 5. Urutan Eksekusi (Step-by-Step)

Kerjakan **per context, satu per satu**, jalankan `go build ./...` setelah tiap langkah supaya kesalahan cepat terlokalisir.

1. **shared dulu** — pindah `config`, `database`, `middleware` ke `internal/shared/`. Update import di `cmd/`. Build.
2. **user** — pindah domain + infrastructure. Update import di auth & di. Build.
3. **auth** — pindah domain/application/infrastructure(jwt)/transport. Build.
4. **organization** — pindah 4 layer. Build.
5. **employee** — pindah 4 layer. Build.
6. **Hapus folder lama** yang sudah kosong: `internal/domain`, `internal/application`, `internal/infrastructure`, `internal/interfaces`.
7. **Regen Wire**: `go generate ./internal/di/...` (atau `go run github.com/google/wire/cmd/wire ./internal/di`). Update `wire.go` ProviderSet import path lebih dulu.
8. **Build final**: `go build ./...`.
9. **Lint**: `golangci-lint run`.
10. **Smoke test**: `make run`, cek `GET /` + 1 endpoint per domain.

---

## 6. Dampak di Luar `internal/`

- **`cmd/api/server.go`** — import `config`, `database`(via main), `security.NewJWTService` → `auth/infrastructure`, `middleware` → `shared/middleware`, `di`.
- **`cmd/api/main.go`** — import `config` → `shared/config`, `database` → `shared/database`.
- **`cmd/seed/main.go`** — cek import (kemungkinan pakai `config`, `database`, `user`).
- **`internal/di/{wire.go, api.go}`** — update semua import path + alias, lalu regen `wire_gen.go`.

---

## 7. Dokumentasi yang Wajib Diupdate Setelah Migrasi

- [`architecture.md`](../.agents/rules/architecture.md) — **diagram folder tree** + contoh path (`employee_postgres.go`, dll.) harus disesuaikan ke layout baru.
- [`persistence-convention.md`](../.agents/rules/persistence-convention.md) & [`uuid-generation.md`](../.agents/rules/uuid-generation.md) — cek referensi path file (`models.go`, dll.).
- Skill `scaffold-domain` & workflow `execute-domain` — template scaffold-nya masih menghasilkan layout layer-first; **wajib disesuaikan** agar domain baru langsung lahir domain-first.
- [`architecture-review.md`](architecture-review.md) — tambahkan catatan bahwa struktur sudah dimigrasi.

---

## 8. Risiko & Rollback

- **Risiko:** rendah secara logika (tidak ada perubahan behavior), tapi churn import besar. Risiko utama = ada import path terlewat → build merah (langsung ketahuan, bukan bug senyap).
- **Mitigasi:** eksekusi per-context + `go build ./...` tiap langkah. Commit atomik per context (`refactor(structure): migrate <ctx> to domain-first`).
- **Rollback:** karena murni move + rename, `git revert`/`git reset` aman. Kerjakan di branch terpisah (mis. `refactor/domain-first`), jangan langsung `main`.

---

## 9. Checklist Ringkas

- [ ] Branch `refactor/domain-first` dibuat.
- [ ] shared/ (config, database, middleware) dipindah + build hijau.
- [ ] user, auth, organization, employee dipindah berurutan + build hijau tiap langkah.
- [ ] Folder lama (`domain`, `application`, `infrastructure`, `interfaces`) dihapus.
- [ ] `wire.go` diupdate + `wire_gen.go` diregen.
- [ ] `go build ./...` + `golangci-lint run` hijau.
- [ ] Smoke test endpoint per domain OK.
- [ ] `architecture.md` + scaffold-domain/execute-domain disesuaikan.
- [ ] Commit atomik per context.
