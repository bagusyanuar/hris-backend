---
module: Organization
version: 1.1.0          # 1.0 = base 3-pilar (Department/JobTitle/JobPosition); 1.1 = tambah Multi-Entity layer (Company & Branch)
status: Draft
owner: bagusyanuar
updated: 2026-07-23 13:23:43
depends_on: []
consumed_by: [employee, rbac@planned, payroll@planned, attendance@planned]
---

# Product Requirements: Organization Module

Dokumen ini menjelaskan fungsionalitas, konsep, dan contoh data untuk modul **Organization** di dalam sistem HRIS. Tujuannya adalah untuk menyamakan pemahaman antara tim Backend dan Frontend agar tidak ada kebingungan saat membangun UI/UX.

> **Status implementasi (grounding 2026-07-23):** kode existing (`internal/domain/organization/entity.go`) baru punya 3 pilar (Department, Job Title, Job Position). **Company & Branch di §0 belum ada di kode sama sekali** — seluruh §0 berstatus *planned/intended behavior*, bukan cerminan kode saat ini.

---

## 0. Perluasan Multi-Entity: Company & Branch (PLANNED)

Bagian ini menambah **dua dimensi struktural baru** di atas 3 pilar existing, supaya satu instalasi aplikasi bisa menampung **banyak PT (badan hukum)** sekaligus **banyak cabang (lokasi)**. Ini fondasi multi-tenant enterprise.

### 0.1. Tujuan & Dampak (Why)
Owner mengelola grup usaha: satu holding membawahi beberapa **PT** (badan hukum terpisah), dan tiap PT punya beberapa **cabang** (lokasi fisik). Satu aplikasi, satu login owner, visibilitas seluruh grup — tanpa deploy aplikasi terpisah per PT.

Masalah yang dipecahkan:
- **Pemisahan legal vs lokasi.** Payroll/pajak (PPh 21, BPJS, NPWP) diikat ke PT, bukan lokasi. Absensi/shift/UMR diikat ke cabang. Dua concern beda, wajib dua entity beda.
- **Isolasi data antar-PT & antar-cabang.** HR cabang A tak boleh lihat data cabang B; HR PT-X tak boleh lihat PT-Y. Owner holding lihat semua.
- **Konsolidasi.** Owner butuh laporan agregat lintas PT/cabang, sementara tiap PT tetap tutup buku sendiri.

### 0.2. Scope & Out-of-Scope

**In-Scope (§0 ini):**
- Entity **Company** (PT / badan hukum) — master data legal (nama legal, NPWP, alamat, BPJS registration).
- Entity **Branch** (cabang / lokasi) — nempel ke satu Company.
- Kolom dimensi `company_id` + `branch_id` sebagai fondasi scoping di seluruh entity operasional (termasuk 3 pilar existing & Employee).
- Aturan integritas: Branch wajib milik satu Company; Employee wajib punya `company_id` **dan** `branch_id` (non-nullable).

**Out-of-Scope (dikerjakan di PRD lain):**
- **RBAC branch-scoped / company-scoped access control** → PRD **`rbac` sendiri** (belum ada, WAJIB dibuat). §0 hanya menyediakan *dimensi* datanya, bukan enforcement-nya.
- Kalkulasi payroll/pajak per-PT → PRD **Payroll**.
- Kebijakan shift/UMR/kalender libur per-cabang → PRD **Attendance/Leave**.
- Transfer karyawan antar-cabang/antar-PT (riwayat penempatan) → fase lanjutan Employee (disebut di §0.4 sebagai kontrak, implementasi menyusul).

### 0.3. User Roles & Permissions (ringkas — detail di PRD RBAC)

| Role | Company scope | Branch scope | Baca | Tulis |
|------|--------------|-------------|------|-------|
| **Owner / Group Admin** | semua PT | semua cabang | ✅ | ✅ |
| **Company Admin (HR PT-X)** | 1 PT | semua cabang PT-X | ✅ PT-X | ✅ PT-X |
| **Branch Admin (HR cabang)** | 1 PT | 1 cabang | ✅ cabang sendiri | ✅ cabang sendiri |
| **Employee** | PT sendiri | cabang sendiri | ✅ diri sendiri | ❌ |

> Enforcement scoping (inject filter `company_id`/`branch_id` di query boundary) = tanggung jawab modul **RBAC**, bukan §0. §0 cuma jamin kolomnya ada & terisi.

### 0.4. Kriteria Penerimaan (Given-When-Then)

- **Company unik per NPWP.**
  *Given* sudah ada Company dengan NPWP `01.234.567.8-901.000`, *When* admin buat Company baru dengan NPWP sama, *Then* tolak `409 Conflict` (`ErrCompanyNPWPDuplicate`).

