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

## Step 3: Generate Technical Specifications (RFC) & DBML
1. Once the user approves the PRD, read and apply the instructions from the `scaffold-rfc` skill.
2. Create a new directory for the domain: `docs/technical/<domain_name>/`.
3. Generate the following files inside that directory based on the PRD:
   - `tech-spec.md` (Architecture, API, Schema, Mermaid ERD)
   - `user-stories.md` (Scenarios, Edge Cases, Mermaid Sequence Diagrams)
   - `decision-log.md` (ADR for architectural decisions made during tech spec generation)
   - (Generate `data-dictionary.md`, `infrastructure.md`, or `test-plan.md` only if the domain is complex enough to warrant them).
4. Generate the database schema file in DBML format at `docs/databases/<domain_name>.dbml`.
5. Present the generated technical documents and DBML to the user for final engineering review.
