# Architecture Review — HRIS Backend

**Tanggal:** 2026-07-22
**Scope:** Clean Architecture & DDD assessment
**Tech stack:** Go 1.26 · Fiber v3 · GORM · google/wire · PostgreSQL · golang-migrate

---

## 1. Ringkasan Eksekutif

Arsitektur secara **struktural sudah solid dan mengikuti Clean Architecture + DDD dengan disiplin yang jarang ditemukan di tahap awal project**. Pemisahan layer (`domain` / `application` / `infrastructure` / `interfaces`) konsisten, dependency rule dihormati, dan domain layer benar-benar murni (tanpa GORM tags, tanpa import framework).

Namun ada beberapa masalah **korektnes berisiko tinggi** (terutama pola `GORM Save()` pada Progressive Save) dan **kesenjangan enterprise-readiness** (nol test, tidak ada panic recovery, tidak ada structured logging, transaction boundary belum di layer yang benar) yang **wajib** dibereskan sebelum modul HRIS berikutnya (Attendance, Payroll) dibangun di atas fondasi ini.

**Verdict:** Fondasi bagus (7.5/10 struktur), tapi belum production-ready (4/10 operasional). Refactor terarah bisa naikkan ke enterprise-grade tanpa merombak arsitektur.

| Aspek | Nilai | Catatan |
|-------|-------|---------|
| Layering & dependency rule | 🟢 Kuat | Domain murni, arah dependensi benar |
| Domain modeling (DDD) | 🟡 Cukup | Aggregate boundary bocor via Progressive Save |
| Correctness / data integrity | 🔴 Berisiko | Pola `Save()` berpotensi silent no-op (lihat 3.1) |
| Error handling | 🟡 Cukup | Mapping benar, tapi bocorkan internal error ke client |
| Observability & ops | 🔴 Kurang | Tanpa recover/logger/CORS terpasang |
| Testability | 🔴 Kurang | 0 test, padahal interface sudah mock-friendly |
| Consistency | 🟡 Cukup | UUID & pola berbeda antar domain |

---

## 2. Yang Sudah Benar (Pertahankan)

- **Domain purity** — `internal/domain/employee/entity.go` bebas GORM/JSON tag. Model DB terpisah di `infrastructure/repository/models/` lengkap dengan mapper `ToDomain()` / `FromDomain()`. Ini tepat sesuai [architecture.md](../.agents/rules/architecture.md).
- **Dependency direction benar** — repository interface didefinisikan di domain (`domain/employee/repository.go`), diimplementasikan di infrastructure. Application depend ke abstraksi, bukan ke GORM.
- **Constructor invariants** — `NewEmployee`, `NewBank`, dll. memvalidasi field wajib dan mengisi UUID. Rule bisnis (`ErrPrimaryBankRequired`) hidup di service/domain, bukan di handler.
- **Custom domain errors** dipetakan ke HTTP status yang tepat di handler (404/409/400) — persis pola yang diinginkan rule coding-convention.
- **Wire DI bersih**, split by ProviderSet (Repository/Service/Handler). `server.go` cukup panggil `di.InitializeAPI`.
- **Migrasi via golang-migrate**, bukan `AutoMigrate` — sesuai rule.
- **Graceful shutdown** dengan signal handling + close DB pool di `cmd/api/server.go`.
- **Response envelope** standar (`pkg/response`) konsisten dipakai handler.

---

## 3. Temuan Kritis (Wajib Fix)

### 3.1 🔴 CRITICAL — `GORM Save()` pada Progressive Save berpotensi silent no-op / data loss

**Lokasi:** `internal/infrastructure/repository/employee_postgres.go` (`SaveCore`, `SavePersonalData`, `SaveContact`, `SaveDocument`).

**Masalah:** Domain constructor **selalu** generate UUID baru setiap dipanggil (`NewPersonalData` → `uuid.NewString()`). Repo lalu memanggil `db.Save(model)`.

GORM `Save()` berperilaku: kalau primary key **terisi**, ia jalankan **UPDATE (semua field)**, bukan INSERT. Karena PK selalu terisi UUID baru yang belum ada di DB:

- **Insert pertama** → `UPDATE ... WHERE id = <uuid-baru>` → 0 rows affected → **tidak error, tapi tidak insert apa pun**.
- **Update kedua** (`UpdatePersonalData` dipanggil ulang) → generate UUID lain lagi → tetap 0 rows → data lama tak tersentuh, data baru tak tersimpan.