- **Branch wajib milik Company valid.**
  *Given* `company_id` tak dikenal, *When* buat Branch, *Then* tolak `422` (validasi FK), tidak boleh Branch yatim.

- **Employee wajib dua dimensi.**
  *Given* payload Employee tanpa `company_id` atau tanpa `branch_id`, *When* create, *Then* tolak `422`. Tidak ada "default company".
  *Catatan implementasi:* entity Employee existing belum punya `company_id`/`branch_id` — ini **gap** yang wajib ditutup saat §0 diimplementasi (butuh migrasi + ubah domain Employee).

- **Branch harus se-Company dengan Employee-nya.**
  *Given* Employee `company_id = PT-A`, *When* set `branch_id` yang cabangnya milik `PT-B`, *Then* tolak `422` (cross-entity mismatch).

- **Job Position terikat Company.**
  *Given* struktur 3 pilar, *When* buat Job Position, *Then* posisi wajib nempel ke satu `company_id` (reporting line & headcount quota berlaku dalam batas satu PT, tidak lintas PT).
  *Catatan implementasi:* 3 pilar existing belum punya `company_id` — gap yang ikut ditutup di §0.

- **Isolasi baca (kontrak untuk RBAC).**
  *Given* Branch Admin cabang A, *When* GET daftar Employee, *Then* hanya Employee `branch_id = A` yang balik. *Catatan:* enforcement di PRD RBAC; §0 hanya menyediakan kolom filternya.

### 0.5. Technical & Architectural Constraints
- **Bounded context:** Company & Branch masuk context **Organization** (`internal/organization/`), bukan context baru. Mereka struktur organisasi, bukan concern lintas-modul. Yang lintas-modul (scoping/RBAC) dipisah ke PRD RBAC — loose coupling.
- **Dimensi non-nullable:** `company_id` (dan `branch_id` untuk entity operasional) **NOT NULL** sejak migrasi awal. Kolom nullable = pintu masuk data ambigu yang meracuni payroll. Retrofit belakangan = neraka migrasi.
- **Isolasi data:** mulai dengan **shared DB + kolom `company_id`/`branch_id`** (row-level). Jangan schema-per-tenant / DB-per-tenant sampai ada tuntutan regulasi nyata.
- **Sentinel error:** `ErrCompanyNotFound`, `ErrBranchNotFound`, `ErrCompanyNPWPDuplicate`, `ErrBranchCompanyMismatch` — patuhi [persistence-convention.md](../../.agents/rules/persistence-convention.md) §3.
- **UUID:** generate di domain constructor (`NewCompany`, `NewBranch`), `BeforeCreate` sebagai safety-net — patuhi [uuid-generation.md](../../.agents/rules/uuid-generation.md).
- **Tier dokumen:** perluasan ini **Kompleks** (relasi multi-entity + jadi dependency banyak modul) → butuh `tech-spec.md` + `decision-log.md` (ADR: kenapa 2 kolom `company_id`+`branch_id`, kenapa row-level bukan schema-per-tenant).

### 0.6. Dependencies
- **Depends on:** — (Company/Branch adalah akar, tak bergantung modul lain).
- **Consumed by:**
  - **Employee** — tambah `company_id` + `branch_id` (breaking untuk skema Employee existing).
  - **RBAC** (PRD baru, WAJIB dibuat) — pakai dua dimensi ini untuk scoping akses. **Ini gap dependency**: belum ada `docs/PRD/rbac.md`.
  - **Payroll** (planned) — group by `company_id` untuk PPh 21/BPJS/NPWP per-PT.
  - **Attendance/Leave** (planned) — konfigurasi per `branch_id` (shift, UMR, kalender libur).
- **External:** — (tidak ada integrasi eksternal di §0).

### 0.7. Data Schema (PLANNED — sample buat FE, bukan pengganti DBML)

**Entity `Company` (PT / badan hukum)**
Aturan bisnis: `npwp` unik; `is_active` soft-toggle; root dari hierarki, tidak punya parent.

| id | code | legal_name | npwp | bpjs_no | is_active |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `co-1` | PTA | PT Alpha Nusantara | `01.234.567.8-901.000` | `JKN-0001` | true |
| `co-2` | PTB | PT Beta Sejahtera | `02.345.678.9-012.000` | `JKN-0002` | true |

**Entity `Branch` (cabang / lokasi)**
Aturan bisnis: wajib `company_id` valid; `code` unik dalam satu Company (boleh sama antar-Company); `is_main` menandai kantor pusat PT.

