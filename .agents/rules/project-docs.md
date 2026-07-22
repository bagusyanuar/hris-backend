# Dokumen Proyek (Requirements & Technical)

Ada dua jenis dokumen yang wajib disimpan ke dalam repositori secara permanen:

1. **Product Requirements Document (PRD):**
   - Jika user meminta rancangan fitur, requirement, atau skema bisnis (bukan teknis database murni), gunakan format PRD.
   - **WAJIB** disimpan di folder `docs/requirement/` (contoh: `docs/requirement/employee.md`).
   - Format penulisan wajib mematuhi panduan dari skill `scaffold-prd`.
2. **Dokumen Teknis & Arsitektur (Enterprise Tech Specs):**
   - Setiap rancangan arsitektur sistem, struktur database, atau *implementation plan* murni teknis, **WAJIB** disimpan di dalam sub-folder per domain di `docs/technical/<domain_name>/`.
   - Dokumentasi di dalam folder tersebut **harus dipecah** menjadi beberapa file spesifik:
     - `tech-spec.md` (Arsitektur inti, API, dan skema DB).
     - `user-stories.md` (Alur logika dan diagram *sequence*).
     - `decision-log.md` (ADR - Mencatat *kenapa* keputusan teknis tertentu diambil).
     - Serta dokumen pendukung opsional seperti `data-dictionary.md`, `infrastructure.md`, dan `test-plan.md`.
3. **Database Markup (DBML):**
   - Skema database relasional **WAJIB** ditulis dalam format DBML (`.dbml`) dan disimpan di folder `docs/databases/` (contoh: `docs/databases/employee.dbml`).
   - Tujuannya agar arsitektur ERD bisa divisualisasikan dengan mudah via dbdiagram.io.
