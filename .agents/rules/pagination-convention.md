# Konvensi Pagination, Sort, Order & Search (`pkg/pagination`)

Setiap endpoint `List`/`FindAll` yang berpotensi mengembalikan banyak row **WAJIB** paginated — jangan `Find()` tanpa `LIMIT` (full scan tak terbatas). Semua modul **WAJIB** pakai `pkg/pagination`, jangan reinvent logic default/offset/sort/search sendiri per domain (lihat implementasi rujukan di `internal/organization/`, `internal/workforce/`).

**`search` adalah query param standar** (sejajar `page`/`limit`/`sort`/`order`) — **WAJIB** ada di tiap endpoint List **kalau entity-nya punya minimal satu kolom text human-readable** (`name`, `code`, `title`, dst). Kalau entity gak punya kolom text sama sekali (mis. cuma numeric/enum/tanggal), `search` boleh di-skip — dokumentasikan alasannya eksplisit di tech-spec (pola sama kayak opt-out "Global master" di scoping-convention.md §1), bukan default diam-diam dihilangkan.

---

## 1. Pembagian Tanggung Jawab per Layer

`pkg/pagination` sengaja dipecah 2 file dalam 1 package supaya batas tanggung jawab jelas:

| File | Isi | Boleh diimport dari layer |
|------|-----|---------------------------|
| `pkg/pagination/pagination.go` | `Request{Page,Limit,Sort,Order,Search}`, `Normalize()`, `Offset()`, `SortMap`, `OrderClause()`, `SearchClause()`, `Meta`, `NewMeta()` — **pure Go, tanpa dependency GORM** | Application, Transport (bebas) |
| `pkg/pagination/gorm.go` | `Query[T any](db *gorm.DB, req Request) ([]T, Meta, error)` — generic Count+Offset+Limit+Find | **HANYA Adapter layer** (repository) |

**Domain Layer TETAP TIDAK BOLEH mengimport `pkg/pagination` sama sekali** (selaras architecture.md — domain pure Go, tanpa dependency infra). Signature repository interface di domain tetap pakai primitif — `search` ikut sebagai string biasa kalau entity-nya searchable (§ intro):
```go
// domain/repository.go
FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*Entity, int64, error)
```
Adapter yang membungkus primitif itu jadi `pagination.Request` di dalam implementasinya.

## 2. Application Layer — Normalize Page/Limit/Order

Application Service menerima `page, limit int, sort, order, search string` mentah dari Transport, bungkus jadi `pagination.Request` dan panggil `.Normalize()` SEBELUM diteruskan ke repository:
```go
req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}.Normalize()
items, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Sort, req.Order, search)
...
return &ListResponse{Items: items, Meta: pagination.NewMeta(req, total)}, nil
```
`Normalize()` mengisi default `page=1`, `limit=20` (cap `MaxLimit=100`), dan `order` jadi `"asc"`/`"desc"` valid. `Sort` (logical key) **TIDAK** divalidasi di layer ini — whitelist-nya baru terjadi di Adapter (§3), karena kolom sortable beda tiap entity. `search` **TIDAK** perlu `Normalize()` (gak ada default/cap — kosong = tanpa filter, sah apa adanya) — cukup diteruskan sebagai parameter primitif terpisah ke repository, sejajar `sort`/`order`.

DTO List Response **WAJIB** pakai `pagination.Meta` untuk field `Meta` — struct `{Items, Meta}` ini adalah *return type* method Application Service, BUKAN bentuk envelope HTTP akhir (lihat pemetaan ke `response.SuccessList` di bawah). Jangan bikin struct `ListMeta` lokal duplikat per domain:
```go
type CompanyListResponse struct {
    Items []CompanyResponse `json:"items"`
    Meta  pagination.Meta   `json:"meta"`
}
```

**Transport layer WAJIB serialize lewat `response.SuccessList(c, code, message, res.Items, res.Meta)`**, bukan `response.Success(c, code, message, res)`. `SuccessList` membongkar DTO `{Items, Meta}` di atas jadi envelope HTTP:
```go
return response.SuccessList(c, fiber.StatusOK, "Companies fetched successfully", res.Items, res.Meta)
```
menghasilkan:
```json
{
  "code": 200, "status": "success", "message": "Companies fetched successfully",
  "data": [ { "id": "...", "name": "..." } ],
  "meta": { "pagination": { "page": 1, "limit": 20, "total": 57, "total_pages": 3 } }
}
```
`data` adalah array item **langsung** (tanpa wrapper `{items, meta}`), dan `meta` jadi sibling top-level dari `data` — **BUKAN** `data.meta`. `meta.pagination` (bukan flat `meta.page`/`meta.limit`) sengaja dikelompokkan supaya metadata lain (request tracing, summary count) bisa ditambah ke `meta` di masa depan tanpa nabrak field pagination. Endpoint single-record (Get/Create/Update) TETAP pakai `response.Success` biasa — `data` object, tanpa key `meta` sama sekali (bukan `{}`).

## 3. Adapter Layer — Whitelist Sort & Search (WAJIB, Security)

**DILARANG KERAS** meneruskan `Request.Sort` mentah dari client ke GORM `.Order(...)` — itu SQL injection vector (`Order()` menerima raw clause string, bukan cuma nama kolom). **WAJIB** definisikan `pagination.SortMap` (logical key client → kolom DB asli) per entity, lalu resolve lewat `OrderClause`. `search` pakai prinsip whitelist yang sama lewat `SearchClause(columns ...string)`:

