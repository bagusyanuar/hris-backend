---
name: auto-commit
description: Guide for agents to handle Git commits automatically and atomically using Conventional Commits.
---

# Auto-Commit Workflow

Trigger this skill whenever the user asks to "commit", "save changes", or explicitly requests the agent to handle git operations.

## 1. Safety & Staging Area Check
If the user asks you to commit, **DO NOT blindly run `git commit` or `git add .`**.
1. Run `git status` to see the current state.
2. If files are already staged (`Changes to be committed`), but they belong to multiple different domains/modules, run `git reset` (or `git restore --staged .`) to unstage everything. We must perform Atomic Commits.

## 2. Linting & Validation
Before grouping and committing, run `make lint` to verify code quality. If the linter reports errors, you MUST fix them before proceeding.

## 3. Analyze & Group Files
1. Run `git status` and `git diff` to understand what was changed.
2. Group the modified/untracked files by their bounded context or domain. For example:
   - `internal/domain/auth/...` + `docs/api/swagger/auth.yaml` -> **Auth Module**
   - `internal/domain/organization/...` -> **Organization Module**
   - `pkg/validator/...` + `cmd/api/...` -> **Core/Infra Module**

## 4. Formulate Atomic Commits
For each module group, formulate a commit message adhering to Conventional Commits format:
`type(scope): subject`

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

**Scope:** The module name (e.g., `auth`, `organization`, `core`).

**Description (Changelog):**
Include a bulleted list describing the changes.
Example:
```text
feat(auth): implement request validation

- Add validate tags to LoginRequest DTO
- Implement ValidateStruct in Login handler
- Document 422 responses in Swagger and Bruno
```

## 5. Execute Commits & Changelog
1. **Minta Persetujuan (Ask User):** Sebelum benar-benar menjalankan `git commit`, tampilkan **Rencana Commit** (daftar grup, file, dan pesan commit) kepada user. Tanyakan apakah mereka setuju.
2. **Update Changelog:** Jika user setuju dan perubahannya signifikan (feat/fix/breaking change), tulis perubahan tersebut ke dalam file `CHANGELOG.md` di *root directory* (di bawah *header* `## [Unreleased]`) jika diminta.
3. **Eksekusi:** Setelah user bilang "Ya", eksekusi commit secara berurutan.
```bash
# Group 1: Auth
git add internal/interfaces/http/auth/handler.go docs/api/swagger/auth.yaml docs/api/bruno/Auth/Login.bru
git commit -m "feat(auth): implement request validation" -m "- Add validate tags to LoginRequest DTO" -m "- Implement ValidateStruct in Login handler" -m "- Document 422 responses in Swagger and Bruno"

# Group 2: Organization
git add internal/interfaces/http/organization/handler.go
git commit -m "feat(organization): add organization endpoints" -m "- Implement CRUD operations for organization"
```

> **IMPORTANT**: Never use `git add .` if changes span multiple domains. Always be specific with file paths to ensure atomic commits.
