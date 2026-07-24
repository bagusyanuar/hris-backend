---
module: Workforce Structure
version: 1.2.0
status: Approved
owner: bagusyanuar
updated: 2026-07-24
depends_on: [organization@2.0.0]
consumed_by: [employee]
---

# Product Requirements: Workforce Structure Module

Modul **Workforce Structure** mengatur *bagan organisasi internal* sebuah perusahaan: unit kerja (Department), pangkat/grade (Job Title), dan jabatan aktual (Job Position) beserta reporting line & headcount. Ini "kerangka kursi" yang nanti diduduki karyawan.

> **Pemisahan dari Organization (grounding 2026-07-23):** dulu konsep ini menyatu di modul `Organization`. Sejak split multi-entity, **Organization** dipersempit jadi legal/lokasi (Company & Branch), dan struktur internal (Department/Title/Position) pindah ke modul ini. Kode existing (`internal/domain/organization/`) masih memuat 3 pilar ini — **belum dipindah** ke `internal/<workforce-structure>/`; itu gap implementasi (lihat §5).

---

## Konsep Dasar (The 3-Pillars)

Bagan organisasi dibangun di atas 3 pilar:
1. **Department** — unit kerja (divisi/departemen/sub-departemen). Relasi hierarki Parent-Child (`parent_id`).
2. **Job Title** — master pangkat/golongan/jenjang karir baku. Menentukan salary band/benefit.
3. **Job Position** — "kursi" jabatan aktual = kombinasi **Department × Job Title**. Mendefinisikan reporting line (siapa lapor ke siapa) & headcount quota.

### Kenapa Position-Based (bukan Person-Based)?
Sistem HR sederhana nempel jabatan langsung ke orang (*Person-Based*). Kita pakai *Position-Based* (3 pilar): struktur & "kursi" dibentuk dulu, baru karyawan menduduki.

**Pros:**
- **Struktur independen:** manajer resign, reporting line bawahan tak rusak (lapor ke *posisi*, bukan *orang*).
- **Manajemen kuota/budget (headcount):** batasi jumlah pegawai per jabatan.
- **Standarisasi gaji (grade):** `Job Title` terpisah menjamin keadilan gaji lintas departemen.

**Cons:**
- **Setup awal lebih berat:** harus buat Department → Job Title → satukan jadi Job Position dulu.
- **Kurang cocok start-up cair** yang pegawainya merangkap banyak peran abstrak.

---

## 1. Tujuan & Dampak (Why)
Menyediakan kerangka organisasi *position-based* yang stabil: reporting line tak rusak saat turnover, headcount terkontrol, grade gaji terstandarisasi. Tanpa ini, penempatan karyawan jadi ad-hoc dan bagan organisasi tak bisa diaudit.

## 2. Scope & Out-of-Scope

**In-Scope:**
- CRUD Department (hierarki tree via `parent_id`).
- CRUD Job Title (grade level).
- CRUD Job Position (Department × Job Title + `reports_to_id` + `headcount_quota`).
- Render data buat Organization Chart.

**Out-of-Scope:**
- **Penempatan karyawan ke Position** (occupancy) → modul **Employee**. Modul ini cuma sediakan "kursi"-nya.
- **Company/Branch** (legal/lokasi) → modul **Organization**.
- **Perhitungan gaji** dari grade → modul **Payroll**.
- **Enforce headcount saat assign** → kontrak di §4 (enforcement-nya di Employee/RBAC).

## 3. User Roles & Permissions
| Role | Baca | Tulis |
|------|------|-------|
| **Owner / Company Admin** | ✅ struktur PT-nya | ✅ |
| **HR Manager** | ✅ | ✅ (dalam scope Company-nya) |
| **Employee** | ✅ (lihat bagan) | ❌ |

> Scoping per `company_id` = tanggung jawab RBAC. Modul ini sediakan kolomnya. FE kirim PT/cabang aktif (hasil pilih di workspace switcher) via header `X-Company-Id`/`X-Branch-Id`, divalidasi server-side terhadap assignment asli user — bukan query param, lihat [scoping-convention.md](../../.agents/rules/scoping-convention.md) §3.1.

## 4. Kriteria Penerimaan (Given-When-Then)

- **Department hierarki valid.**
  *Given* `parent_id` menunjuk Department lain di Company sama, *When* create, *Then* sukses. *Given* `parent_id = null`, *Then* jadi root.

- **Department milik satu Company.**
  *Given* payload Department tanpa `company_id` valid, *When* create, *Then* tolak `422`.
  *Catatan implementasi:* entity Department existing **belum punya** `company_id` — gap yang ditutup saat implementasi split (butuh migrasi + ubah domain).

- **Job Position wajib Department + Job Title valid.**
  *Given* salah satu FK tak dikenal, *When* create Position, *Then* tolak `422`.

- **Reporting line tak lintas Company.**
  *Given* `reports_to_id` menunjuk Position milik Company lain, *When* set, *Then* tolak `422` (reporting line dalam satu PT).

- **Headcount quota terdefinisi.**
  *Given* create Position, *When* tanpa `headcount_quota`, *Then* default sesuai aturan (mis. 1) atau tolak — dipertegas di tech-spec.
  *Catatan:* enforcement "assign tak boleh lebihi quota" ada di modul Employee, bukan di sini.

- **Reporting line anti-siklus.**
  *Given* set `reports_to_id` yang bikin loop (A→B→A), *When* simpan, *Then* tolak (cycle detection).

