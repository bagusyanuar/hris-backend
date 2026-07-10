# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Employee Module:** Implemented complete employee domain module including domain entities, Postgres repository, Application Services, and REST API via Fiber Handlers.
- **Documentation:** Generated PRD, Technical Specs, and DBML for Employee module.
- **AI Automation:** Added `/git-commit` workflow, Conventional Commits rules, and integrated Cross-Domain Bounded Context rules to `scaffold-rfc`.

### Changed
- **Core Architecture:** Migrated Dependency Injection from manual bootstrap to compile-time injection using `google/wire`.
- Consolidated HTTP routes inside `APIHandlers` struct to keep `cmd/api/server.go` clean.
