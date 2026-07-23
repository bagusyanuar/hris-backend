---
name: scaffold-dbml
description: Guide for generating or updating the DBML database schema doc for a module, sourced from its approved PRD, as the single source of truth for SQL migrations.
---

# Scaffolding DBML (Database Schema Documentation)

Trigger this skill whenever the user references an existing PRD (`docs/PRD/<module>.md`) and asks to create/update the database schema — e.g. "buatkan skema database nya", "generate DBML", "bikin schema buat migrasi".

## 1. Non-Negotiable Rule (see [project-docs.md](../../rules/project-docs.md) §3)
DBML is **MANDATORY for every module**, regardless of tier (Simpel/Sedang/Kompleks). It is the **single source of truth for SQL migrations** — `AutoMigrate` is forbidden in production (see [architecture.md](../../rules/architecture.md) §C), so the physical schema must be pinned explicitly in DBML before/alongside any migration file.

## 2. Location & Naming Convention
- **Location:** `docs/databases/<module>.dbml` — same base name as the module's PRD file (`docs/PRD/<module>.md` → `docs/databases/<module>.dbml`).
- **Anti-duplication:** Before creating a new file, `ls docs/databases/` to confirm one doesn't already exist for this module. If it exists, **update in place** — do not create a second file with a different name pattern for the same module.
- **Table naming:** snake_case, plural (e.g. `users`, `job_positions`), matching the real/planned SQL table name — not the Go entity name.

## 3. Source of Truth Priority
1. **Section 7 ("Data Schema & Business Rules") of the module's PRD** — defines entities, fields, and business rules to translate into DBML.
2. **Existing migration file** (`migrations/*.up.sql`), if one already exists for this domain — the DBML **MUST match the real migration exactly** (types, constraints, defaults). Never invent a DBML schema that diverges from what's already migrated; if they disagree, the migration is the ground truth and the DBML must be corrected to match it (or vice versa, if the PRD's Acceptance Criteria demand a change — flag this explicitly, don't silently pick one).
3. If **no migration exists yet**, the DBML is the spec the future migration (`make migrate-create`) should be authored from.

## 4. Writing the DBML
- Reflect exact column types/constraints as they are (or will be) in Postgres: `uuid`, `varchar(n)`, `not null`, `default`.
- **DBML cannot express every real-world SQL constraint** (partial unique indexes, `WHERE`-conditioned constraints, `CHECK` enums, composite uniqueness). When the real/planned constraint is more specific than what DBML's shorthand (`unique`, `pk`) can say, **do not use the flat shorthand alone** — annotate with `note:` explaining the actual rule and point to the migration file. (Example: `docs/databases/user.dbml` — `email` uses a partial unique index `WHERE deleted_at IS NULL`, not a plain unique constraint, because emails are reusable after soft-delete.)
- Add `Ref:` blocks for foreign keys to other domains' tables (e.g., `workforce_structure`'s `job_positions.department_id` → `departments.id`).
- Cross-check field name casing against [coding-convention.md](../../rules/coding-convention.md) §6 (acronym consistency) so DB column, GORM model, and Domain Entity don't drift (e.g. `ktp_number` vs `KtpNumber`, not mixed casing).

### 4a. Kolom Scope Multi-Entity (WAJIB — [scoping-convention.md](../../rules/scoping-convention.md))
Setiap tabel entity operasional WAJIB bawa kolom scope sesuai kelasnya (scoping-convention.md §1):
- **Company-owned** (default) → `company_id uuid [not null]` + `Ref: <table>.company_id > companies.id`.
- **Company + Location** → tambah `branch_id uuid [not null]` + `Ref: <table>.branch_id > branches.id`.
- **Index wajib** di kolom scope: `Indexes { company_id }` (dan `(company_id, branch_id)` untuk entity dua-dimensi) — semua query difilter lewat kolom ini.
- **Integritas silang** `branch_id` se-`company_id` tak bisa diekspresikan DBML polos → kasih `note:` yang rujuk aturan + migration.
- **Staged (§4 scoping-convention):** tabel `companies`/`branches` mungkin belum ada saat DBML modul lain ditulis. Tetap **deklarasikan** `Ref:`-nya (dokumentasi kontrak); catat di akhir bahwa migrasi FK harus diurutkan setelah migrasi Organization.

## 5. What This Skill Does NOT Do
- **Does not write or run SQL migrations.** DBML is documentation/spec only. Creating the actual `migrations/*.up.sql`/`*.down.sql` pair is a separate, explicit step via `make migrate-create` — only do it if the user asks for the migration itself, not just the schema doc.
- **Does not scaffold Go code** (entity/model/repository). That's `scaffold-domain` / `execute-domain`'s job.
- If the DBML introduces a field that has no corresponding migration yet, say so explicitly at the end: *"Field X belum ada di migration existing — perlu `make migrate-create` dulu sebelum bisa dipakai."*

## 6. Checklist Before Reporting Done
- [ ] File di `docs/databases/<module>.dbml`, bukan file baru kalau udah ada.
- [ ] Setiap field balik ke PRD §7 module bersangkutan — gak ada field nebak-nebak yang gak dijustifikasi di PRD.
- [ ] Constraint yang gak bisa direpresentasikan shorthand DBML (partial index, CHECK enum, composite unique) diberi `note:`, bukan diam-diam disederhanakan jadi `[unique]` polos.
- [ ] Kalau ada migration existing, DBML match persis — bukan versi idealis yang beda dari real schema.
- [ ] Foreign key ke domain lain pakai `Ref:`.
- [ ] Kolom scope (`company_id`/`branch_id`) ada sesuai kelas entity + ber-index + `Ref:` ke `companies`/`branches` ([scoping-convention.md](../../rules/scoping-convention.md) §1–§2). Kalau "Global master" tanpa scope → ada justifikasi eksplisit.