| id | company_id | code | name | city | is_main | is_active |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `br-1` | `co-1` | JKT | Kantor Pusat Jakarta | Jakarta | true | true |
| `br-2` | `co-1` | SBY | Cabang Surabaya | Surabaya | false | true |
| `br-3` | `co-2` | BDG | Kantor Pusat Bandung | Bandung | true | true |

> Setelah §0 diimplementasi, tabel sample 3 pilar (§1–§3) & Employee bertambah kolom `company_id` (dan `branch_id` untuk Employee). Belum direfleksikan di tabel §1–§3 di bawah karena masih *planned*.

---

## Ringkasan Kompleksitas (buat owner — seberapa berat §0)

| Area | Dampak | Berat |
|------|--------|-------|
| Entity baru Company + Branch (domain, adapter, transport, DTO, CRUD, test) | 2 entity full-stack | **Sedang** |
| Migrasi tambah `company_id`+`branch_id` ke 3 pilar + Employee | breaking skema existing, butuh backfill | **Berat** |
| Ubah domain Employee & 3 pilar (constructor, validasi cross-entity) | sentuh kode existing yang sudah jalan | **Berat** |
| PRD + implementasi **RBAC** scoping (modul baru terpisah) | fondasi access-control lintas semua modul | **Berat** — proyek sendiri |
| Payroll/Attendance per-PT/per-cabang | belum ada modulnya, ikut arah §0 | ditunda |

**Verdict analis:** §0 (Company+Branch entity + kolom dimensi) = tier **Kompleks**, tapi *tractable* kalau dipecah bertahap:
1. **Fase 1** — bikin entity Company + Branch (CRUD murni, belum sentuh Employee). Aman, tak breaking.
2. **Fase 2** — migrasi `company_id`/`branch_id` ke Employee + 3 pilar + validasi cross-entity. Ini fase breaking, butuh hati-hati + backfill data existing.
3. **Fase 3** — PRD & modul **RBAC** untuk enforce scoping. Ini yang paling berat, layak jadi track terpisah.

Jangan gabung tiga fase jadi satu PR. Fase 1 bisa jalan sekarang tanpa ganggu apa pun.

---

## Konsep Dasar (The 3-Pillars)

Sistem organisasi perusahaan di dalam HRIS ini dibangun di atas 3 pilar utama:
1. **Department**: Struktur unit kerja perusahaan. Bisa berupa divisi utama, departemen, atau sub-departemen. Memiliki relasi *hierarki* (Parent-Child).
2. **Job Title**: Master data pangkat, golongan, atau jenjang karir baku yang berlaku secara umum di perusahaan.
3. **Job Position**: Slot jabatan aktual / definitif di dalam struktur organisasi yang merupakan kombinasi antara **Department** dan **Job Title**. Jabatan ini mendefinisikan "Siapa lapor ke siapa" (*reporting line*) serta jumlah batasan kuota pegawai (*headcount quota*).

### Mengapa Menggunakan Konsep 3 Pilar? (Position-Based vs Person-Based)
Di banyak sistem HR sederhana (UMKM), biasanya karyawan langsung ditempelkan nama jabatannya (*Person-Based*). Namun, untuk sistem *Enterprise-Grade*, kita menggunakan pendekatan *Position-Based* (3 Pilar). Artinya, struktur organisasi dan "kursi" jabatannya dibentuk terlebih dahulu, baru kemudian karyawan menduduki kursi tersebut.

**Pros (Kelebihan):**
- **Struktur Independen:** Jika seorang manajer *resign*, struktur pelaporan (*reporting line*) di bawahnya tidak rusak karena bawahan melapor ke *Posisi* manajer, bukan ke *Orang*-nya.
- **Manajemen Kuota & Budget (Headcount):** Memudahkan finance dan HR untuk membatasi jumlah pegawai (contoh: Posisi "Staf IT" hanya boleh diisi maksimal 5 orang).
- **Standarisasi Gaji (Grade):** Pemisahan `Job Title` (Pangkat) memastikan standarisasi gaji/fasilitas yang adil di lintas departemen (contoh: Manajer IT dan Manajer HR berada di *grade* yang sama).

**Cons (Kekurangan):**
- **Kompleksitas Awal (Setup):** Butuh waktu ekstra di awal untuk mengatur master data. Admin HR tidak bisa langsung menambah karyawan, mereka harus membuat Departemen, lalu Job Title, dan menyatukannya jadi Job Position terlebih dahulu.
- **Kurang Cocok untuk Start-Up Kecil:** Perusahaan dengan struktur yang sangat cair (pegawai merangkap banyak peran abstrak) mungkin merasa sistem ini terlalu kaku.

---

