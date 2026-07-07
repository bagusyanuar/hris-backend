# Project Rules: HRIS Backend DDD

Semua instruksi coding, pembuatan file, refactoring, dan pengerjaan tugas di workspace `/Users/dystopia/go/hris-backend` harus mematuhi aturan arsitektur yang tertera pada file [GEMINI.md](file:///Users/dystopia/go/hris-backend/GEMINI.md).

## Aturan Utama:
1. **Domain-Driven Design (DDD)**: Patuhi batasan layer (`internal/domain`, `internal/application`, `internal/infrastructure`, `internal/interfaces`).
2. **Tanpa CQRS**: Gunakan application service terpadu di `internal/application/` (`service.go` untuk read & write logic), jangan memisahkannya menjadi command & query.
3. **No External Imports / ORM Tags in Domain**: Domain layer tidak boleh mengimpor framework atau database library. Entity Domain **tidak boleh** memiliki tag GORM (seperti `gorm:"..."`). Jika perlu pemetaan ORM, buat struct Model khusus di layer **Infrastructure** lalu map dari/ke Entity Domain.
4. **Database Migrations via golang-migrate**: Semua perubahan skema database harus menggunakan file migrasi SQL (`.up.sql` dan `.down.sql`) melalui perintah `make migrate-create`. Jangan menggunakan GORM `AutoMigrate` di production/code utama.
5. **Context & Error Handling**: Selalu gunakan `context.Context` dan gunakan custom domain error.
