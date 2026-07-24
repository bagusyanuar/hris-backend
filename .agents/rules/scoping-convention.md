# Konvensi Scoping Multi-Entity (Company & Branch)

Aturan ini mengunci **aksioma multi-entity** dari [Global PRD / product-vision.md](../../docs/PRD/product-vision.md) §2 ke level kode & tooling: aplikasi = **satu grup usaha (single-owner holding)** yang menaungi banyak **Company (PT)**, tiap PT banyak **Branch (cabang)**. Supaya data satu PT/cabang tidak bocor ke yang lain, **setiap entity operasional WAJIB scope-aware sejak lahir** — bukan retrofit belakangan (retrofit = neraka migrasi).

Semua aturan di bawah bersifat **WAJIB (STRICT)**.

---

## 1. Klasifikasi Entity (tentukan SEBELUM scaffold)

Sebelum bikin entity baru, klasifikasikan dulu kelas scope-nya. Ini menentukan kolom apa yang wajib ada:

| Kelas | Kolom scope wajib | Contoh | Aturan |
|-------|-------------------|--------|--------|
| **Legal root** | — (dia *adalah* scope) | `Company` | tidak punya `company_id`; root hierarki |
| **Location root** | `company_id` | `Branch` | milik satu Company; belum punya `branch_id` |
| **Company-owned** | `company_id` NOT NULL | Department, Job Title, Job Position | scope per-PT; tak lintas Company |
| **Company + Location bound** | `company_id` + `branch_id` NOT NULL | Employee, Attendance, Leave record | scope per-PT **dan** per-cabang |
| **Global master** | — (eksplisit, JARANG) | lookup lintas-PT sejati (mis. daftar negara) | HANYA kalau data identik untuk semua PT; default BUKAN ini |

**Default untuk modul operasional baru = `company_id` NOT NULL** (kelas Company-owned). Tambahkan `branch_id` bila entity lokasi-spesifik. Jangan pilih "Global master" kecuali datanya benar-benar identik lintas seluruh PT — ini keputusan sadar, bukan default malas.

> Contoh keputusan project: **Job Title = per-PT** (Company-owned, punya `company_id`), bukan master global. Tiap PT punya grade/pangkat sendiri.

---

## 2. Kolom & Foreign Key

- `company_id` → FK ke `companies(id)`, **NOT NULL** untuk kelas Company-owned & Company+Location.
- `branch_id` → FK ke `branches(id)`, **NOT NULL** untuk kelas Company+Location.
- **Nullable DILARANG** untuk kolom scope pada kelas yang mewajibkannya. Kolom scope nullable = pintu masuk data ambigu yang meracuni payroll/pajak (tak jelas milik PT mana).
- **Index wajib** di kolom scope (`company_id`, dan `(company_id, branch_id)` untuk entity dua-dimensi) — semua query difilter lewat kolom ini, tanpa index = full scan.
- **Integritas silang:** `branch_id` yang dipilih WAJIB milik `company_id` yang sama. Repository/domain WAJIB tolak mismatch dengan sentinel error (mis. `ErrBranchCompanyMismatch`).

---

## 3. Filter di Query Boundary (WAJIB)

Scope **tidak boleh** mengandalkan caller ingat nambah `WHERE company_id = ?` manual tiap query. Pola wajib, selaras propagasi context di [logging-convention.md](logging-convention.md) §2:

1. **RBAC middleware** menaruh scope user (allowed `company_id` + allowed `branch_id`) ke `context.Context` (mis. `scope.WithContext`), sejajar cara `request_id` diinjeksi.
2. **Repository** `FindXxx`/`ListXxx` WAJIB membaca scope dari context dan meng-inject filter `WHERE company_id IN (...)` (dan `branch_id` bila relevan) ke setiap query baca.
3. **Write** (Create/Update) WAJIB memvalidasi entity yang ditulis berada dalam scope yang diizinkan — cegah user cabang A bikin/ubah data cabang B.

Kontrak minimal abstraksi scope (shared kernel, bukan detail GORM):
```go
type Scope struct {
    CompanyIDs []string // PT yang boleh diakses; kosong = semua (owner/group admin)
    BranchIDs  []string // cabang yang boleh diakses; kosong = semua cabang dalam CompanyIDs
}
// scope.FromContext(ctx) Scope
```

> **Owner / Group Admin** = scope kosong (akses semua PT & cabang). **Company Admin** = 1 `company_id`, semua branch. **Branch Admin** = 1 company + 1 branch. Detail role di PRD **RBAC**.

