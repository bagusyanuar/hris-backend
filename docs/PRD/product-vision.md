---
doc_type: Product Vision (Global PRD / North Star)
version: 1.0.0
status: Draft
owner: bagusyanuar
updated: 2026-07-23 13:35:51
scope: platform-wide          # bukan PRD per-bounded-context; ini payung semua modul
---

# HRIS Backend — Global PRD (Product Vision / North Star)

Dokumen ini adalah **payung** di atas semua PRD per-modul. Fungsinya menetapkan *arah produk* (GOALS), *model deployment*, dan *prinsip arsitektur global* yang **tidak boleh dilanggar** oleh PRD modul manapun. Kalau ada konflik antara PRD modul dan dokumen ini, dokumen ini menang (atau dokumen ini yang di-revisi lebih dulu secara sadar).

> Beda peran dokumen:
> - **Global PRD (ini)** = arah & aksioma seluruh platform.
> - **PRD modul** (`organization.md`, `employee.md`, …) = WHAT/WHY per bounded context.
> - **tech-spec / decision-log** = HOW teknis per modul.

---

## 1. Visi Produk (North Star)

> **Satu platform HRIS untuk satu grup usaha (holding) yang menaungi beberapa PT, masing-masing dengan banyak cabang — dikelola satu owner, dari satu login, dengan visibilitas konsolidasi penuh atas seluruh grup.**

Sistem ini melayani **enterprise skala menengah–besar** dengan struktur kompleks (multi-PT, multi-cabang, multi-departemen) dan jumlah karyawan yang bertumbuh. Arsitektur DDD dipilih supaya tiap domain bisa berkembang mandiri, bahkan diekstrak jadi service terpisah nanti.

---

## 2. Model Bisnis & Deployment (AKSIOMA — mengunci semua desain)

Keputusan paling fundamental yang menyetir seluruh arsitektur:

| Aspek | Keputusan | Konsekuensi |
|-------|-----------|-------------|
| **Model** | **Group / Holding, single-owner** | Satu instalasi = satu grup usaha milik satu owner. |
| **BUKAN** | **BUKAN multi-tenant SaaS** | Tidak dijual ke banyak klien tak-saling-kenal. Tidak ada lapis `Tenant`, tidak ada onboarding self-service, tidak ada billing per-tenant, tidak ada sub-domain per-tenant. |
| **Legal entity** | **Banyak `Company` (PT)** di bawah satu owner | Tiap PT punya NPWP/BPJS/payroll sendiri. Owner lihat konsolidasi lintas PT. |
| **Lokasi** | **Banyak `Branch` (cabang)** per Company | Operasional (absensi/shift/UMR) di-scope per cabang. |
| **Isolasi data** | **Shared DB + row-level scoping** (`company_id` / `branch_id`) | Cukup untuk group tunggal. TIDAK pakai schema-per-tenant / DB-per-tenant kecuali ada tuntutan regulasi nyata. |

> **Kalau suatu hari model berubah jadi SaaS** — itu perubahan MAJOR pada dokumen ini, wajib tambah lapis `Tenant` di atas `Company` dan review ulang strategi isolasi. Jangan diam-diam nyelundupin asumsi SaaS ke modul.

---

## 3. Hierarki Struktural Global (Kanonik)

Semua modul WAJIB mengacu ke hierarki ini. Company/Branch didefinisikan di [organization.md](organization.md); Department/Position di [workforce-structure.md](workforce-structure.md).

```text
Group / Holding (implisit = 1 instalasi aplikasi, 1 owner)
  └── Company (PT / badan hukum)        ← company_id  · NPWP, BPJS, payroll, pajak   [Organization]
        └── Branch (cabang / lokasi)    ← branch_id   · absensi, shift, UMR, libur   [Organization]
              └── Department            ← struktur unit kerja (tree)                 [Workforce Structure]
                    └── Job Position    ← "kursi" (Department × Job Title)           [Workforce Structure]
                          └── Employee  ← menduduki posisi; wajib company_id + branch_id  [Employee]
```

**Dua dimensi scoping wajib** (jangan digabung jadi satu kolom):
- `company_id` → dimensi **legal** (payroll, pajak, kontrak). Non-nullable di semua entity operasional.
- `branch_id` → dimensi **lokasi** (operasional harian). Non-nullable di Employee & transaksi operasional.

---

## 4. Prinsip Arsitektur Global (Non-Negotiable)

