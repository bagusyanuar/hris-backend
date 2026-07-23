# Architecture Decision Records (ADR): Organization Module (v2.0.0)

## ADR-001: Company & Branch sebagai Dua Kolom Scope Terpisah (`company_id` + `branch_id`)
- **Date:** 2026-07-23
- **Status:** Accepted
- **Context:** Payroll/pajak (PPh 21, BPJS, NPWP) legal-nya diikat ke badan hukum (PT), sementara jam kerja/shift/UMR diikat ke lokasi fisik (cabang). Kalau digabung jadi satu dimensi (mis. hanya `branch_id`, company diturunkan dari branch), query lintas-cabang dalam satu PT (mis. laporan payroll konsolidasi per PT) jadi butuh JOIN tambahan di setiap tempat, dan entity yang company-scoped-tapi-bukan-lokasi-spesifik (mis. Department di level PT) jadi terpaksa nempel ke satu cabang secara salah kaprah.
- **Decision:** Dua kolom scope independen: `company_id` (WAJIB di semua entity operasional) + `branch_id` (WAJIB hanya untuk entity yang benar-benar lokasi-spesifik). Lihat klasifikasi kelas di [scoping-convention.md](../../../.agents/rules/scoping-convention.md) §1.
- **Consequence:** Entity company-owned-tapi-bukan-lokasi-spesifik (Department, Job Title) tidak perlu `branch_id` sama sekali — hemat kolom, hindari ambiguitas "cabang mana" untuk data yang bisnisnya memang per-PT.

## ADR-002: Row-Level Scoping (Shared DB + Kolom), bukan Schema/DB-per-Tenant
- **Date:** 2026-07-23
- **Status:** Accepted
- **Context:** Multi-PT bisa diisolasi dengan 3 pendekatan: (a) DB terpisah per PT, (b) schema Postgres terpisah per PT, (c) shared table + kolom `company_id`. Opsi (a)/(b) memberi isolasi lebih kuat tapi migration/maintenance overhead naik linear dengan jumlah PT (owner bisa punya >10 PT dalam grup usaha), dan query konsolidasi lintas-PT (kebutuhan Owner/Group Admin) jadi butuh cross-database query yang jauh lebih mahal.
- **Decision:** Shared database + kolom `company_id`/`branch_id` row-level (per [product-vision.md](../../PRD/product-vision.md) §2 dan [organization.md](../../PRD/organization.md) §5). Isolasi ditegakkan di query boundary (aplikasi), bukan di level infrastruktur DB.
- **Consequence:** Enforcement isolasi jadi tanggung jawab kode (RBAC middleware + `scope.FromContext` di setiap repository read) — bukan otomatis dijamin DB. Kalau suatu saat ada tuntutan regulasi yang butuh isolasi fisik penuh (jarang), perlu migrasi arsitektur terpisah — bukan default sekarang.

## ADR-003: Branch adalah Aggregate Root Sendiri, Bukan Child Entity Company
- **Date:** 2026-07-23
- **Status:** Accepted
- **Context:** Kalau Branch dimodelkan sebagai child collection yang selalu dimuat bersama Company (mis. `Company.Branches []Branch`), setiap load Company jadi query mahal + N+1 risk ketika Company punya puluhan cabang. Selain itu, operasi Branch (create/update/delete satu cabang) tidak butuh memuat ulang seluruh Company.
- **Decision:** Branch = aggregate root independen dengan `CompanyID` sebagai foreign key referensi (bukan embedded slice). `BranchRepository` terpisah dari `CompanyRepository`.
- **Consequence:** Konsistensi antar aggregate (mis. "Company tidak boleh dihapus kalau masih punya Branch aktif") jadi tanggung jawab Application Service, bukan otomatis dijamin satu aggregate boundary. Untuk scope 2.0.0 aturan itu belum diimplementasikan (lihat tech-spec.md §6.1 poin 5) — Company delete tidak mem-validasi Branch anak dulu, ditandai sebagai gap eksplisit, bukan silent decision.

