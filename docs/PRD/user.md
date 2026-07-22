---
module: User
version: 1.0.0
status: Draft
owner: bagusyanuar
updated: 2026-07-22 22:48:47
depends_on: []
---

# Product Requirements: User Module

## 1. Tujuan & Dampak (The "Why")
Menjadi sumber tunggal (Single Source of Truth) identitas akun sistem — pasangan `email` + `password` + `status` — yang dipakai modul Auth untuk otentikasi dan modul lain untuk gate akses. Tanpa modul ini terpisah dari Employee, data karyawan (person) dan data akun login (credential) akan tercampur, padahal keduanya punya siklus hidup berbeda (mis. akun bisa di-suspend tanpa menghapus data kepegawaian, atau ada akun sistem yang bukan karyawan seperti `admin@hris.local`).

## 2. Scope & Out-of-Scope (Batasan Tegas)

**In-Scope (Dikerjakan):**
- Entity akun: `email` (unique), `password` (hash), `status` (`active` / `inactive` / `suspended`).
- Provisioning akun **otomatis** saat proses onboarding Employee disetujui (system-to-system, bukan form pendaftaran terbuka) — memenuhi kontrak yang dijanjikan [employee.md](employee.md) §6 ("Karyawan wajib memiliki `user_id`").
- Perubahan status akun (`active` ↔ `suspended` ↔ `inactive`), dipicu oleh event dari modul lain (mis. offboarding Employee → `inactive`) atau tindakan Admin langsung.
- Penyediaan kontrak baca (`FindByEmail`, `FindByID`) untuk dikonsumsi modul Auth.
- Penyimpanan password **wajib** dalam bentuk hash (`bcrypt`), tidak pernah plaintext, di titik manapun proses provisioning/ubah password terjadi.

**Out-of-Scope (TIDAK Dikerjakan di modul ini, untuk saat ini):**
- **Self sign-up / registrasi publik** — akun hanya dibuat via provisioning sistem (Employee onboarding) atau Admin, bukan form publik.
- **RBAC / permission matrix** — `role` adalah tanggung jawab bersama Auth+modul Access Control masa depan, bukan atribut inti User.
- **Reset password via email (forgot password flow)** — belum ada integrasi email/notifikasi.
- **Ganti password self-service oleh end-user** — belum ada endpoint; saat ini password hanya bisa di-set ulang lewat skrip seed (`cmd/seed/main.go`) atau Admin manual di DB.
- **Audit log percobaan login** — kalau dibutuhkan, jadi domain Auth atau modul Security terpisah, bukan bagian inti User.
- **Multiple email / phone-based login**, **akun non-email (username-only)**.

## 3. User Roles & Permissions
- **Superadmin / HR Admin**: Berhak membuat akun sistem (di luar jalur Employee onboarding — mis. akun admin), mengubah `status` akun manapun (suspend/reaktivasi), tapi **tidak** boleh melihat password (hash tidak pernah diekspos ke response API manapun).
- **Sistem (internal, dipicu modul lain)**: Employee module memicu pembuatan akun otomatis saat onboarding disetujui; Employee module memicu perubahan status ke `inactive` saat offboarding. Ini bukan aksi manusia, tapi *service-to-service call* (Application Service Employee memanggil Application Service User, sesuai [coding-convention.md](../../.agents/rules/coding-convention.md) §4 — dilarang injeksi Repository lintas modul langsung).
- **Pemilik akun (User biasa)**: Untuk saat ini **tidak** ada hak akses langsung ke modul ini (tidak bisa ubah profil/password sendiri) — hanya jadi subjek yang datanya dikonsumsi Auth saat login.

## 4. Kriteria Penerimaan (Acceptance Criteria)

**Skenario 1: Provisioning Akun Otomatis saat Onboarding**
- **Given** HR Admin menyetujui data karyawan baru di modul Employee.
- **When** proses onboarding selesai.
- **Then** sistem membuat 1 baris `users` baru dengan `status = active`, hash password default/temporary, dan mengaitkan `id`-nya sebagai `user_id` ke record Employee bersangkutan.
- *Catatan implementasi:* alur ini **belum ada** di kode saat ini — `internal/application/employee` belum pernah memanggil apapun dari modul User meskipun `Employee.UserID` sudah jadi field di entity. Ini gap yang harus ditutup agar kontrak [employee.md](employee.md) §6 benar-benar terpenuhi, bukan cuma asumsi di dokumen.

**Skenario 2: Duplikasi Email Ditolak**
- **Given** email yang akan dipakai provisioning sudah terpakai user lain (`deleted_at IS NULL`).
- **When** proses pembuatan akun dijalankan.
- **Then** sistem menolak dengan sentinel error (mis. `ErrEmailDuplicate`), tidak membuat baris baru, dan mengembalikan status HTTP `409 Conflict` di layer yang memanggilnya.

