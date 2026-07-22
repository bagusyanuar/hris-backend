# HRIS Backend - Domain-Driven Design (DDD) & Coding Guidelines

Dokumen ini adalah index aturan project HRIS Backend. Semua agent dan developer harus mematuhi aturan di bawah ini secara ketat. Detail tiap topik dipecah di folder [`rules/`](rules/).

---

## Daftar Rules

1. [Architecture & Struktur Folder (DDD)](rules/architecture.md) — struktur folder, dependency rules, aturan coding per layer (Domain/Application/Infrastructure/Interfaces).
2. [Konvensi Kode Go](rules/coding-convention.md) — context, error handling, Wire DI, cross-domain communication, config, acronym naming, mandatory build check.
3. [Dokumentasi API (Swagger & Bruno)](rules/api-documentation.md) — aturan wajib dokumentasi endpoint, anti-duplikasi, versioning.
4. [Git Commit & Versioning](rules/commit-convention.md) — Conventional Commits, aturan atomik, changelog.
5. [Dokumen Proyek (PRD, Tech Spec, DBML)](rules/project-docs.md) — lokasi dan format dokumen requirement & teknis.
6. [UUID Generation (Primary Key)](rules/uuid-generation.md) — pola auto-generate UUID di Domain & Infrastructure layer.

---

## Referensi Lain
- Skills: [`skills/`](skills/) — auto-commit, api-validation, go-best-practices, scaffold-prd, scaffold-rfc, scaffold-domain, scaffold-api-docs.
- Workflows / Slash Commands: [`workflows/`](workflows/) — execute-domain, git-commit, scaffold-docs.