## ADR-004: `is_main` Branch — Demote Otomatis, Bukan Tolak (Reject)
- **Date:** 2026-07-23
- **Status:** Accepted
- **Context:** PRD §4 acceptance criteria: *"Given sudah ada Branch is_main=true di Company X, When set Branch lain jadi is_main di Company X, Then pindahkan status / tolak sesuai aturan tech-spec (satu main per Company)"* — PRD sengaja menyerahkan pilihan pindahkan-vs-tolak ke tech-spec ini. Dari sisi UX, "tolak dengan error" berarti Admin harus manual: (1) unset main branch lama, (2) baru set main branch baru — dua request terpisah dengan window race (sempat tidak ada main branch sama sekali, atau race dua request paralel bikin dua main branch).
- **Decision:** Demote otomatis. Saat create/update Branch dengan `is_main=true`, dalam satu `TxManager.Do`: (1) `DemoteMainBranch(companyID)` — set `is_main=false` untuk main branch lama di company yang sama, (2) simpan branch baru/ubah dengan `is_main=true`. Partial unique index `idx_branches_company_main` di DB tetap jadi jaring pengaman terakhir kalau ada bug di application layer yang lewatkan step demote.
- **Consequence:** Admin tidak perlu dua langkah manual. Trade-off: aksi "pindahkan kantor pusat" jadi implicit side-effect dari update biasa (tidak ada endpoint terpisah `PATCH /branches/{id}/set-main`) — didokumentasikan di sini supaya FE tahu efek sampingnya saat toggle `is_main`.

## ADR-005: `npwp` & `bpjs_no` Nullable (Bukan NOT NULL) di Scope Awal
- **Date:** 2026-07-23
- **Status:** Accepted
- **Context:** PRD §7.1 mendeskripsikan `npwp`/`bpjs_no` sebagai bagian data legal Company, tapi belum ada modul konsumen (Payroll) yang benar-benar butuh nilainya sekarang — mewajibkan `NOT NULL` di awal berarti tim harus punya data NPWP/BPJS valid untuk *setiap* PT sebelum bisa input Company sama sekali, padahal onboarding PT baru sering data itu menyusul belakangan.
- **Decision:** Kolom nullable, unique constraint tetap ada (partial index `WHERE npwp IS NOT NULL`, aman untuk multi-NULL) — bukan penghapusan field. Keputusan tim, dikonfirmasi 2026-07-23.
- **Consequence:** Validasi "npwp wajib diisi" (kalau suatu saat dibutuhkan Payroll) jadi tanggung jawab modul consumer atau validasi tambahan di masa depan, bukan constraint DB. `ErrCompanyNPWPDuplicate` di PRD §4 tetap berlaku HANYA ketika `npwp` diisi (dua Company boleh sama-sama `npwp = NULL`).

## ADR-006: Nested `branches` di `GET /companies` — Batch Query Manual, Bukan GORM `Preload`
- **Date:** 2026-07-23
- **Status:** Accepted
- **Context:** FE butuh tampilkan Company + daftar Branch-nya sekaligus di satu halaman (list view), tanpa round-trip kedua ke `GET /companies/{id}/branches` per row. ADR-003 sudah tetapkan Branch = aggregate root terpisah (bukan child collection Company) justru untuk hindari pola "selalu ikut ke-load bareng Company" yang beresiko N+1. Kebutuhan UI ini sekilas kontradiksi ADR-003, jadi perlu didokumentasikan eksplisit kenapa TIDAK ubah aggregate boundary.
- **Decision:**
  1. Nested `branches` di response `GET /companies` adalah **komposisi read-model di Application Layer**, bukan perubahan aggregate — `BranchRepository` tetap terpisah dari `CompanyRepository`, tidak ada foreign-key embed di domain entity `Company`.
  2. Diambil lewat method baru `BranchRepository.FindAllByCompanyIDs(ctx, companyIDs []string)` — SATU query `WHERE company_id IN (...)` atas seluruh company di halaman itu (bukan loop per company / N+1), lalu di-group manual per `company_id` di Application service.
  3. **Bukan GORM `Preload`** meski secara mekanisme sama-sama batch query di balik layar — alasannya bukan performa (sama), tapi kontrol: `Preload` butuh association field (`Company.Branches []BranchModel`) nempel di GORM model, yang bikin adapter model "tau" relasi lintas aggregate (nyimpang dari ADR-003), dan susah disisipin filter `scope.FromContext` custom nanti pas RBAC landing (scoping-convention.md §3) dibanding query manual yang eksplisit di adapter.
  4. Saat `search` match branch name (bukan `legal_name`), `branches` yang di-embed tetap FULL LIST milik company itu — TIDAK difilter cuma yang match. `search` cuma nentuin company mana yang lolos filter, bukan menyaring isi nested-nya (behavior konsisten, gampang didokumentasikan ke FE, tanpa `CompanyResponse` jadi punya dua bentuk berbeda tergantung ada-tidaknya `search`).
- **Consequence:** `GET /companies` sekarang selalu jalankan 1 query tambahan (branch batch) per request list, walau FE kadang tidak butuh data branch-nya (trade-off diterima, volume data kecil — lihat tech-spec.md §8). Kalau nanti butuh skip branch (mis. modul lain consume endpoint ini tanpa perlu nested), pertimbangkan `?include=branches` opsional saat itu — belum dibutuhkan sekarang (YAGNI).