**Skenario 3: Offboarding Menonaktifkan Akun**
- **Given** karyawan diproses resign/PHK di modul Employee.
- **When** offboarding selesai.
- **Then** `users.status` milik akun terkait berubah jadi `inactive`, dan sejak saat itu login lewat modul Auth ditolak (lihat [auth.md](auth.md) §4 Skenario 3).

**Skenario 4: Password Selalu Ter-hash**
- **Given** proses apapun yang menyimpan/mengubah password (provisioning, reset manual).
- **When** nilai password disimpan ke `users.password`.
- **Then** nilai yang tersimpan **wajib** hasil `bcrypt.GenerateFromPassword`, tidak pernah nilai plaintext dari input.

**Skenario 5: Pencarian Akun Tidak Ditemukan**
- **Given** `FindByEmail` atau `FindByID` dipanggil dengan email/id yang tidak ada (atau sudah soft-deleted).
- **When** query dijalankan.
- **Then** repository mengembalikan sentinel error `ErrUserNotFound` (bukan `nil, nil`), sesuai [persistence-convention.md](../../.agents/rules/persistence-convention.md) §3 — kontrak ini **sudah** terimplementasi dengan benar di kode saat ini.

## 5. Technical & Architectural Constraints
- **Domain-Driven Design (DDD)**: Domain layer (`domain/user/entity.go`) sudah pure Go, tanpa GORM tag — pertahankan. Model DB terpisah di `infrastructure/repository/models/user_model.go` dengan mapper `ToDomain()`/`FromDomain()`.
- **UUID Generation Gap**: Sesuai [uuid-generation.md](../../.agents/rules/uuid-generation.md), constructor `NewUser` seharusnya auto-generate UUID kalau `id` kosong. Saat ini `NewUser` justru mengembalikan `ErrInvalidInput` kalau `id == ""` — **tidak konsisten** dengan pola domain lain (Employee, Organization). Wajib diperbaiki saat implementasi fitur Create/provisioning.
- **Persistence**: Repository interface saat ini hanya `FindByEmail` & `FindByID` (read-only). Method `Create`/`UpdateStatus` yang akan ditambah **wajib** ikut [persistence-convention.md](../../.agents/rules/persistence-convention.md) — gunakan `Create()` utuh (bukan `Save()`) untuk insert baru, dan constraint unique violation pada `email` wajib diterjemahkan ke sentinel error, bukan bocor error driver Postgres mentah ke client.
- **Cross-Domain Communication**: Pemicu dari Employee ke User (provisioning, status change) **wajib** lewat Application Service User (`di/wire.go` injection), bukan akses Repository User langsung dari Employee — sesuai [coding-convention.md](../../.agents/rules/coding-convention.md) §4.
- **Data Deletion**: Soft delete (`deleted_at`), konsisten dengan pola Employee — jangan hard delete akun (jejak audit login harus tetap bisa ditelusuri).

## 6. Dependencies (Ketergantungan)
- **Depends on**: Tidak ada — User adalah root/foundational context, tidak butuh data modul lain untuk berfungsi.
- **Consumed by — Modul Auth**: Auth memanggil `FindByEmail`/`FindByID` dan membaca `status` sebagai gate login. Lihat [auth.md](auth.md) §6.
- **Consumed by — Modul Employee**: Employee butuh User untuk provisioning akun saat onboarding dan menonaktifkan akun saat offboarding (lihat [employee.md](employee.md) §6). Arah panggilan: Employee → User Application Service (bukan sebaliknya).
- **External integrations**: Tidak ada saat ini. Kalau nanti ada *forgot password* flow, akan butuh integrasi email/notification service (di luar scope dokumen ini).

---

## 7. Data Schema & Business Rules

Skema fisik sudah ada di [docs/databases/user.dbml](../databases/user.dbml) (tidak perlu file DBML baru — tinggal dipastikan konsisten dengan Acceptance Criteria di atas).

### 7.1. User (Akun Sistem)
**Aturan Bisnis:**
- `email` unique, case-sensitivity ikut default Postgres (`citext` bisa dipertimbangkan kalau perlu case-insensitive — belum diimplementasikan).
- `status` hanya boleh salah satu dari: `active`, `inactive`, `suspended` — validasi enum sebaiknya di level domain constructor, bukan cuma dokumentasi.
- `password` tidak pernah diekspos di response API manapun (DTO Response wajib exclude field ini).
- Satu `email` = satu akun, satu akun bisa dipakai lintas konteks (Employee ATAU akun non-employee seperti admin).

**Sample Data:**

| id | email | status | created_at |
|---|---|---|---|
| `a1b2...` | `admin@hris.local` | `active` | `2026-01-10T08:00:00Z` |
| `c3d4...` | `employee@hris.local` | `active` | `2026-01-10T08:00:00Z` |
| `e5f6...` | `resigned.user@hris.local` | `inactive` | `2025-11-02T09:15:00Z` |
