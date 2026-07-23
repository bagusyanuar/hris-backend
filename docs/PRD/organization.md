---
module: Organization
version: 2.0.0          # 2.0 = scope dipersempit jadi legal/lokasi (Company & Branch); Dept/Title/Position pindah ke workforce-structure.md
status: Draft
owner: bagusyanuar
updated: 2026-07-23 13:46:17
depends_on: []
consumed_by: [workforce-structure, employee, rbac@planned, payroll@planned, attendance@planned]
---

# Product Requirements: Organization Module

Modul **Organization** mengatur *legal entity & lokasi* sebuah grup usaha: **Company** (PT / badan hukum) dan **Branch** (cabang / lokasi fisik). Ini fondasi multi-entity yang bikin satu platform menaungi beberapa PT, tiap PT beberapa cabang — sesuai [Global PRD / product-vision.md](product-vision.md).

> **⚠️ Breaking scope change (2.0.0, 2026-07-23):** versi 1.x modul ini mencakup Department, Job Title, Job Position. Konsep itu **dipindah** ke modul baru [workforce-structure.md](workforce-structure.md). Sejak 2.0.0, Organization **hanya** soal Company & Branch. Kode existing (`internal/domain/organization/`) masih memuat 3 pilar lama — pemindahan folder ke `internal/workforce/` + entity Company/Branch baru = **gap implementasi** (lihat §5).

---

## 1. Tujuan & Dampak (Why)
Owner mengelola grup usaha: satu holding membawahi beberapa **PT** (badan hukum terpisah), tiap PT punya beberapa **cabang**. Satu aplikasi, satu login, visibilitas seluruh grup — tanpa deploy terpisah per PT.

Masalah yang dipecahkan:
- **Pemisahan legal vs lokasi.** Payroll/pajak (PPh 21, BPJS, NPWP) diikat ke **Company**, bukan lokasi. Absensi/shift/UMR diikat ke **Branch**. Dua concern beda, dua entity beda.
- **Isolasi data antar-PT & antar-cabang.** HR cabang A tak lihat cabang B; HR PT-X tak lihat PT-Y. Owner holding lihat semua.
- **Konsolidasi.** Owner butuh agregat lintas PT/cabang; tiap PT tetap tutup buku sendiri.

## 2. Scope & Out-of-Scope

**In-Scope:**
- Entity **Company** (PT / badan hukum) — master legal: nama legal, NPWP, alamat, BPJS registration.
- Entity **Branch** (cabang / lokasi) — nempel ke satu Company.
- Kolom dimensi `company_id` + `branch_id` sebagai fondasi scoping di seluruh entity operasional (Workforce Structure, Employee, dst).
- Integritas: Branch wajib milik satu Company; entity operasional wajib `company_id` (+ `branch_id` untuk yang lokasi-spesifik).

**Out-of-Scope (PRD lain):**
- **Department, Job Title, Job Position** → [workforce-structure.md](workforce-structure.md).
- **RBAC company/branch-scoped access control** → PRD **`rbac`** (belum ada, WAJIB dibuat). Organization sediakan *dimensi* datanya, bukan enforcement-nya.
- **Kalkulasi payroll/pajak per-PT** → PRD **Payroll**.
- **Shift/UMR/kalender libur per-cabang** → PRD **Attendance/Leave**.
- **Transfer karyawan antar-cabang/PT** (riwayat penempatan) → fase lanjutan **Employee**.

## 3. User Roles & Permissions (ringkas — detail di PRD RBAC)

| Role | Company scope | Branch scope | Baca | Tulis |
|------|--------------|-------------|------|-------|
| **Owner / Group Admin** | semua PT | semua cabang | ✅ | ✅ |
| **Company Admin (HR PT-X)** | 1 PT | semua cabang PT-X | ✅ PT-X | ✅ PT-X |
| **Branch Admin (HR cabang)** | 1 PT | 1 cabang | ✅ cabang sendiri | ✅ cabang sendiri |
| **Employee** | PT sendiri | cabang sendiri | ✅ diri sendiri | ❌ |

> Enforcement scoping (inject filter `company_id`/`branch_id` di query boundary) = tanggung jawab modul **RBAC**. Organization cuma jamin kolomnya ada & terisi.

## 4. Kriteria Penerimaan (Given-When-Then)

- **Company unik per NPWP.**
  *Given* sudah ada Company NPWP `01.234.567.8-901.000`, *When* buat Company baru NPWP sama, *Then* tolak `409` (`ErrCompanyNPWPDuplicate`).

- **Branch wajib Company valid.**
  *Given* `company_id` tak dikenal, *When* buat Branch, *Then* tolak `422`. Tidak boleh Branch yatim.

- **Satu Company punya satu kantor pusat.**
  *Given* sudah ada Branch `is_main = true` di Company X, *When* set Branch lain jadi `is_main` di Company X, *Then* pindahkan status / tolak sesuai aturan tech-spec (satu main per Company).

- **Branch code unik dalam Company.**
  *Given* Company X sudah punya Branch code `JKT`, *When* buat Branch `JKT` lagi di X, *Then* tolak `409`. (Code sama boleh di Company beda.)