Efek: perubahan bisa **hilang diam-diam tanpa error**. Ini bug paling berbahaya karena lolos dari happy-path testing manual sekilas.

> **Verifikasi dulu:** jalankan `UpdatePersonalData` dua kali untuk employee yang sama, cek row di tabel `employee_personal_data`. Kalau hanya 0/1 row atau data tidak berubah pada call kedua → confirmed.

**Fix yang disarankan** (pilih salah satu, konsisten):
1. **Upsert eksplisit by business key.** Untuk entity 1-1 (personal data, contact) yang di-key oleh `employee_id`, pakai:
   ```go
   r.db.WithContext(ctx).
     Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "employee_id"}},
        UpdateAll: true,
     }).Create(model)
   ```
2. Atau **fetch dulu, reuse PK**: repo cari row existing by `employee_id`, kalau ada pakai PK-nya, kalau tidak `Create`. Domain constructor jangan paksa generate UUID saat update.

Prinsip: **jangan pakai `Save()` untuk semantik "insert-or-update by non-PK key"** — itu bukan yang dilakukan `Save`.

### 3.2 🔴 HIGH — Transaction boundary ada di layer yang salah

Rule [architecture.md](../.agents/rules/architecture.md) §B: *Application Layer bertanggung jawab atas Transaction Management.* Sekarang transaksi hidup di dalam **repository** (`SaveBanks`, `SaveEducations` buka `db.Transaction`), sedangkan `UpdatePersonalData` di service melakukan `FindByID` → `FindByKTP` → `SavePersonalData` sebagai **3 operasi tanpa transaksi**.

Akibat:
- **TOCTOU** pada cek duplikat KTP — dua request paralel bisa lolos cek lalu dua-duanya insert (untung ada unique constraint DB sebagai jaring pengaman, tapi error DB-nya **tidak** dipetakan ke `ErrKTPDuplicate` → user dapat 500, bukan 409).
- Tidak ada atomicity lintas sub-entity saat butuh.

**Fix:** perkenalkan **Unit of Work / transaction manager** yang di-inject ke application service. Pola umum di Go:
```go
type TxManager interface {
    Do(ctx context.Context, fn func(ctx context.Context) error) error
}
```
Repo baca `tx` dari context. Service yang membuka & menutup transaksi. Repository kembali murni jadi data-access.

### 3.3 🟠 MEDIUM — Internal error bocor ke client (info leak)

Di semua handler: `return response.Error(c, 500, err.Error(), nil)`. `err.Error()` bisa berisi pesan GORM/driver (nama kolom, SQL) — ini **information disclosure** (flag `gosec` kamu aktif tapi ini lolos karena bukan pattern yang dikenali).

**Fix:** untuk 500, log error lengkap (structured), balikan pesan generik: `"internal server error"`. Sisakan `err.Error()` hanya untuk error domain yang memang aman dibaca user.

### 3.4 🟠 MEDIUM — Tidak ada panic recovery middleware

`fiber.New()` dipanggil tanpa `recover` middleware. Satu nil-deref (mis. `s.db` nil karena `main.go` cuma warning saat DB gagal connect) → **seluruh proses crash**. Untuk enterprise ini tidak boleh.

**Fix:** pasang `recover.New()` + tolak start kalau DB nil di production (`main.go` sekarang `log.Printf("Warning...")` lalu lanjut — sebaiknya `log.Fatal` di env production).

---

## 4. Temuan Konsistensi & Kualitas

### 4.1 UUID generation tidak konsisten antar domain
- **Employee:** UUID digenerate di **domain constructor** (`NewEmployee`).
- **Organization:** UUID digenerate di **application service** (`uuid.New().String()` lalu dioper ke `NewDepartment(id, ...)`).

Rule [uuid-generation.md](../.agents/rules/uuid-generation.md) mewajibkan pola domain-constructor + GORM `BeforeCreate`. Samakan: **application layer jangan generate UUID**. Serahkan ke domain constructor (atau biarkan kosong dan andalkan `BeforeCreate`). Pilih satu, terapkan seragam.

> Catatan: `BeforeCreate` hook hanya jalan pada `Create()`, bukan `Save()`-yang-jadi-`Update`. Selama 3.1 belum dibereskan, hook UUID di model praktis tidak pernah terpakai untuk entity yang lewat `Save()`.

### 4.2 Repository contract bocor: `FindByKTP` return `(nil, nil)`
`FindByKTP` menelan `ErrRecordNotFound` jadi `return nil, nil`. Ini "nil-nil footgun" — caller wajib ingat cek nil manual (`existingKtp != nil && ...`). Lebih aman kembalikan `ErrPersonalDataNotFound` dan biarkan service `errors.Is`-cek, atau kembalikan `(bool, error)` yang eksplisit.

