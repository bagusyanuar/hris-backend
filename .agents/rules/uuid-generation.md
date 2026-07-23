# UUID Generation (Primary Key)

Semua entitas yang menggunakan UUID sebagai Primary Key wajib mengimplementasikan pola *auto-generate* UUID pada dua *layer* berikut:
1. **Domain Layer (`entity.go`)**: Di dalam *constructor function* (`NewEntityName(...)`), pastikan ada pengecekan jika ID kosong, maka diisi dengan UUID baru (`if id == "" { id = uuid.New().String() }`).
2. **Adapter Layer (`models.go`)**: Tambahkan *hook* GORM `BeforeCreate` pada model yang bersangkutan untuk mengisi `m.ID` dengan UUID baru jika masih kosong.

## Single Source of Generation (WAJIB)

UUID **hanya boleh** di-generate di satu tempat: **Domain Layer (constructor)**, dengan `BeforeCreate` sebagai jaring pengaman di Adapter. **Application Layer DILARANG** meng-generate UUID sendiri lalu mengopernya ke constructor. Terapkan pola ini **seragam di semua domain** (jangan Employee generate di domain sementara Organization generate di service).

```go
// ❌ DILARANG — application layer generate UUID
id := uuid.New().String()
dept, _ := organization.NewDepartment(id, ...)

// ✅ WAJIB — constructor yang generate
dept, _ := organization.NewDepartment(req.Code, req.Name, ...) // ID diisi di dalam constructor
```

> Catatan: `BeforeCreate` hanya jalan pada `Create()`, bukan pada `Save()`-yang-menjadi-`Update`. Lihat [persistence-convention.md](persistence-convention.md) §1 — pakai `Create()`/`OnConflict`, bukan `Save()`, agar hook & generation berperilaku benar.
