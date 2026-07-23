# Konvensi Pagination, Sort & Order (`pkg/pagination`)

Setiap endpoint `List`/`FindAll` yang berpotensi mengembalikan banyak row **WAJIB** paginated — jangan `Find()` tanpa `LIMIT` (full scan tak terbatas). Semua modul **WAJIB** pakai `pkg/pagination`, jangan reinvent logic default/offset/sort sendiri per domain (lihat implementasi rujukan di `internal/organization/`).

---

## 1. Pembagian Tanggung Jawab per Layer

`pkg/pagination` sengaja dipecah 2 file dalam 1 package supaya batas tanggung jawab jelas:

| File | Isi | Boleh diimport dari layer |
|------|-----|---------------------------|
| `pkg/pagination/pagination.go` | `Request{Page,Limit,Sort,Order}`, `Normalize()`, `Offset()`, `SortMap`, `OrderClause()`, `Meta`, `NewMeta()` — **pure Go, tanpa dependency GORM** | Application, Transport (bebas) |
| `pkg/pagination/gorm.go` | `Query[T any](db *gorm.DB, req Request) ([]T, Meta, error)` — generic Count+Offset+Limit+Find | **HANYA Adapter layer** (repository) |

**Domain Layer TETAP TIDAK BOLEH mengimport `pkg/pagination` sama sekali** (selaras architecture.md — domain pure Go, tanpa dependency infra). Signature repository interface di domain tetap pakai primitif:
```go
// domain/repository.go
FindAll(ctx context.Context, page, limit int, sort, order string) ([]*Entity, int64, error)
```
Adapter yang membungkus primitif itu jadi `pagination.Request` di dalam implementasinya.

## 2. Application Layer — Normalize Page/Limit/Order

Application Service menerima `page, limit int, sort, order string` mentah dari Transport, bungkus jadi `pagination.Request` dan panggil `.Normalize()` SEBELUM diteruskan ke repository:
```go
req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
items, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order)
...
return &ListResponse{Items: items, Meta: pagination.NewMeta(req, total)}, nil
```
`Normalize()` mengisi default `page=1`, `limit=20` (cap `MaxLimit=100`), dan `order` jadi `"asc"`/`"desc"` valid. `Sort` (logical key) **TIDAK** divalidasi di layer ini — whitelist-nya baru terjadi di Adapter (§3), karena kolom sortable beda tiap entity.

DTO List Response **WAJIB** pakai `pagination.Meta` untuk field `meta`, jangan bikin struct `ListMeta` lokal duplikat per domain:
```go
type CompanyListResponse struct {
    Items []CompanyResponse `json:"items"`
    Meta  pagination.Meta   `json:"meta"`
}
```

## 3. Adapter Layer — Whitelist Sort (WAJIB, Security)

**DILARANG KERAS** meneruskan `Request.Sort` mentah dari client ke GORM `.Order(...)` — itu SQL injection vector (`Order()` menerima raw clause string, bukan cuma nama kolom). **WAJIB** definisikan `pagination.SortMap` (logical key client → kolom DB asli) per entity, lalu resolve lewat `OrderClause`:

```go
var companySortMap = pagination.SortMap{
    "code":       "code",
    "legal_name": "legal_name",
    "created_at": "created_at",
    "updated_at": "updated_at",
}

func (r *companyRepository) FindAll(ctx context.Context, page, limit int, sort, order string) ([]*domain.Company, int64, error) {
    req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}
    db := dbFromContext(ctx, r.db).Order(req.OrderClause(companySortMap, "created_at"))
    rows, meta, err := pagination.Query[models.CompanyModel](db, req)
    ...
}
```
`OrderClause(allowed, defaultSort)` fallback ke `defaultSort` kalau `Sort` client gak ada di whitelist — silent fallback, BUKAN error 422 (dokumentasikan perilaku ini di Swagger/Bruno, lihat api-documentation.md).

`pagination.Query[T]` menjalankan `Count` + `Offset().Limit().Find()` di atas `*gorm.DB` yang sudah di-scope caller (`.Where(...)` dsb sebelum dipanggil) — filter scope (`company_id`, dst, scoping-convention.md §3) tetap di-chain di `db` SEBELUM masuk ke `Query[T]`.

## 4. Transport Layer — Parsing Query Param

Handler parse `page`, `limit`, `sort`, `order` dari query string dan teruskan APA ADANYA (tanpa validasi/default di sini — itu tanggung jawab Application §2):
```go
func parsePagination(c fiber.Ctx) (page, limit int, sort, order string) {
    page, _ = strconv.Atoi(c.Query("page"))
    limit, _ = strconv.Atoi(c.Query("limit"))
    sort = c.Query("sort")
    order = c.Query("order")
    return
}
```

## 5. Dokumentasi API (selaras api-documentation.md)

Endpoint List **WAJIB** mendokumentasikan 4 query param (`page`, `limit`, `sort`, `order`) di Swagger (`enum` buat `sort` sesuai `SortMap` yang tersedia, `enum: [asc, desc]` buat `order`) dan Bruno. Catat eksplisit di docs bahwa `sort` invalid **fallback diam-diam ke default**, bukan `422` — supaya FE gak salah asumsi.

## 6. Checklist Review

- [ ] Endpoint List/FindAll pakai `pkg/pagination`, bukan reinvent `page/limit` manual.
- [ ] Domain layer (`internal/<domain>/domain/`) TIDAK mengimport `pkg/pagination` — signature repository tetap primitif (`page, limit int, sort, order string`).
- [ ] Adapter TIDAK PERNAH pass `Request.Sort` mentah ke `.Order()` — WAJIB lewat `SortMap` + `OrderClause`.
- [ ] `pagination.Query[T]` hanya dipanggil di Adapter layer, bukan Application.
- [ ] DTO List Response pakai `pagination.Meta`, tidak ada `ListMeta` duplikat lokal.
- [ ] Swagger + Bruno endpoint List mendokumentasikan `page`/`limit`/`sort`/`order` + catatan fallback silent untuk `sort` invalid.