## 5. Technical & Architectural Constraints
- **Bounded context baru:** `workforce-structure` (folder kode `internal/workforce/` atau setara — package Go tanpa tanda hubung). Existing dept/title/position code **dipindah** ke sini dari `internal/domain/organization/`.
- **Depend ke Organization:** Department punya `company_id` FK → `organization.Company`. Cross-domain komunikasi via **Application Service** Organization (bukan inject repo langsung) — [coding-convention.md](../../.agents/rules/coding-convention.md) §4.
- **`company_id` NOT NULL** di Department (turun ke Position lewat Department). Migrasi + backfill wajib.
- **Sentinel error:** `ErrDepartmentNotFound`, `ErrJobTitleNotFound`, `ErrJobPositionNotFound`, `ErrReportingCycle` — [persistence-convention.md](../../.agents/rules/persistence-convention.md) §3.
- **UUID** generate di domain constructor — [uuid-generation.md](../../.agents/rules/uuid-generation.md).
- **Tier dokumen:** **Sedang** (relasi antar-entity + hierarki + depend Organization) → butuh `tech-spec.md`. `decision-log.md` kalau ada keputusan non-trivial (mis. algoritma cycle detection). `user-stories.md` bukan wajib tier ini, tapi dibuat juga atas permintaan eksplisit — dua entity self-referencing (Department/JobPosition) dianggap cukup berisiko buat divisualisasikan lewat sequence diagram.

## 6. Dependencies
- **Depends on:** `organization@2.0.0` — konsumsi `company_id` untuk scope Department (Organization §Company).
- **Consumed by:** `employee` — Employee menduduki Job Position; konsumsi `job_position_id`, reporting line, headcount.
- **External:** —

---

## 7. Data Schema & Business Rules

> Sample buat FE, bukan pengganti DBML. Kolom `company_id` ditambahkan saat split (planned).

### 7.1. Department (`departments`)
Aturan: `parent_id = null` → root; hierarki tree; wajib `company_id`.

| id | company_id | code | name | parent_id | is_active |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `dept-1` | `co-1` | DIR | Direksi | `null` | true |
| `dept-2` | `co-1` | TI | Divisi Teknologi Informasi | `dept-1` | true |
| `dept-3` | `co-1` | DEV | Departemen Pengembangan | `dept-2` | true |
| `dept-4` | `co-1` | OPR | Divisi Operasional | `dept-1` | true |
| `dept-5` | `co-1` | SDM | Divisi Sumber Daya Manusia | `dept-1` | true |

### 7.2. Job Title (`job_titles`)
Aturan: `grade_level` makin tinggi = pangkat makin tinggi; independen dari departemen tapi **per-PT** (Company-owned, wajib `company_id` — tiap PT punya grade/pangkat sendiri, [scoping-convention.md](../../.agents/rules/scoping-convention.md) §1). Code unik dalam satu Company.

| id | company_id | code | name | grade_level | is_active |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `title-1` | `co-1` | DIR | Direktur | 10 | true |
| `title-2` | `co-1` | KDV | Kepala Divisi / GM | 9 | true |
| `title-3` | `co-1` | MGR | Manajer | 7 | true |
| `title-4` | `co-1` | SPV | Supervisor | 5 | true |
| `title-5` | `co-1` | STF | Staf | 3 | true |

### 7.3. Job Position (`job_positions`)
Aturan: kombinasi 1 Department + 1 Job Title; `reports_to_id` → Position atasan (org chart); `headcount_quota` batas pegawai.

| id | name | department_id | job_title_id | reports_to_id | headcount_quota |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `pos-1` | Direktur Utama | `dept-1` | `title-1` | `null` | 1 |
| `pos-2` | Direktur Teknologi (CTO) | `dept-2` | `title-1` | `pos-1` | 1 |
| `pos-3` | Kepala Divisi TI | `dept-2` | `title-2` | `pos-2` | 1 |
| `pos-4` | Manajer Pengembangan | `dept-3` | `title-3` | `pos-3` | 3 |
| `pos-5` | Supervisor Backend | `dept-3` | `title-4` | `pos-4` | 5 |
| `pos-6` | Staf Programmer Backend | `dept-3` | `title-5` | `pos-5` | 10 |

---

## Catatan untuk Frontend (FE)

1. **Form Create Job Position:** butuh dropdown `Department` (`GET /workforce/departments`) & `Job Title` (`GET /workforce/job-titles`). `Reports To` = autocomplete dari `GET /workforce/job-positions`.
2. **Organization Chart:** render `job_positions` via `reports_to_id`. Root = `reports_to_id = null`.
3. **`is_active`:** default API balikin yang aktif; yang nonaktif bisa di-grey-out/sembunyikan.
4. **Search:** ketiga endpoint List (`departments`, `job-titles`, `job-positions`) dukung query param `search` opsional — match `code` ATAU `name` (Department/Job Title), `name` saja (Job Position, gak punya `code`). Case-insensitive substring, dipakai buat autocomplete di poin 1.
5. **Workspace switcher (PT/cabang aktif):** kirim header `X-Company-Id`/`X-Branch-Id` sesuai pilihan switcher, bukan query param — lihat catatan §3.
6. **Department Tabel (nested row) & Bagan (tree diagram):** pakai `GET /workforce/departments/tree` (TANPA pagination, beda dari `GET /workforce/departments` yang paginated) — FE assembly tree dari `parent_id` sendiri (group by `parent_id`, root = `parent_id: null`), dipakai buat expand/collapse row bertingkat maupun render diagram. Pagination server-side gak cocok buat expand/collapse (parent & anak bisa kepisah halaman), jadi endpoint ini sengaja unpaginated + client-side pagination di FE.
7. **`department`/`job_title` di response Job Position:** bukan lagi `department_id`/`job_title_id` string doang — sekarang objek `{id, name}` (preload, hindari FE nampilin raw UUID di tabel/form). `id` di dalamnya tetap dipakai buat submit form edit.

> Path endpoint di atas ilustratif (`/workforce/...`) — final URL ditetapkan saat scaffold API docs.
