---
module: Auth
version: 1.0.0
status: Draft
owner: bagusyanuar
updated: 2026-07-22 22:48:47
depends_on: [user@1.0.0]
---

# Product Requirements: Auth Module

## 1. Tujuan & Dampak (The "Why")
Menyediakan satu pintu masuk (Single Point of Entry) otentikasi untuk seluruh sistem HRIS, menggantikan kebutuhan tiap modul membuat mekanisme login sendiri-sendiri. Modul ini menjamin hanya pengguna dengan kredensial valid dan akun aktif yang bisa mengakses API, sekaligus jadi lapisan pertahanan pertama (perimeter security) sebelum request menyentuh domain bisnis lain (Employee, Organization, dst).

## 2. Scope & Out-of-Scope (Batasan Tegas)

**In-Scope (Dikerjakan):**
- Login menggunakan `email` + `password` (dicocokkan dengan hash `bcrypt` di tabel `users`).
- Penerbitan token pasangan (*Token Pair*): `access_token` (umur pendek, dikirim via JSON response) dan `refresh_token` (umur panjang, diset via **HttpOnly Secure Cookie**).
- Endpoint refresh token dengan pola **rotation** (tiap refresh berhasil, refresh token lama diganti baru).
- Middleware `AuthProtected` untuk memvalidasi `access_token` (`Authorization: Bearer <token>`) pada endpoint yang butuh proteksi, dan menyisipkan `userID` + `role` ke context request.
- Validasi status akun (`users.status`) — akun yang tidak `active` (mis. `inactive`, `suspended`) **wajib** ditolak login, meskipun kredensial benar. Ini adalah kontrak yang dijanjikan modul Employee ([employee.md](employee.md) §4 Skenario 2 — offboarding memblokir akses login).

**Out-of-Scope (TIDAK Dikerjakan di modul ini, untuk saat ini):**
- **Registrasi akun baru (Sign Up)** — pembuatan baris `users` adalah tanggung jawab proses onboarding Employee, bukan Auth.
- **RBAC granular / permission matrix** — token membawa klaim `role`, tapi otorisasi berbasis role (mis. "hanya HR Manager boleh akses endpoint X") **belum** diimplementasikan di middleware manapun. Baru sebatas propagasi klaim ke context.
- **Logout server-side / token revocation** — karena token bersifat *stateless JWT* tanpa token store/blacklist, tidak ada mekanisme invalidasi token sebelum masa berlakunya habis. Endpoint logout (jika ada nanti) hanya akan menghapus cookie sisi client, bukan mencabut validitas token.
- **Reset password / lupa password** — belum ada endpoint maupun flow email verifikasi.
- **Account lockout / brute-force protection** (mis. limit percobaan login gagal) — belum ada rate limiting di layer ini.
- **Multi-Factor Authentication (MFA/OTP)**.
- **Social Login / SSO (Google, Microsoft, dll)**.

## 3. User Roles & Permissions
Modul Auth tidak mendefinisikan role bisnis sendiri — ia hanya **membawa** (carrier) klaim `role` yang sumber datanya dari modul lain (saat ini masih *hardcoded* `"employee"` di `application/auth/service.go` sampai RBAC modul User/Organization matang).

- **Semua pengguna terotentikasi**: berhak memperoleh `access_token` + `refresh_token` selama kredensial valid dan `status = active`. Tidak ada perbedaan hak akses antar role di level Auth — pembatasan hak akses granular adalah tanggung jawab modul consumer (Employee, Organization, dst) atau RBAC module di masa depan.
- **Pengguna nonaktif/suspended**: ditolak login sepenuhnya (lihat §2), tanpa pengecualian role.

## 4. Kriteria Penerimaan (Acceptance Criteria)

**Skenario 1: Login Berhasil**
- **Given** pengguna dengan akun `status = active` memasukkan email dan password yang benar.
- **When** request `POST /api/v1/auth/login` dikirim.
- **Then** sistem merespon `200 OK` berisi `access_token`, `expires_in`, `token_type`, dan menyisipkan `refresh_token` sebagai cookie `HttpOnly`, `Secure`, `SameSite=Strict`.

**Skenario 2: Login Gagal — Kredensial Salah**
- **Given** email tidak terdaftar, atau email terdaftar tapi password salah.
- **When** request login dikirim.
- **Then** sistem merespon `401 Unauthorized` dengan pesan generik *"Invalid credentials"* — **tidak** membedakan pesan antara "email tidak ditemukan" vs "password salah" (mencegah *user enumeration*).

**Skenario 3: Login Ditolak — Akun Tidak Aktif**
- **Given** akun dengan email/password benar tapi `status != active` (mis. sudah di-offboard oleh modul Employee).
- **When** request login dikirim.
- **Then** sistem menolak dengan `401 Unauthorized`, tanpa menerbitkan token apapun.
- *Catatan implementasi:* pengecekan status ini **belum ada** di kode saat ini (`application/auth/service.go` hanya cek password, tidak cek `u.Status()`) — dicatat sebagai gap yang harus ditutup, karena kontrak ini sudah dijanjikan ke modul Employee.

