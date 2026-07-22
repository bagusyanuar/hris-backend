---
description: Trigger this workflow to generate complete enterprise documentation (PRD and Technical Specs) for a new HRIS domain sequentially.
---

# Scaffold Domain Documentation Workflow

This workflow ensures that every time a new domain/module is requested, the AI follows a strict, sequential pipeline to generate Enterprise-Grade documentation based on the project's rules.

## Step 1: Initialization & Gathering Context
1. If the user invokes this workflow without specifying a domain, ask them: "Domain/Fitur apa yang mau kita buat dokumentasinya hari ini? Boleh kasih gambaran singkat alur bisnisnya?"
2. STOP and wait for the user's response before proceeding.

## Step 2: Generate Product Requirements Document (PRD)
1. Read and apply the instructions from the `scaffold-prd` skill.
2. Create the PRD document and save it exactly at `docs/PRD/<domain_name>.md`.
3. Inform the user that the PRD has been generated.
4. STOP and ask the user to review the PRD. You MUST NOT proceed to generate technical documents until the user explicitly approves the PRD.

## Step 3: Classify Complexity (Tiering Gate)
Technical docs are **tiered** — not every module needs the full set (see `rules/project-docs.md`). After PRD approval, classify the module and confirm with the user before generating.

1. Propose a tier based on the PRD, using these criteria:
   - **Simpel** — CRUD lurus, 1–2 entity, tanpa integrasi luar, tanpa kalkulasi/state machine (mis. master data, lookup).
   - **Sedang** — ada relasi antar-entity, state/status flow, atau depend antar-modul.
   - **Kompleks** — kalkulasi berlapis, state machine, atau integrasi eksternal (mis. Payroll, Attendance, Leave).
2. State your proposed tier + reason, then ask the user to confirm or override. STOP and wait.

## Step 4: Generate Technical Specifications (RFC) — by tier
Once tier is confirmed, read and apply the `scaffold-rfc` skill and generate **only** what the tier requires:

- **Simpel:** Skip `docs/technical/`. No tech-spec. (Kode boleh langsung di-scaffold dari PRD + DBML.)
- **Sedang:** Create `docs/technical/<domain_name>/` with `tech-spec.md` (Architecture, API contract, Mermaid ERD). Add `decision-log.md` if any non-trivial technical decision needs its *why* recorded.
- **Kompleks:** Create `docs/technical/<domain_name>/` with the full set:
  - `tech-spec.md` (Architecture, API, Schema, Mermaid ERD)
  - `user-stories.md` (Scenarios, Edge Cases, Mermaid Sequence Diagrams)
  - `decision-log.md` (ADR for architectural decisions made during tech spec generation)
  - Add `data-dictionary.md`, `infrastructure.md`, or `test-plan.md` only if the module warrants them.

## Step 5: Generate DBML (ALL tiers — mandatory)
1. Generate the database schema file in DBML format at `docs/databases/<domain_name>.dbml`. This is **required for every tier**, including Simpel — DBML is the single source for SQL migrations (`AutoMigrate` is banned in prod).
2. Present the generated documents (whatever the tier produced) + DBML to the user for final engineering review.
