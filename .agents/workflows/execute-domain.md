---
description: Trigger this workflow to automatically generate Go code (DDD layers) for a specific domain based on approved technical specs and DBML.
---

# Execute Domain Code Generation Workflow

This workflow automates the translation of approved Technical Specifications and DBML into actual Go code following the strict DDD rules in this repository.

## Step 1: Initialization & Context Gathering
1. If the user invokes this workflow without specifying a domain, ask them: "Domain/Modul apa yang mau kita koding eksekusinya hari ini?"
2. STOP and wait for the user's response before proceeding.
3. Read whatever context documents exist for the domain (technical docs are **tiered** — a Simpel module may have no `docs/technical/` folder at all; see `rules/project-docs.md`). Gather the fullest available context before writing any code, in this order of preference:
   - `docs/PRD/<domain>.md` — **always read** (source of business rules & acceptance criteria).
   - `docs/technical/<domain>/tech-spec.md` — if it exists (Sedang/Kompleks tier).
   - `docs/technical/<domain>/user-stories.md` — if it exists (Kompleks tier).
   - `docs/technical/<domain>/decision-log.md` — if it exists (constraints/ADR).
   - `docs/databases/<domain>.dbml` — **always read** (mandatory for every tier; source of the physical schema).
   > Below, "the spec" means the strongest source available: tech-spec if present, otherwise the PRD + DBML. For a Simpel module, generate directly from PRD + DBML.
4. Read and deeply understand the coding guidelines in the `scaffold-domain` skill.

## Step 2: Layer-by-Layer Code Generation
You MUST generate the code systematically, from the innermost layer (core) to the outermost layer, ensuring dependencies are correct.

> **CRITICAL RULE (LAYER CONSISTENCY)**: Do not output partial implementations to save tokens. You MUST write the FULL logic for ALL operations defined in the Technical Specification. If the spec dictates 3 endpoints, you must implement exactly those 3 endpoints consistently across the Repository, Application Service, HTTP Handlers, and Routers. Do not skip or use placeholders.

This project uses a **domain-first** layout: one bounded context = one folder `internal/<domain>/` containing all four layers. Follow the `scaffold-domain` skill templates exactly.

1. **Domain Layer:** Generate `entity.go` and `repository.go` (interfaces) inside `internal/<domain>/domain/` (package `domain`). UUID generated inside the constructor only; not-found is a sentinel error, never `(nil, nil)`.
2. **Adapter Layer:** Generate GORM Models inside `internal/<domain>/adapter/models/` (package `models`) and implement the Postgres Repository in `internal/<domain>/adapter/postgres.go` (package `adapter`). Insert with `Create()`, never `Save()` (persistence-convention.md §1).
3. **Application Layer:** Generate Request/Response DTOs (`dto.go`) and Application Service (`service.go`) inside `internal/<domain>/application/` (package `application`). Handle transactions here if required by the Tech Spec; do NOT generate UUIDs here.
4. **Transport Layer:** Generate HTTP handlers (`handler.go`) and routing (`router.go`) inside `internal/<domain>/transport/http/` (package `http`).
5. **Bootstrap/Wiring:** You MUST update `internal/di/api.go` (to register the handler to `APIHandlers` struct) and `internal/di/wire.go` (to add the Repository, Service, and Handler to their respective `ProviderSets`, using descriptive import aliases like `<domain>App`, `<domain>Infra`, `<domain>HTTP`). Then run `go run github.com/google/wire/cmd/wire@latest ./internal/di` to regenerate the DI code.
6. **API Documentation:** Generate the Swagger YAML and Bruno Collection for the new module as per the `scaffold-api-docs` skill.

## Step 3: Verification & Walkthrough
1. **MANDATORY**: Run `go build ./...` to verify there are no syntax or type errors in the newly generated code. If it fails, fix the errors first.
2. After generating and verifying all the files, present a summary of the new files to the user.
3. Ensure you remind the user to run database migrations if there are new DBML changes.
4. Ask the user if they want to run a quick test on the new endpoints.
