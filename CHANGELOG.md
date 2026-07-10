# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Employee Module CRUD:** Completed missing endpoints (`FindAll`, `Update`, `Delete`) in Application Service and Fiber Handlers, along with Swagger YAML and Bruno collections for each.
- **Employee Module:** Implemented complete employee domain module including domain entities, Postgres repository, Application Services, and REST API via Fiber Handlers.
- **Documentation:** Generated PRD, Technical Specs, and DBML for Employee module.
- **AI Automation:** Added `/git-commit` workflow, Conventional Commits rules, and integrated Cross-Domain Bounded Context rules to `scaffold-rfc`.

### Changed
- **AI Automation:** Refined `scaffold-domain` skill and `execute-domain` workflow to mandate **Layer Consistency** (implementing only exactly what's in the Tech Spec fully without placeholders) rather than forcing 5 CRUD endpoints blindly.
- **Core Architecture:** Migrated Dependency Injection from manual bootstrap to compile-time injection using `google/wire`.
- Consolidated HTTP routes inside `APIHandlers` struct to keep `cmd/api/server.go` clean.
