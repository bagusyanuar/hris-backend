# Dokumen Proyek (Requirements & Technical)

Ada dua jenis dokumen yang wajib disimpan ke dalam repositori secara permanen:

1. **Product Requirements Document (PRD):**
   - Jika user meminta rancangan fitur, requirement, atau skema bisnis (bukan teknis database murni), gunakan format PRD.
   - **WAJIB** disimpan di folder `docs/PRD/` (contoh: `docs/PRD/employee.md`).
   - Format penulisan wajib mematuhi panduan dari skill `scaffold-prd`.
   - **Module boundary heuristic:** Concern lintas-modul yang bakal dikonsumsi *banyak* bounded context masa depan (mis. authorization/roles, audit trail, notifikasi) **WAJIB** jadi PRD/modul sendiri — jangan ditempel ke modul terdekat cuma karena kebetulan berhubungan (contoh: RBAC jadi PRD sendiri, bukan ditempel ke User cuma karena User yang punya tabel akun). Loose coupling di dokumen, sama seperti prinsip di kode.
2. **Dokumen Teknis & Arsitektur (Enterprise Tech Specs) — BERTINGKAT (Tiered):**
   - Dokumen teknis disimpan di sub-folder per domain di `docs/technical/<domain_name>/`. **Kelengkapannya menyesuaikan kompleksitas modul** — bukan all-or-nothing. Tujuannya: cepat untuk modul remeh, punya jaring pengaman kontrak untuk modul berisiko.
   - **Prinsip:** PRD = WHAT/WHY (bisnis). Technical = HOW (kontrak API, DDL, sequence, ADR). Jangan campur.
   - **Penentuan tingkat** dilakukan setelah PRD di-approve (lihat workflow `scaffold-docs`). Tiga tingkat:

     | Tingkat | Kriteria modul | Dokumen teknis WAJIB |
     |---------|----------------|----------------------|
     | **Simpel** | CRUD lurus, 1–2 entity, tanpa integrasi luar, tanpa kalkulasi/state machine (mis. master data, lookup) | *Tidak ada tech-spec.* Cukup PRD + DBML. Boleh langsung scaffold kode dari PRD. |
     | **Sedang** | Ada relasi antar-entity, state/status flow, atau depend antar-modul | `tech-spec.md` (arsitektur inti + kontrak API) |
     | **Kompleks** | Kalkulasi berlapis, state machine, atau integrasi eksternal (mis. Payroll, Attendance, Leave) | Set penuh: `tech-spec.md` + `user-stories.md` + `decision-log.md` |

   - Dokumen pendukung opsional (`data-dictionary.md`, `infrastructure.md`, `test-plan.md`) ditambahkan hanya bila modul kompleks menuntut.
   - **`decision-log.md` (ADR)** wajib begitu ada keputusan teknis non-trivial yang perlu dijelaskan *kenapa*-nya — meskipun modul tergolong sedang. PRD tidak merekam alasan teknis; jangan hilangkan konteks ini.
3. **Database Markup (DBML) — WAJIB SEMUA MODUL (non-negotiable):**
   - Skema database relasional **WAJIB** ditulis dalam format DBML (`.dbml`) dan disimpan di folder `docs/databases/` (contoh: `docs/databases/employee.dbml`), **tanpa terkecuali tingkat kompleksitas apa pun**.
   - DBML adalah **sumber tunggal migrasi SQL**. "Sample data" di PRD bukan skema fisik dan **tidak** menggantikan DBML. Ingat: `AutoMigrate` dilarang di produksi (lihat [architecture.md](architecture.md) §C), jadi skema fisik harus dipaku eksplisit di DBML.
   - Tujuannya agar arsitektur ERD bisa divisualisasikan dengan mudah via dbdiagram.io.
   - **Trigger:** Kalau user merujuk PRD existing lalu minta "buatkan skema database nya" (atau sejenisnya), gunakan skill `scaffold-dbml` — mendefinisikan lokasi file, konvensi penamaan, dan cara menerjemahkan §7 PRD (+ migration existing kalau ada) jadi DBML yang akurat.
