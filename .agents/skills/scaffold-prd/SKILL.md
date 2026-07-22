---
name: scaffold-prd
description: Guide for scaffolding Product Requirements Documents (PRD) in the HRIS project to ensure consistency between Backend and Frontend teams.
---

# Scaffolding PRD (Product Requirements Document)

When the user asks to create or brainstorm a requirement document, PRD, or specifications for a new module/feature, you MUST follow this format and convention.

## 1. Directory & Naming Convention
- **One file per module (bounded context).** DO NOT put multiple modules in a single monolithic file. Each module gets its own `.md`, mirroring the domain-first code layout.
- **Location:** All PRD documents MUST be saved in the `docs/PRD/` directory.
- **Naming Convention:** Use lowercase with hyphens (e.g., `employee.md`, `attendance-tracking.md`).
- **Do not** use `docs/technical/` for PRDs. `docs/technical/` is strictly for implementation plans and system architecture documents.
- **Index maintenance (MANDATORY):** Every time a PRD is created or its status/version changes, you MUST update `docs/PRD/README.md` (module table + dependency graph). The index is the single "see everything at a glance" entry point.
- **Shared concepts:** Cross-module terminology or concepts (e.g., "Working Day", holiday calendar) live in `docs/PRD/_shared/glossary.md`. Modules **reference** the glossary — never copy-paste shared rules into individual PRDs.
- **Module boundary heuristic:** A cross-cutting concern that will be consumed by *many* future bounded contexts (e.g., authorization/roles, audit trail, notifications) gets its **own** PRD/module — do not fold it into the nearest existing module just because it's convenient right now (e.g., RBAC belongs in its own PRD, not bolted onto User just because User owns the account table). Loose coupling in docs mirrors loose coupling in code.
- **Ground truth before writing (MANDATORY for already-implemented modules):** If the module already has code (`internal/<module>/`), you MUST read the real `domain`/`application`/`infrastructure`/`transport` files before writing or updating a single line of the PRD. Never infer behavior from function/file names alone. Acceptance Criteria and Constraints must reflect what the code *actually does*, not what it's assumed to do.

## 1a. Mandatory Frontmatter Header
Every PRD file MUST begin with a YAML frontmatter block for versioning and traceability:

```yaml
---
module: Payroll
version: 1.0.0          # SemVer — bump on every business-rule change
status: Draft           # Draft | In Review | Approved | Deprecated
owner: <name>
updated: 2026-07-22 14:35:07   # yyyy-MM-dd HH:mm:ss, Asia/Jakarta timezone, down to the second (not date-only)
depends_on: [attendance@1.2, leave@1.0]   # empty [] if none
---
```

Bump `version` (SemVer) and refresh `updated` (full timestamp, not just the date) whenever the PRD content changes, and reflect it in the README index. Concrete bump criteria:
- **PATCH** (1.0.0→1.0.1): wording clarification, typo fix — no change in contract meaning.
- **MINOR** (1.0.0→1.1.0): new scope/feature added that's backward-compatible — existing consumers aren't broken.
- **MAJOR** (1.0.0→2.0.0): breaking change to a contract another module references via `depends_on` (field removed/renamed, business rule reversed, etc.).

## 2. Mandatory PRD Structure (The 6 Pillars)
An Enterprise-Grade PRD must act as a single source of truth for Business, QA, and Engineering. Every PRD MUST contain the following 6 core sections:

### 2.1. Tujuan & Dampak (The "Why")
Explain *why* this module is being built. What business problem does it solve?
*Example: "Mempercepat proses input data dari 10 menit menjadi 2 menit."*

### 2.2. Scope & Out-of-Scope (Batasan Tegas)
Clearly define what is being built and, critically, what is **NOT** being built right now to prevent feature creep.

### 2.3. User Roles & Permissions
Define who will use this feature and their access levels (e.g., Superadmin, HR Manager, Regular Employee). Detail what each role can Read, Write, or Approve.

### 2.4. Kriteria Penerimaan (Acceptance Criteria)
The strict definition of "Done" to prevent debates between QA and Engineering. Use the **Given-When-Then** format to ensure boundaries are black-and-white and easily convertible into unit tests.

**Implementation gap flagging (MANDATORY):** If, while grounding the PRD in real code (see §1 "Ground truth before writing"), you find the code does **not** actually do what a scenario requires — or does not fulfill a contract another module's PRD already promised depends on it — write the scenario as the *intended/required* behavior anyway, then add a `*Catatan implementasi:*` line directly under it stating what the code currently does instead and that it's a gap to close. Never silently omit the scenario, and never mark it done just because "that's what the code happens to do."

### 2.5. Technical & Architectural Constraints
Define the engineering rules for this module.
- **Backend:** DDD isolation rules, strict typing, or soft-delete mandates.
- **Frontend:** Form structure (e.g., Wizard/Multi-step), UI constraints, data masking, or client-side validations.

### 2.6. Dependencies (Ketergantungan)
Make coupling explicit and versioned. Two directions are MANDATORY:
- **Depends on** — modules this one consumes, with version. State *which field/output* is consumed and **reference** the source PRD section instead of restating its rules. *Example: "Payroll consumes `total_work_hours` from Attendance PRD §4.2 (v1.2)."*
- **Consumed by** — modules that depend on this one (reverse edge). Keeps the impact radius visible when this PRD changes.
- **External integrations** — 3rd-party APIs, SSO, Payment Gateways.

> Principle: loose coupling in docs, same as in code. Never duplicate a parent module's business rule; link to it. If a rule is truly shared, promote it to `docs/PRD/_shared/glossary.md`.

## 3. Data Schema & Business Rules
After the 6 core pillars, provide the logical breakdown of the data models (Entities):
- **Header:** `## [Entity Name]`
- **Aturan Bisnis:** Unique constraints, relations.
- **Sample Data:** A Markdown table illustrating what the data looks like. Columns should reflect the actual fields (e.g., `id`, `name`, `status`, `created_at`). This helps Frontend developers mock the UI.

## 4. Updating an Existing PRD
When the user references an existing PRD and asks for a change/new feature, follow this sequence — don't skip straight to editing:
1. **Brainstorm first.** Align on the scope of the change with the user before rewriting anything.
2. **Re-ground in real code** if the module is implemented (see §1) — the code may have moved since the PRD was last written.
3. Edit only the relevant sections. Don't delete still-valid history unless it's genuinely deprecated.
4. **Bump `version`** per the SemVer criteria in §1a, and refresh `updated`.
5. **Sync `docs/PRD/README.md`** — registry row (version/status/updated) and dependency graph if an edge changed.
6. **Ripple-check dependents (MANDATORY):** search other PRDs' `depends_on` for this module. If the bump is MAJOR (breaking), update the version pointer in every dependent PRD and re-verify the section they reference still says what they think it says.
7. If the module has a `tech-spec.md`/`decision-log.md` (Sedang/Kompleks tier), sync any technical-decision change there too — PRD stays WHAT/WHY, tech-spec stays HOW, never blend them.
8. **Commit atomically** with a `docs:` prefix, separate from any `feat:`/`fix:` code commit ([commit-convention.md](../../rules/commit-convention.md)).
