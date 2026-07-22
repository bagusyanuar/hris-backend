# HRIS Backend - Domain-Driven Design (DDD) & Coding Guidelines

Dokumen ini adalah index aturan project HRIS Backend. Semua agent dan developer harus mematuhi aturan di bawah ini secara ketat. Detail tiap topik dipecah di folder [`rules/`](rules/).

---

## Project Overview & Goals

HRIS Backend adalah sistem Human Resource Information System yang dirancang untuk kebutuhan **enterprise** — perusahaan skala menengah-besar dengan struktur organisasi kompleks (multi-departemen/multi-cabang) dan jumlah karyawan yang terus bertumbuh. Arsitektur DDD yang dipakai (lihat [architecture.md](rules/architecture.md)) sengaja dipilih agar sistem **scalable**: tiap domain adalah bounded context independen yang bisa berkembang sendiri-sendiri, bahkan diekstrak jadi service terpisah di masa depan tanpa merombak keseluruhan sistem.

**Modul yang sudah berjalan** (tiap context = satu folder di `internal/`, pola domain-first):
- **Auth** — autentikasi & access control.
- **User** — manajemen akun pengguna sistem.
- **Organization** — struktur perusahaan/departemen/cabang.
- **Employee** — data & profil karyawan.

**Modul basic HRIS yang jadi arah pengembangan** (belum ada PRD detail, cakupan final menyusul dari tim):
- **Attendance & Time Tracking** — presensi, jam kerja, shift.
- **Leave / Time-off Management** — pengajuan & approval cuti/izin.
- **Payroll & Compensation** — perhitungan gaji, benefit, potongan.
- **Performance Management** — penilaian kinerja karyawan.
- **Recruitment & Onboarding** — proses rekrutmen hingga onboarding karyawan baru.

> Daftar ini adalah gambaran awal scope, bukan PRD resmi. Setiap modul baru tetap **WAJIB** melalui proses PRD/Tech Spec lengkap (lihat [project-docs.md](rules/project-docs.md)) sebelum diimplementasikan.

---

## Daftar Rules

1. [Architecture & Struktur Folder (DDD)](rules/architecture.md) — struktur folder, dependency rules, aturan coding per layer (Domain/Application/Infrastructure/Interfaces).
2. [Konvensi Kode Go](rules/coding-convention.md) — context, error handling, Wire DI, cross-domain communication, config, acronym naming, mandatory build check.
3. [Dokumentasi API (Swagger & Bruno)](rules/api-documentation.md) — aturan wajib dokumentasi endpoint, anti-duplikasi, versioning.
4. [Git Commit & Versioning](rules/commit-convention.md) — Conventional Commits, aturan atomik, changelog.
5. [Dokumen Proyek (PRD, Tech Spec, DBML)](rules/project-docs.md) — lokasi dan format dokumen requirement & teknis.
6. [UUID Generation (Primary Key)](rules/uuid-generation.md) — pola auto-generate UUID di Domain & Infrastructure layer (single source of generation).
7. [Konvensi Persistensi (Repository, Transaction, Data Integrity)](rules/persistence-convention.md) — larangan `db.Save()` untuk upsert, transaction ownership di application layer, kontrak not-found sentinel error, larangan bocor error internal ke client.

---

## Referensi Lain
- Skills: [`skills/`](skills/) — auto-commit, api-validation, go-best-practices, scaffold-prd, scaffold-rfc, scaffold-domain, scaffold-api-docs, scaffold-dbml.
- Workflows / Slash Commands: [`workflows/`](workflows/) — execute-domain, git-commit, scaffold-docs.