1. **DDD domain-first** — tiap bounded context = satu folder utuh di `internal/`, siap diekstrak jadi service. Lihat [architecture.md](../../.agents/rules/architecture.md).
2. **Loose coupling di dokumen & kode** — modul merujuk field/section modul lain (+ versi), bukan copy-paste aturan. Konsep lintas-modul naik ke [_shared/glossary.md](_shared/glossary.md).
3. **Concern lintas-modul = modul sendiri** — RBAC, audit trail, notifikasi TIDAK ditempel ke modul terdekat. Masing-masing PRD sendiri.
4. **Dua dimensi scoping (`company_id`/`branch_id`) sejak baris pertama** — kolom non-nullable dari migrasi awal. Retrofit belakangan = neraka migrasi.
5. **Payroll/pajak selalu per-Company** — tidak ada agregasi lintas PT dalam satu slip/laporan pajak. Konsolidasi hanya di layer reporting owner.
6. **Cross-domain via Application Service** (sinkron, Wire DI) — bukan inject repository modul lain, bukan message broker (untuk sekarang). Lihat [coding-convention.md](../../.agents/rules/coding-convention.md) §4.

---

## 5. Peta Modul & Roadmap

### 5.1. Modul Fondasi (sudah/sedang berjalan)
| Modul | Peran | Status kode |
|-------|-------|-------------|
| **Auth** | autentikasi & access control dasar | ada |
| **User** | akun pengguna sistem | ada |
| **Organization** | legal & lokasi: **Company (PT), Branch** | **planned** ([organization.md](organization.md)) |
| **Workforce Structure** | struktur internal: Department, Job Title, Job Position | 3 pilar ada di kode; **pindah dari Organization** (planned) ([workforce-structure.md](workforce-structure.md)) |
| **Employee** | data & profil karyawan | ada; **belum punya `company_id`/`branch_id`** (gap) |

### 5.2. Concern Lintas-Modul (jadi modul/PRD sendiri)
| Modul | Peran | Kenapa dipisah |
|-------|-------|----------------|
| **RBAC** | enforce scoping `company_id`/`branch_id` + role/permission | dikonsumsi SEMUA modul; fondasi access-control |
| **Audit Trail** (nanti) | jejak siapa-ubah-apa-kapan lintas modul | requirement enterprise |
| **Notification** (nanti) | email/push lintas modul | dikonsumsi banyak modul |

### 5.3. Modul HRIS Inti (arah pengembangan, belum PRD)
Attendance & Time Tracking · Leave/Time-off · **Payroll & Compensation** · Performance Management · Recruitment & Onboarding.

> Semua modul di §5.3 WAJIB lewat proses PRD/Tech-Spec penuh sebelum diimplementasi. Payroll & Attendance = tier **Kompleks** (kalkulasi berlapis / state machine).

### 5.4. Urutan Eksekusi yang Disarankan (fase besar)
1. **Fase 1 — Fondasi Multi-Entity.** Company + Branch (CRUD murni) di modul Organization. Aman, tak breaking. → [organization.md](organization.md).
2. **Fase 2 — Split Workforce Structure.** Pindah Dept/Title/Position dari `organization` ke modul `workforce-structure` + tambah `company_id`. Breaking (rename + migrasi). → [workforce-structure.md](workforce-structure.md).
3. **Fase 3 — Scoping Employee.** Migrasi `company_id`/`branch_id` ke Employee + validasi cross-entity. Breaking, butuh backfill.
4. **Fase 4 — RBAC.** PRD + modul enforce scoping. Track terpisah, paling berat.
5. **Fase 5 — Modul HRIS inti** (Attendance → Leave → Payroll → dst), tiap-tiap lewat PRD sendiri, semua sadar dua dimensi scoping.

---

## 6. Non-Goals (Batas Tegas — mencegah scope creep)

- ❌ **BUKAN SaaS multi-tenant.** Tidak melayani banyak owner/klien tak-saling-kenal. Tidak ada lapis Tenant, billing, self-service onboarding.
- ❌ **Bukan** konsolidasi payroll lintas-PT dalam satu perhitungan pajak. Tiap PT tutup buku sendiri.
- ❌ **Bukan** schema-per-tenant / DB-per-tenant (sampai ada tuntutan regulasi nyata).
- ❌ **Bukan** message broker / event-driven antar-modul untuk sekarang (sinkron via Application Service).
- ❌ **Bukan** GORM `AutoMigrate` di produksi — skema via SQL migration + DBML.

---

## 7. Metrik Sukses (indikatif — dipertajam per modul)

- Owner bisa lihat **konsolidasi lintas-PT** (jumlah karyawan, headcount) dari satu dashboard.
- Data satu cabang/PT **tidak bocor** ke cabang/PT lain (enforced RBAC) — 0 insiden kebocoran lintas-scope.
- Tambah PT/cabang baru **tanpa deploy ulang** aplikasi — cukup master data.
- Tiap PT bisa jalankan payroll/pajak **independen** tanpa saling ganggu.

---

## 8. Referensi

- Index PRD: [README.md](README.md)
- Glossary lintas-modul: [_shared/glossary.md](_shared/glossary.md)
- Rules arsitektur & konvensi: [.agents/rules/](../../.agents/rules/)
- PRD fondasi multi-entity: [organization.md](organization.md) §0