## 1. Department (Unit Kerja)
Menyimpan struktur divisi atau departemen. Relasi bersifat *Tree* atau hierarkis menggunakan `parent_id`.

### Aturan Bisnis:
- Jika `parent_id` adalah `null`, berarti ini adalah departemen level tertinggi (Root).
- Di frontend, tampilan ini biasanya direpresentasikan sebagai **Tree View** atau nested list.

### Sample Data (Tabel `departments`):

| id | code | name | parent_id | is_active |
| :--- | :--- | :--- | :--- | :--- |
| `dept-1` | DIR | Direksi | `null` | true |
| `dept-2` | TI | Divisi Teknologi Informasi | `dept-1` | true |
| `dept-3` | DEV | Departemen Pengembangan (Engineering) | `dept-2` | true |
| `dept-4` | OPR | Divisi Operasional | `dept-1` | true |
| `dept-5` | SDM | Divisi Sumber Daya Manusia (HR) | `dept-1` | true |

---

## 2. Job Title (Pangkat / Grade)
Master data standarisasi jabatan atau jenjang karir. Biasanya digunakan untuk menentukan standar gaji (*Salary Band*) atau fasilitas (*Benefit*).

### Aturan Bisnis:
- `grade_level` menentukan tinggi/rendahnya pangkat secara angka (misalnya makin tinggi angkanya, makin tinggi pangkatnya).
- Tidak terkait dengan departemen tertentu (independen).

### Sample Data (Tabel `job_titles`):

| id | code | name | grade_level | is_active |
| :--- | :--- | :--- | :--- | :--- |
| `title-1` | DIR | Direktur | 10 | true |
| `title-2` | KDV | Kepala Divisi / GM | 9 | true |
| `title-3` | MGR | Manajer | 7 | true |
| `title-4` | SPV | Supervisor | 5 | true |
| `title-5` | STF | Staf | 3 | true |

---

## 3. Job Position (Jabatan Aktif / Posisi)
Ini adalah "kursi" aktual yang diduduki oleh pegawai di dalam struktur organisasi.

### Aturan Bisnis:
- **Kombinasi**: Setiap Job Position harus menempel pada satu `Department` dan satu `Job Title`.
- **Reporting Line**: `reports_to_id` menunjuk ke ID Job Position lain sebagai atasannya, membentuk **Organization Chart** (Bagan Struktur Organisasi).
- **Headcount Quota**: Menentukan batas maksimal pegawai yang bisa menduduki jabatan ini (misalnya, CEO kuotanya 1, tapi Software Engineer kuotanya bisa 10).

### Sample Data (Tabel `job_positions`):

| id | name (Posisi) | department_id | job_title_id | reports_to_id (Atasan) | headcount_quota |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `pos-1` | Direktur Utama | `dept-1` (DIR) | `title-1` (DIR) | `null` | 1 |
| `pos-2` | Direktur Teknologi (CTO) | `dept-2` (TI) | `title-1` (DIR) | `pos-1` (Dirut) | 1 |
| `pos-3` | Kepala Divisi TI | `dept-2` (TI) | `title-2` (KDV) | `pos-2` (Dirtek) | 1 |
| `pos-4` | Manajer Pengembangan | `dept-3` (DEV) | `title-3` (MGR) | `pos-3` (Kadiv TI) | 3 |
| `pos-5` | Supervisor Backend | `dept-3` (DEV) | `title-4` (SPV) | `pos-4` (Mgr DEV) | 5 |
| `pos-6` | Staf Programmer Backend | `dept-3` (DEV) | `title-5` (STF) | `pos-5` (Spv BE) | 10 |

---

## Catatan Khusus untuk Frontend (FE)

1. **Pembuatan Form `Job Position`**:
   - Di form "Create Job Position", FE membutuhkan dropdown untuk memilih `Department` dan `Job Title`. Oleh karena itu, FE harus memanggil API `GET /organization/departments` dan `GET /organization/job-titles` terlebih dahulu untuk mengisi *dropdown option*.
   - Input `Reports To` adalah *autocomplete dropdown* yang mengambil data dari `GET /organization/job-positions`.
   
2. **Organization Chart (Bagan Organisasi)**:
   - Data `job_positions` yang saling terkait lewat `reports_to_id` bisa dirender menjadi **Organization Chart** visual.
   - Root (puncak) dari chart adalah posisi dengan `reports_to_id` bernilai `null` (seperti contoh `pos-1` CEO di atas).

3. **Status `is_active`**:
   - Secara *default*, data yang dikembalikan oleh API adalah yang aktif (jika tidak difilter). Jika suatu departemen/posisi di-nonaktifkan, UI dapat menampilkannya dengan warna abu-abu (greyed out) atau disembunyikan.
