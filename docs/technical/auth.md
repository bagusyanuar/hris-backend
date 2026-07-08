# Technical Documentation: Authentication (Auth)

## 1. Overview
Modul Authentication bertanggung jawab untuk proses verifikasi identitas (Login) dan penerbitan JSON Web Token (JWT) yang akan digunakan untuk otorisasi endpoint-endpoint lain di HRIS Backend.

Sesuai aturan arsitektur project, modul ini akan mengikuti prinsip Domain-Driven Design (DDD).
Sistem otentikasi akan menggunakan pendekatan **Hybrid**, di mana `access_token` akan dikirim via JSON response dan `refresh_token` akan diset secara otomatis ke dalam **HttpOnly Secure Cookie** untuk perlindungan maksimal terhadap eksploitasi serangan XSS.

## 2. Arsitektur & Layering (DDD)

### A. Domain Layer (`internal/domain/auth`)
*   **Interface**: `TokenGenerator` (Interface untuk men-generate dan memvalidasi token JWT). Tidak melakukan koneksi ke DB langsung, murni kontrak bisnis terkait token.
*   *Catatan*: Untuk verifikasi user, `auth.Service` di Application Layer akan berinteraksi dengan `user.Repository` (melalui domain `user`) untuk mengecek eksistensi email dan validasi password.

### B. Application Layer (`internal/application/auth`)
*   **Service (`service.go`)**: 
    *   Menerima DTO `LoginRequest` (email, password).
    *   Memanggil `user.Repository.FindByEmail(email)`.
    *   Memvalidasi password yang dikirim dengan yang ada di database menggunakan `bcrypt.CompareHashAndPassword`.
    *   Jika valid, memanggil `TokenGenerator.GenerateToken(userID)`. (Harus menghasilkan pasangan Access Token dan Refresh Token).
    *   Mengembalikan token pair ke Handler.

### C. Infrastructure Layer (`internal/infrastructure/security`)
*   **JWT Implementation (`jwt.go`)**:
    *   Tempat diletakkannya implementasi dari interface `TokenGenerator` yang ada di domain auth.
    *   Menggunakan package eksternal `github.com/golang-jwt/jwt/v5`.
    *   Membaca dan menggunakan secret key dari `config.Config`.
    *   Menghasilkan 2 buah token terpisah: *Access Token* (umur pendek) dan *Refresh Token* (umur panjang).

### D. Interfaces / Presentation Layer (`internal/interfaces/http`)
*   **Handler (`auth_handler.go`)**:
    *   Endpoint: `POST /api/v1/auth/login`
        *   Menerima request JSON, melakukan binding & validasi input, memanggil `auth.Service.Login`.
        *   Menyisipkan Refresh Token ke HTTP Cookie (`HttpOnly`, `Secure`, `SameSite=Strict`).
        *   Mengirimkan response JSON berisi Access Token.
    *   Endpoint: `POST /api/v1/auth/refresh`
        *   Mengekstrak Refresh Token dari Cookie.
        *   Memvalidasi Refresh Token dan menerbitkan Access Token baru.
*   **Middleware (`middleware/auth_middleware.go`)**:
    *   Memeriksa header `Authorization: Bearer <token>` untuk Access Token.
    *   Memvalidasi JWT token melalui `TokenGenerator.ValidateToken`.
    *   Menyisipkan `userID` ke dalam context Fiber (`c.Locals("userID", claims.UserID)`).

## 3. Alur Proses (Flow)

### Flow Login
1.  **Client** mengirim `POST /api/v1/auth/login` dengan body `email` dan `password`.
2.  **AuthHandler** menangkap request, melakukan validasi dasar, dan memanggil `AuthService.Login(ctx, req)`.
3.  **AuthService**:
    *   Meminta data ke `UserRepository.FindByEmail(email)`.
    *   Jika tidak ketemu -> kembalikan error `ErrInvalidCredentials`.
    *   Jika ketemu, bandingkan hash password.
    *   Jika password salah -> kembalikan error `ErrInvalidCredentials`.
    *   Jika benar, panggil `TokenGenerator.GenerateToken(user.ID)` untuk membuat Access dan Refresh Token.
4.  **AuthService** mengembalikan pasangan token ke **AuthHandler**.
5.  **AuthHandler** menge-set header `Set-Cookie` berisi Refresh Token (`HttpOnly`).
6.  **AuthHandler** merespon client dengan HTTP 200 OK dengan payload berisi JSON `access_token`.

### Flow Refresh Token
1.  **Client** (saat sadar `access_token` habis) memanggil endpoint `POST /api/v1/auth/refresh`. Browser secara otomatis melampirkan cookie Refresh Token yang sebelumnya diset.
2.  **AuthHandler** membaca `refresh_token` dari request cookies.
3.  **AuthService** memvalidasi `refresh_token`. Jika valid, generate `access_token` baru.
4.  **AuthHandler** merespon dengan `access_token` baru di JSON Payload.

### Flow Middleware (Proteksi Endpoint)
1.  **Client** mengirim request ke endpoint terlindungi (misal `GET /api/v1/users/me`) dengan header `Authorization: Bearer <token>`.
2.  **AuthMiddleware** mencegat request tersebut.
3.  Middleware mengekstrak token dari header, lalu memvalidasi signature dan masa berlakunya.
4.  Jika invalid / expired -> tolak dengan response HTTP 401 Unauthorized.
5.  Jika valid -> ekstrak `userID` dari JWT claims, simpan di Context lokal request, dan lanjutkan eksekusi ke handler berikutnya (`c.Next()`).

## 4. Struktur Data DTO (Data Transfer Object)

```go
package auth

// Digunakan di Application Layer dan di-binding oleh Handler
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Hanya mengirimkan access token, karena refresh token dikirim via cookie.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // waktu kadaluarsa (dalam detik)
	TokenType   string `json:"token_type"` // default: "Bearer"
}

// Domain model internal untuk transfer balikan dari service (belum tentu semuanya terekspos ke respon API)
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int
}
```

## 5. Next Steps / Action Items
Untuk merealisasikan fitur di atas, berikut adalah urutan pengerjaannya:
1.  **Domain User**: Buat `internal/domain/user/entity.go` dan `internal/domain/user/repository.go` terlebih dahulu karena login butuh ngecek user.
2.  **Infrastructure User**: Buat implementasi `UserRepository` di `internal/infrastructure/repository/user_postgres.go`.
3.  **Domain Auth & Security**: Buat interface generator token di `domain/auth` dan implementasikan JWT-nya di `infrastructure/security/jwt.go` (beserta konfigurasi untuk Access & Refresh Token duration).
4.  **Application Auth**: Buat `AuthService` di `internal/application/auth/service.go`.
5.  **Interfaces Auth**: Buat `AuthHandler` (termasuk urusan parsing/setting Cookie) dan `AuthMiddleware`.
6.  **Routing**: Hubungkan semuanya di `cmd/api/main.go` (Dependency Injection).
