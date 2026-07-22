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

## 1a. Mandatory Frontmatter Header
Every PRD file MUST begin with a YAML frontmatter block for versioning and traceability:

```yaml
---
module: Payroll
version: 1.0.0          # SemVer — bump on every business-rule change
status: Draft           # Draft | In Review | Approved | Deprecated
owner: <name>
updated: 2026-07-22     # ISO date, absolute
depends_on: [attendance@1.2, leave@1.0]   # empty [] if none
---
```

Bump `version` (SemVer) and refresh `updated` whenever the PRD content changes, and reflect it in the README index.

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