### 3.1. Active Scope Selector — Header (`X-Company-Id`/`X-Branch-Id`), BUKAN Query Param

FE (mis. workspace switcher di sidebar — pilih PT & cabang aktif) butuh cara ngirim tahu backend "PT/cabang mana yang lagi aktif dilihat user". Ini **beda konsep** dari `Scope` di atas:

- **`scope.FromContext(ctx)`** = *allowed set* — batas otorisasi user, WAJIB diturunkan dari JWT/RBAC assignment server-side (DB lookup role→company/branch). **TIDAK PERNAH** dari input client.
- **Active selector** = pilihan "lagi kerja di PT/cabang mana" saat ini — dikirim client via header **`X-Company-Id`** dan **`X-Branch-Id`** (opsional, cuma diisi kalau relevan buat endpoint itu).

**Kenapa header, bukan query param:** pilihan ini cross-cutting — kepake di hampir semua request (`POST`/`PUT`/`DELETE`, bukan cuma `GET` List/filter). Query param berarti tiap endpoint URL kudu bawa `?company_id=...` manual & rawan kelupaan di satu endpoint; header cukup di-set sekali di HTTP client interceptor (axios `defaults.headers`/fetch wrapper) sekali pas switch + sekali pas app load (baca dari localStorage/cookie), otomatis nempel ke semua request abis itu.

**WAJIB divalidasi server-side, jangan dipercaya mentah** — alur middleware:
1. JWT diverifikasi → dapat `scope.FromContext(ctx)` asli (allowed set, dari assignment DB, bukan header).
2. Middleware baca header `X-Company-Id`/`X-Branch-Id`.
3. Validasi: nilai header **WAJIB subset** dari allowed set (`scope.CompanyIDs` kosong = Owner = header apa aja boleh; kalau `scope.CompanyIDs` terisi, header WAJIB salah satu isinya — kalau tidak, `403 Forbidden`, jangan diam-diam di-ignore).
4. Kalau valid, itu jadi **active scope** request ini — dipakai buat filter query DB (gantiin/mempersempit allowed set, bukan menambah akses).

Endpoint yang gak perlu tau active company/branch (mis. `GET /companies` punya Owner) boleh abaikan header ini sepenuhnya.

---

## 4. Sequencing — Enforcement Bertahap (baca sebelum scaffold modul sekarang)

Modul **Organization** (Company/Branch) & **RBAC** **BELUM ada di kode** (per 2026-07-23). Jadi:

- **Yang WAJIB sekarang** saat scaffold modul baru: entity/model/DTO/DBML **sudah bawa kolom** `company_id`(+`branch_id`) sesuai kelasnya, dan repository `Find`/`List` **sudah scope-aware** (baca `scope.FromContext`).
- **Yang STAGED** sampai Organization landing (Fase 1) & RBAC landing (Fase 4, lihat [product-vision.md](../../docs/PRD/product-vision.md) §5.4):
  - FK fisik `company_id → companies(id)` baru bisa jalan setelah tabel `companies`/`branches` ada. Sampai itu, DBML tetap **deklarasikan** FK-nya (dokumentasi), migrasi SQL urutkan setelah migrasi Organization.
  - Isi `scope.FromContext` di-populate penuh oleh RBAC middleware. Sebelum RBAC ada, scope boleh default kosong (owner-mode) — TAPI kolom & signature repository tetap wajib scope-aware supaya tak perlu bongkar ulang nanti.

Intinya: **struktur (kolom + signature) dipaku sekarang; pengisian enforcement nyusul.** Jangan tunda kolomnya.

---

## 5. Checklist Review / PR

- [ ] Tiap entity baru sudah diklasifikasikan kelas scope-nya (§1); `company_id`(+`branch_id`) ada sesuai kelas.
- [ ] Kolom scope **NOT NULL** + ber-index; FK dideklarasikan di DBML.
- [ ] `FindXxx`/`ListXxx` membaca `scope.FromContext` dan inject filter — tidak ada query baca tanpa filter scope.
- [ ] Write memvalidasi entity dalam scope; `branch_id` se-`company_id` (mismatch → sentinel error).
- [ ] Tidak ada kolom scope nullable pada kelas yang mewajibkannya.
- [ ] Kalau entity dipilih "Global master" (tanpa scope) — ada justifikasi eksplisit datanya identik lintas PT.
- [ ] Kalau endpoint baca header `X-Company-Id`/`X-Branch-Id` (§3.1) — nilainya divalidasi subset dari `scope.FromContext(ctx)` sebelum dipakai filter, TIDAK dipercaya mentah dari client.
