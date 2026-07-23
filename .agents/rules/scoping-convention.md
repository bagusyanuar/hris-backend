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