- **Entity operasional wajib dimensi.**
  *Given* payload Employee tanpa `company_id`/`branch_id`, *When* create, *Then* tolak `422`. Tidak ada "default company".
  *Catatan implementasi:* entity Employee & Department existing **belum punya** `company_id`/`branch_id` — **gap** yang wajib ditutup saat implementasi (migrasi + ubah domain).

- **Branch harus se-Company dengan pemakainya.**
  *Given* Employee `company_id = PT-A`, *When* set `branch_id` yang cabangnya milik `PT-B`, *Then* tolak `422` (cross-entity mismatch).

## 5. Technical & Architectural Constraints
- **Bounded context:** Company & Branch = isi modul `organization` (`internal/organization/`). Dept/Title/Position **keluar** ke `workforce-structure`.
- **Dimensi non-nullable:** `company_id` (+ `branch_id` untuk entity lokasi-spesifik) **NOT NULL** sejak migrasi awal. Kolom nullable = pintu data ambigu yang meracuni payroll. Retrofit belakangan = neraka migrasi.
- **Isolasi data:** **shared DB + kolom `company_id`/`branch_id`** (row-level). Bukan schema-per-tenant / DB-per-tenant, kecuali tuntutan regulasi nyata (lihat [product-vision.md](product-vision.md) §2).
- **Sentinel error:** `ErrCompanyNotFound`, `ErrBranchNotFound`, `ErrCompanyNPWPDuplicate`, `ErrBranchCompanyMismatch`, `ErrBranchCodeDuplicate` — [persistence-convention.md](../../.agents/rules/persistence-convention.md) §3.
- **UUID:** generate di domain constructor (`NewCompany`, `NewBranch`), `BeforeCreate` safety-net — [uuid-generation.md](../../.agents/rules/uuid-generation.md).
- **Tier dokumen:** **Kompleks** (relasi multi-entity + jadi dependency banyak modul) → butuh `tech-spec.md` + `decision-log.md` (ADR: kenapa 2 kolom `company_id`+`branch_id`, kenapa row-level bukan schema-per-tenant).

## 6. Dependencies
- **Depends on:** — (Company/Branch adalah akar).
- **Consumed by:**
  - **Workforce Structure** — Department konsumsi `company_id` untuk scope struktur.
  - **Employee** — tambah `company_id` + `branch_id` (breaking untuk skema Employee existing).
  - **RBAC** (PRD baru, WAJIB dibuat) — pakai dua dimensi ini untuk scoping. **Gap dependency:** belum ada `docs/PRD/rbac.md`.
  - **Payroll** (planned) — group by `company_id` untuk PPh 21/BPJS/NPWP per-PT.
  - **Attendance/Leave** (planned) — konfigurasi per `branch_id`.
- **External:** —

---

## 7. Data Schema & Business Rules

> Sample buat FE, bukan pengganti DBML.

### 7.1. Company (`companies`) — PT / badan hukum
Aturan: `npwp` unik; root hierarki (tanpa parent); `is_active` soft-toggle.

| id | code | legal_name | npwp | bpjs_no | is_active |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `co-1` | PTA | PT Alpha Nusantara | `01.234.567.8-901.000` | `JKN-0001` | true |
| `co-2` | PTB | PT Beta Sejahtera | `02.345.678.9-012.000` | `JKN-0002` | true |

### 7.2. Branch (`branches`) — cabang / lokasi
Aturan: wajib `company_id` valid; `code` unik dalam satu Company; `is_main` = kantor pusat PT (satu per Company).

| id | company_id | code | name | city | is_main | is_active |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `br-1` | `co-1` | JKT | Kantor Pusat Jakarta | Jakarta | true | true |
| `br-2` | `co-1` | SBY | Cabang Surabaya | Surabaya | false | true |
| `br-3` | `co-2` | BDG | Kantor Pusat Bandung | Bandung | true | true |

---

## Ringkasan Kompleksitas & Fase (buat owner)

| Area | Dampak | Berat |
|------|--------|-------|
| Entity Company + Branch (domain/adapter/transport/DTO/CRUD/test) | 2 entity full-stack baru | **Sedang** |
| Pindah Dept/Title/Position ke modul `workforce-structure` | rename folder/package/DI/import | **Sedang** |
| Migrasi tambah `company_id`/`branch_id` ke Employee + Department | breaking skema existing + backfill | **Berat** |
| PRD + implementasi **RBAC** scoping | fondasi access-control lintas modul | **Berat** — track sendiri |

**Fase eksekusi:**
1. **Fase 1** — entity Company + Branch (CRUD murni). Aman, tak breaking.
2. **Fase 2** — pindah Dept/Title/Position ke `workforce-structure` + tambah `company_id` (migrasi breaking).
3. **Fase 3** — migrasi `company_id`/`branch_id` ke Employee + validasi cross-entity.
4. **Fase 4** — PRD & modul **RBAC** enforce scoping (track terpisah, paling berat).

Jangan gabung fase jadi satu PR. Fase 1 bisa jalan sekarang tanpa ganggu apa pun.