**Skenario 4: Refresh Token Rotation**
- **Given** client memiliki `refresh_token` valid dan belum kedaluwarsa (dikirim otomatis via cookie).
- **When** request `POST /api/v1/auth/refresh` dikirim.
- **Then** sistem menerbitkan `access_token` baru dan `refresh_token` baru (rotation), meng-overwrite cookie lama.

**Skenario 5: Refresh Token Invalid/Kedaluwarsa**
- **Given** `refresh_token` tidak ada di cookie, atau ada tapi invalid/expired/salah tipe (mis. access token dipakai di endpoint refresh).
- **When** request refresh dikirim.
- **Then** sistem merespon `401 Unauthorized` dengan pesan *"Invalid or expired refresh token"*.

**Skenario 6: Akses Endpoint Terproteksi Tanpa Token**
- **Given** request ke endpoint yang dibungkus middleware `AuthProtected` tanpa header `Authorization`.
- **When** request dikirim.
- **Then** sistem merespon `401 Unauthorized` — *"missing authorization header"*.

**Skenario 7: Akses Endpoint Terproteksi dengan Access Token Kedaluwarsa/Invalid**
- **Given** header `Authorization: Bearer <token>` terisi tapi token expired, salah signature, atau bertipe `refresh` (bukan `access`).
- **When** request dikirim ke endpoint terproteksi.
- **Then** sistem merespon `401 Unauthorized` — *"invalid or expired token"*, request tidak diteruskan ke handler.

## 5. Technical & Architectural Constraints
- **Domain-Driven Design (DDD)**: Domain Layer Auth (`domain/auth`) hanya berisi kontrak `TokenGenerator` + `TokenPair`/`TokenClaims` — murni abstraksi token, tidak menyentuh DB. Verifikasi user dilakukan Application Layer dengan memanggil `user.Repository` milik domain `user` (bukan bypass query langsung).
- **Hybrid Token Storage**: `access_token` dikirim via JSON body (disimpan FE di memory, bukan localStorage, untuk kurangi risiko XSS-exfiltration); `refresh_token` **wajib** HttpOnly Secure Cookie (`SameSite=Strict`) agar tidak terjangkau JavaScript sisi client.
- **Stateless JWT**: Token ditandatangani `HS256` dengan secret dari config, TIDAK ada session/token store di DB — konsekuensinya tidak ada revocation sebelum expiry (lihat §2 Out-of-Scope).
- **Password Hashing**: Wajib `bcrypt.CompareHashAndPassword`, dilarang membandingkan password plaintext.
- **Error Handling**: Kredensial salah dan akun tidak ditemukan **wajib** mengembalikan pesan generik yang sama (`ErrInvalidCredentials`) untuk mencegah user enumeration — jangan bocorkan mana yang salah.

## 6. Dependencies (Ketergantungan)
- **Depends on — Modul User** ([user.md](user.md) §7.1, v1.0.0): Auth mengonsumsi `user.Repository.FindByEmail` dan `FindByID` untuk ambil `id`, `password` (hash), `status` dari tabel `users`.
- **Consumed by — Modul Employee**: Employee PRD ([employee.md](employee.md) §4 Skenario 2 & §6) menjanjikan bahwa proses offboarding memblokir login lewat modul Auth — kontrak ini dipenuhi lewat pengecekan `users.status` di §2/§4 Skenario 3 dokumen ini.
- **Consumed by — Semua modul terproteksi**: Organization, Employee, dan modul HRIS masa depan (Attendance, Leave, Payroll, dll) bergantung pada middleware `AuthProtected` untuk gerbang otentikasi endpoint mereka.
- **External integrations**: Tidak ada (belum ada SSO/OAuth pihak ketiga).

---

## 7. Data Schema & Business Rules

Auth **tidak memiliki tabel sendiri** — modul ini murni konsumen tabel `users` (dimiliki modul User) dan penerbit JWT stateless (tidak disimpan ke DB).

### 7.1. `users` (Dikonsumsi, Dimiliki Modul User)
Skema fisik ada di [docs/databases/user.dbml](../databases/user.dbml). Field yang relevan bagi Auth:

| Field | Tipe | Dipakai Untuk |
|---|---|---|
| `id` | uuid | Klaim `user_id` di JWT |
| `email` | varchar(255), unique | Kredensial login |
| `password` | varchar(255) | Hash `bcrypt`, dicocokkan saat login |
| `status` | varchar(50) | Gate login — hanya `active` yang boleh masuk |

### 7.2. Token Pair (Bukan Tabel — Struktur JWT Payload)
Tidak disimpan ke database. Struktur klaim di dalam JWT:

| Field | Tipe | Keterangan |
|---|---|---|
| `user_id` | string | ID pengguna dari tabel `users` |
| `role` | string | Saat ini hardcoded `"employee"`, menunggu RBAC |
| `type` | string | `"access"` atau `"refresh"` — mencegah token tipe salah dipakai di endpoint lain |
| `exp` / `iat` / `nbf` | JWT standard claims | Masa berlaku token |

**Sample response Login (`200 OK`):**

| access_token | expires_in | token_type |
|---|---|---|
| `eyJhbGciOi...` | `3600` | `Bearer` |

*(`refresh_token` tidak muncul di body — dikirim via `Set-Cookie` header, `HttpOnly`.)*
