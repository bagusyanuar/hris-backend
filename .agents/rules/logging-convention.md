# Konvensi Logging (Zap)

Semua structured logging di project ini **WAJIB** lewat `pkg/logger` (wrapper `go.uber.org/zap`). Dilarang pakai `log` stdlib (`log.Println`/`log.Printf`) di kode baru — satu-satunya pengecualian adalah bootstrap sebelum `logger.Init()` sempat dipanggil (mis. `config.LoadConfig()` gagal baca `.env`, lihat [internal/shared/config/config.go](../../internal/shared/config/config.go)).

---

## 1. Inisialisasi — Single Source

`logger.Init(env, debug)` dipanggil **sekali** di `cmd/api/main.go`, setelah `config.LoadConfig()` dan sebelum `database.InitDB()`. Environment `"production"` pakai JSON encoder (level Info); selain itu console encoder (level Debug). Jangan panggil `Init` di tempat lain — package lain cukup akses via `logger.L()` atau `logger.FromContext(ctx)`.

```go
zapLogger := logger.Init(cfg.AppEnv, cfg.AppDebug)
defer logger.Sync() // flush buffer sebelum process exit
```

## 2. Context Propagation (WAJIB, selaras coding-convention.md §1)

`internal/shared/middleware.RequestLogger()` men-generate/forward `request_id` (header `X-Request-Id`) dan menyisipkan logger ber-request_id ke `context.Context` lewat `logger.WithContext`. Setiap layer yang menerima `ctx context.Context` (application, domain service, adapter) **WAJIB** ambil logger dari situ, bukan dari `logger.L()` langsung, supaya semua log baris untuk satu request bisa dikorelasikan lewat `request_id`:

```go
logger.FromContext(ctx).Error("failed to find user by email during login", zap.Error(err))
```

`logger.L()` global hanya dipakai di tempat yang belum punya `ctx` request-scoped (lifecycle server: startup, shutdown, koneksi DB awal — lihat `cmd/api/server.go`, `internal/shared/database/postgres.go`).

## 3. Log Sekali di Boundary — Jangan Duplikasi di Tiap Layer

**DILARANG** menaruh `logger.FromContext(ctx).Error(...)` di setiap `return nil, err` passthrough (domain repo, application service). Kalau tiap layer log ulang error yang sama, satu kegagalan DB bisa menghasilkan 3-4 baris log identik dan bikin observability berisik.

**WAJIB** log persis di titik keputusan:
- **Transport/HTTP layer**, tepat sebelum error diubah jadi response `5xx` generik (lihat pola `serverError(...)` di `internal/interfaces/http/organization/handler.go` dan `internal/interfaces/http/employee/handler.go`). Ini titik yang benar karena di sinilah detail asli (`err`) tentang hilang diganti pesan generik ke client — kalau tidak dicatat di sini, hilang selamanya.
- **Titik swallow di Application layer**, HANYA kalau error asli sengaja "dibungkus ulang" jadi sentinel error lain sebelum dikembalikan (mis. `auth/application/service.go` — kegagalan DB di `FindByEmail`/`FindByID` dibungkus jadi `ErrInvalidCredentials`/`ErrInvalidToken` supaya tidak bocor detail infra ke client sebagai pesan auth; tanpa log di titik ini, root cause hilang karena error asli tidak pernah nyampe ke boundary manapun).

Passthrough murni (`return nil, err` tanpa reklasifikasi) **TIDAK** perlu log tambahan — cukup dicatat sekali di boundary transport.

## 4. Jangan Bocorkan Detail Error ke Client (selaras persistence-convention.md §4)

Pola wajib di titik boundary:

```go
func serverError(c fiber.Ctx, err error, message string) error {
    logger.FromContext(c.Context()).Error(message, zap.Error(err))
    return response.Error(c, fiber.StatusInternalServerError, message, nil)
}
```

`message` yang dikirim ke client harus generik (mis. `"Failed to create department"`), **bukan** `err.Error()`. Detail asli (`zap.Error(err)`) hanya masuk log server.

## 5. Structured Fields, Bukan String Interpolation

Pakai field zap (`zap.String`, `zap.Int`, `zap.Duration`, `zap.Error`), jangan `fmt.Sprintf` ke dalam message log:

```go
// ❌ DILARANG
logger.L().Info(fmt.Sprintf("connected to db %s on %s:%s", name, host, port))

// ✅ WAJIB
logger.L().Info("connected to database successfully",
    zap.String("db_name", name), zap.String("host", host), zap.String("port", port))
```

## 6. Checklist Review

- [ ] Tidak ada `log.Println`/`log.Printf` baru di luar bootstrap sebelum `logger.Init()`.
- [ ] Fungsi yang terima `ctx` pakai `logger.FromContext(ctx)`, bukan `logger.L()` langsung.
- [ ] Tidak ada log duplikat di tiap layer untuk error yang sama — cukup di boundary transport atau titik reklasifikasi error.
- [ ] Response `5xx` ke client tetap pesan generik; `err` asli hanya di `zap.Error(err)`.
- [ ] Log pakai structured field, bukan string yang sudah di-`Sprintf`.