```go
var departmentSortMap = pagination.SortMap{
    "code":       "code",
    "name":       "name",
    "created_at": "created_at",
    "updated_at": "updated_at",
}

func (r *departmentRepository) FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*domain.Department, int64, error) {
    req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order, Search: search}
    db := dbFromContext(ctx, r.db).Order(req.OrderClause(departmentSortMap, "created_at"))
    if clause, args := req.SearchClause("code", "name"); clause != "" {
        db = db.Where(clause, args...)
    }
    rows, meta, err := pagination.Query[models.DepartmentModel](db, req)
    ...
}
```
`OrderClause(allowed, defaultSort)` fallback ke `defaultSort` kalau `Sort` client gak ada di whitelist — silent fallback, BUKAN error 422 (dokumentasikan perilaku ini di Swagger/Bruno, lihat api-documentation.md).

`SearchClause(columns ...string)` mengembalikan `("(col1 ILIKE ? OR col2 ILIKE ? ...)", args)` — kolom yang di-search **WAJIB** whitelist eksplisit sebagai argumen (persis alasan `SortMap`: jangan pernah biarkan nama kolom mentah dari client nyampe ke query string). Kalau `Search` kosong, `SearchClause` balikin `("", nil)` — cek `clause != ""` sebelum `.Where(...)`, biar gak nambah `WHERE ("")` kosong.

**Pengecualian search lintas-table** (mis. `Company` match `legal_name` miliknya sendiri ATAU nama `Branch` anaknya via `EXISTS` subquery — lihat `internal/organization/adapter/postgres.go`): `SearchClause` cuma cover same-table. Kasus lintas-table **BOLEH** ditulis manual (raw `.Where(...)` dengan subquery), itu bukan pelanggaran — dokumentasikan di tech-spec kenapa gak pakai `SearchClause` polos.

`pagination.Query[T]` menjalankan `Count` + `Offset().Limit().Find()` di atas `*gorm.DB` yang sudah di-scope caller (`.Where(...)` dsb sebelum dipanggil) — filter scope (`company_id`, dst, scoping-convention.md §3) tetap di-chain di `db` SEBELUM masuk ke `Query[T]`.

## 4. Transport Layer — Parsing Query Param

Handler parse `page`, `limit`, `sort`, `order`, `search` dari query string dan teruskan APA ADANYA (tanpa validasi/default di sini — itu tanggung jawab Application §2):
```go
func parsePagination(c fiber.Ctx) (page, limit int, sort, order string) {
    page, _ = strconv.Atoi(c.Query("page"))
    limit, _ = strconv.Atoi(c.Query("limit"))
    sort = c.Query("sort")
    order = c.Query("order")
    return
}
```
`search` di-parse terpisah (`search := c.Query("search")`) — bukan bagian `parsePagination` biasa, karena gak semua endpoint List punya `search` (§ intro), jadi jangan dipaksa satu signature buat semuanya.

## 5. Dokumentasi API (selaras api-documentation.md)

Endpoint List **WAJIB** mendokumentasikan `page`, `limit`, `sort`, `order` di Swagger (`enum` buat `sort` sesuai `SortMap` yang tersedia, `enum: [asc, desc]` buat `order`) dan Bruno. Catat eksplisit di docs bahwa `sort` invalid **fallback diam-diam ke default**, bukan `422` — supaya FE gak salah asumsi.

Skema response 200 di Swagger **WAJIB** merefleksikan envelope `response.SuccessList` yang sebenarnya: `data` bertipe `array` (`items: $ref: '#/components/schemas/<Entity>Response'`), dan `meta.pagination` (`{page, limit, total, total_pages}`) sebagai property top-level `meta` di object response — **BUKAN** `data.items`/`data.meta`. Bruno `docs { }` contoh JSON-nya harus sinkron persis (lihat api-documentation.md untuk aturan schema reusable).

Kalau entity-nya searchable (§ intro), `search` **WAJIB** ikut didokumentasikan sejajar 4 param di atas — sebutkan eksplisit kolom apa aja yang di-match (mis. "match `code` ATAU `name`") dan bahwa kosong = tanpa filter. Kalau entity TIDAK searchable, gak perlu nyantumin `search` sama sekali (jangan dokumentasikan param yang gak ada).

## 6. Checklist Review

- [ ] Endpoint List/FindAll pakai `pkg/pagination`, bukan reinvent `page/limit` manual.
- [ ] Domain layer (`internal/<domain>/domain/`) TIDAK mengimport `pkg/pagination` — signature repository tetap primitif (`page, limit int, sort, order, search string`).
- [ ] Adapter TIDAK PERNAH pass `Request.Sort` mentah ke `.Order()` — WAJIB lewat `SortMap` + `OrderClause`.
- [ ] Entity punya kolom text human-readable (`name`/`code`/dst) → endpoint List-nya WAJIB `search`, lewat `Request.SearchClause(columns...)` dengan whitelist eksplisit (kecuali kasus lintas-table, §3 pengecualian). Entity TANPA kolom text boleh skip, tapi dijustifikasi di tech-spec.
- [ ] `pagination.Query[T]` hanya dipanggil di Adapter layer, bukan Application.
- [ ] DTO List Response (application layer) pakai `pagination.Meta`, tidak ada `ListMeta` duplikat lokal.
- [ ] Handler Transport layer serialize lewat `response.SuccessList(c, code, message, res.Items, res.Meta)` — bukan `response.Success` — supaya `data` jadi array langsung dan `meta.pagination` jadi sibling top-level, bukan `data.meta`.
- [ ] Swagger + Bruno endpoint List mendokumentasikan `page`/`limit`/`sort`/`order` (+ `search` kalau applicable) + catatan fallback silent untuk `sort` invalid, DAN skema response-nya `data: array` + `meta.pagination` (bukan `data.items`/`data.meta`).
