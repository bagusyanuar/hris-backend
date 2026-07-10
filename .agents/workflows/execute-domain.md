---
description: Trigger this workflow to automatically generate Go code (DDD layers) for a specific domain based on approved technical specs and DBML.
---

# Execute Domain Code Generation Workflow

This workflow automates the translation of approved Technical Specifications and DBML into actual Go code following the strict DDD rules in this repository.

## Step 1: Initialization & Context Gathering
1. If the user invokes this workflow without specifying a domain, ask them: "Domain/Modul apa yang mau kita koding eksekusinya hari ini?"
2. STOP and wait for the user's response before proceeding.
3. Read the following documents to gather full context before writing any code:
   - `docs/technical/<domain>/tech-spec.md`
   - `docs/technical/<domain>/user-stories.md`
   - `docs/databases/<domain>.dbml`
4. Read and deeply understand the coding guidelines in the `scaffold-domain` skill.

## Step 2: Layer-by-Layer Code Generation
You MUST generate the code systematically, from the innermost layer (core) to the outermost layer, ensuring dependencies are correct.

1. **Domain Layer:** Generate `entity.go` and `repository.go` (interfaces) inside `internal/domain/<domain>/`.
2. **Infrastructure Layer:** Generate GORM Models inside `internal/infrastructure/repository/models/` and implement the Postgres Repository in `internal/infrastructure/repository/`.
3. **Application Layer:** Generate Request/Response DTOs (`dto.go`) and Application Service (`service.go`) inside `internal/application/<domain>/`. Ensure transactions are handled here if required by the Tech Spec.
4. **Interfaces Layer:** Generate HTTP handlers (`handler.go`) and routing (`router.go`) inside `internal/interfaces/http/<domain>/`.
5. **Bootstrap/Wiring:** You MUST update `internal/di/api.go` (to register the handler to `APIHandlers` struct) and `internal/di/wire.go` (to add the Repository, Service, and Handler to their respective `ProviderSets`). Then run `go run github.com/google/wire/cmd/wire@latest ./internal/di` to regenerate the DI code.

## Step 3: Verification & Walkthrough
1. After generating all the files, present a summary of the new files to the user.
2. Ensure you remind the user to run database migrations if there are new DBML changes.
3. Ask the user if they want to run a quick test on the new endpoints.
