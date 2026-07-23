# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Observability:** Introduced `pkg/logger` (structured logging via `zap`), request-scoped logging middleware with `request_id` correlation, and a `logging-convention.md` rule documenting the pattern (log once at boundary, structured fields, no raw error leaks).
- **Database Resiliency:** `internal/shared/database` now configures connection pool (`SetMaxOpenConns`/`SetMaxIdleConns`/lifetime), retries connection with backoff + timeout on startup, and routes GORM query logs through `zap` (slow-query threshold, env-aware log level).
- **Health Check:** Added `GET /health` endpoint pinging the database, for liveness/readiness probes.
- **Employee Module CRUD:** Completed missing endpoints (`FindAll`, `Update`, `Delete`) in Application Service and Fiber Handlers, along with Swagger YAML and Bruno collections for each.
- **Employee Module:** Implemented complete employee domain module including domain entities, Postgres repository, Application Services, and REST API via Fiber Handlers.
- **Documentation:** Generated PRD, Technical Specs, and DBML for Employee module.
- **AI Automation:** Added `/git-commit` workflow, Conventional Commits rules, and integrated Cross-Domain Bounded Context rules to `scaffold-rfc`.

### Changed
- **AI Automation:** Refined `scaffold-domain` skill and `execute-domain` workflow to mandate **Layer Consistency** (implementing only exactly what's in the Tech Spec fully without placeholders) rather than forcing 5 CRUD endpoints blindly.
- **Core Architecture:** Migrated Dependency Injection from manual bootstrap to compile-time injection using `google/wire`.
- Consolidated HTTP routes inside `APIHandlers` struct to keep `cmd/api/server.go` clean.

### Fixed
- **Organization & Employee:** 500 responses no longer leak raw `err.Error()` to the client; the real error is now logged server-side instead.
- **Auth:** Errors swallowed during login/token-refresh (DB lookup failures reclassified into sentinel errors) are now logged before being discarded.
- **Auth Middleware:** `AuthProtected` now returns the standard `pkg/response` envelope instead of an ad-hoc `fiber.Map`, matching every other endpoint's error shape.