### 4.3 Aggregate boundary bocor (Progressive Save vs DDD)
Employee adalah aggregate root, tapi sub-entity (PersonalData, Contact, Bank) disimpan lewat method repo terpisah dan dipanggil independen. Ini pragmatis untuk form multi-step FE, tapi secara DDD melemahkan invariant "aggregate disimpan sebagai satu unit konsisten". Selama sub-entity tidak punya invariant lintas satu sama lain, ini **acceptable trade-off** — tapi dokumentasikan keputusannya di `decision-log.md` (ADR) supaya sadar-pilihan, bukan kebetulan.

### 4.4 Middleware auth tidak pakai response envelope standar
`auth_middleware.go` balikan `fiber.Map{"error": ...}` — beda format dari `pkg/response`. FE dapat dua bentuk error berbeda. Samakan pakai `response.Error(c, 401, ...)`.

### 4.5 CORS diparse tapi tidak dipasang
`config.AppCorsAllowedOrigins` diparse rapi di config, tapi **tidak ada** `cors` middleware terpasang di `server.go`. Either pasang, atau buang config-nya (dead config membingungkan).

### 4.6 Boilerplate mapping & error-mapping berulang
Tiap handler mengulang blok `if errors.Is(...) { 404 } if errors.Is(...){409} ... 500`. Saat domain bertambah ini meledak. Pertimbangkan **central error mapper**: satu fungsi `mapDomainError(err) (code int, msg string)` atau Fiber `ErrorHandler` global yang menerjemahkan sentinel error ke HTTP.

### 4.7 `interface{}` → `any`
`pkg/response` pakai `interface{}`. Go 1.26 — pakai `any` untuk konsistensi modern (kosmetik).

---

## 5. Enterprise Readiness Gaps

| Gap | Dampak | Rekomendasi |
|-----|--------|-------------|
| **0 unit test** (`find -name '*_test.go'` = 0) | Regresi tak terdeteksi; invariant domain tak terjaga | Mulai dari domain (constructor invariants) & application service (pakai mock repo — interface sudah siap). Target coverage domain+application dulu. |
| Tidak ada structured logging | Debugging produksi buta | Adopsi `slog` (stdlib Go 1.21+), inject logger, request-ID middleware |
| Tidak ada request logging / metrics | Tak ada observability | Fiber logger middleware + (nanti) Prometheus |
| Tidak ada rate limit / timeout | Rentan abuse & hanging query | `timeout` middleware, context deadline di DB |
| RBAC hardcoded `"employee"` | Access control belum nyata | Rancang role/permission sebelum Payroll masuk (data sensitif) |
| DB nil-tolerant startup | Server "jalan" tapi semua query 500 | Fail-fast di production |
| Secret dari `.env` tanpa validasi | JWT secret kosong = token tak aman | Validasi config wajib saat boot (fatal kalau `JWT_SECRET` kosong) |

---

## 6. Prioritas Refactor (Roadmap)

**P0 — sebelum menambah modul baru:**
1. Fix pola `Save()` → upsert eksplisit (3.1). **Ini bug data-loss, paling utama.**
2. Pasang `recover` + fail-fast DB di production (3.4).
3. Berhenti bocorkan `err.Error()` di 500 + log terstruktur (3.3).

**P1 — fondasi enterprise:**
4. Transaction manager / Unit of Work di application layer (3.2).
5. Central domain-error → HTTP mapper (4.6).
6. Test suite domain + application dengan mock repo.
7. Samakan UUID generation antar domain (4.1) & envelope error middleware (4.4).

**P2 — hardening:**
8. Structured logging + request-ID + request logger middleware.
9. Pasang CORS, timeout, rate limit.
10. RBAC nyata sebelum Payroll/Performance.
11. ADR untuk keputusan Progressive Save (4.3).

---

## 7. Kesimpulan

Kerangka Clean Architecture + DDD-nya **sudah benar dan layak jadi fondasi enterprise** — ini modal besar. Yang membedakan "terlihat rapi" vs "production-ready" ada di: **integritas data (3.1), transaction ownership (3.2), dan disiplin operasional (test, recovery, logging)**.

Selesaikan P0 dulu (khususnya bug `Save()`), lalu P1, sebelum membangun Attendance/Payroll. Menumpuk modul di atas pola persistence yang silent-no-op akan menggandakan bug yang sama di setiap domain baru.
