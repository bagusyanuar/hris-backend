# Konvensi Persistensi (Repository, Transaction, Data Integrity)

Aturan ini mengatur **cara menyimpan & mengambil data** di Infrastructure Layer dan **siapa yang memegang transaksi**. Tujuannya mencegah bug integritas data yang tidak tertangkap oleh compiler maupun linter (mis. silent no-op saat update). Semua aturan di bawah bersifat **WAJIB (STRICT)**.

---

## 1. DILARANG `db.Save()` untuk Semantik Insert-or-Update by Non-PK Key

`gorm.DB.Save()` **melakukan UPDATE (semua field) jika primary key terisi**, dan hanya `INSERT` jika PK kosong. Karena pola project ini meng-generate UUID di domain constructor (PK **selalu** terisi), memanggil `Save()` pada data yang belum ada di DB akan menghasilkan `UPDATE ... WHERE id = <uuid-baru>` → **0 rows affected, TANPA error** → data hilang diam-diam.

**DILARANG:**
```go
// BUG: model.ID = UUID baru dari constructor → Save() jadi UPDATE 0 rows → tidak tersimpan
func (r *Repo) SavePersonalData(ctx context.Context, d *employee.PersonalData) error {
    model := models.PersonalDataFromDomain(d)
    return r.db.WithContext(ctx).Save(model).Error
}
```

**WAJIB — pilih salah satu, konsisten per entity:**

**(a) Upsert eksplisit by business key** (untuk relasi 1-1 yang di-key oleh kolom selain PK, mis. `employee_id`):
```go
import "gorm.io/gorm/clause"

func (r *Repo) SavePersonalData(ctx context.Context, d *employee.PersonalData) error {
    model := models.PersonalDataFromDomain(d)
    return r.db.WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "employee_id"}},
            UpdateAll: true,
        }).
        Create(model).Error
}
```

**(b) Fetch-then-reuse-PK** — cari row existing, kalau ada reuse PK-nya lalu `Save`/`Updates`, kalau tidak `Create`. Konsekuensinya domain constructor tidak boleh memaksa generate UUID pada alur update (lihat rule UUID).

**Pemetaan operasi → method GORM yang benar:**
| Maksud | Method WAJIB | Catatan |
|--------|-------------|---------|
| Insert baru | `Create()` | PK boleh kosong → di-generate `BeforeCreate` |
| Update record yang PK-nya sudah pasti ada di DB | `Updates()` / `Save()` | Cek `RowsAffected` bila perlu deteksi not-found |
| Insert-or-update by unique non-PK key | `Clauses(clause.OnConflict{...}).Create()` | Jangan `Save()` |
| Ganti-seluruh koleksi anak (banks, educations) | `Transaction`: delete-by-parent lalu `Create` batch | Wajib dalam satu transaksi |

---

## 2. Transaction Ownership Ada di Application Layer

Sesuai [architecture.md](architecture.md) §B, **transaksi dimiliki Application Layer**, bukan repository.

- **DILARANG** repository membuka `db.Transaction(...)` sendiri untuk mengoordinasi beberapa operasi bisnis. Repository hanya melakukan satu unit data-access.
- **WAJIB** menyediakan abstraksi transaksi (Unit of Work / `TxManager`) yang di-inject ke application service. Service yang membuka & menutup transaksi; repository membaca handle transaksi dari `context.Context`.

Kontrak minimal:
```go
// domain (atau shared kernel) — abstraksi, bukan detail GORM
type TxManager interface {
    Do(ctx context.Context, fn func(ctx context.Context) error) error
}
```

Semua langkah yang harus atomik (mis. cek duplikat + insert, atau update induk + ganti anak) **WAJIB** dibungkus dalam satu `TxManager.Do`. Jangan mengandalkan urutan beberapa call repo tanpa transaksi (rawan TOCTOU / partial write).

---

## 3. Kontrak Repository: Not-Found WAJIB Sentinel Error

Method `FindXxx` **DILARANG** mengembalikan `(nil, nil)` untuk menandai "tidak ditemukan". Pola nil-nil memaksa caller mengingat pengecekan manual dan mudah menimbulkan nil-dereference.

**DILARANG:**
```go
func (r *Repo) FindByKTP(ctx context.Context, ktp string) (*employee.PersonalData, error) {
    ...
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil // ❌ leaky contract
    }
}
```

**WAJIB** — kembalikan sentinel domain error, biarkan caller `errors.Is`:
```go
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, employee.ErrPersonalDataNotFound
}
```

Repository juga **WAJIB** menerjemahkan pelanggaran constraint DB (mis. unique violation) menjadi sentinel domain error yang tepat (mis. `ErrKTPDuplicate`) supaya interface layer bisa memetakannya ke HTTP 409, bukan 500.

---

## 4. Error Internal DILARANG Bocor ke Client

Pada response `5xx`, **DILARANG** mengirim `err.Error()` mentah ke client (berpotensi membocorkan nama kolom/SQL/driver — information disclosure).

- **WAJIB** log error lengkap (structured) di sisi server, lalu balikan pesan generik ke client (mis. `"internal server error"`).
- `err.Error()` hanya boleh dikirim untuk **sentinel domain error** yang memang aman & bermakna bagi user (mis. `"employee not found"`).

Sinkron dengan aturan Error Handling di [coding-convention.md](coding-convention.md).

---

## 5. Ringkasan Checklist (untuk Review / PR)

- [ ] Tidak ada `db.Save()` untuk insert-or-update by non-PK key → pakai `OnConflict` upsert atau fetch-reuse-PK.
- [ ] Operasi multi-langkah yang harus atomik dibungkus `TxManager.Do` di application layer.
- [ ] Repository tidak membuka transaksi bisnis sendiri.
- [ ] Tidak ada `FindXxx` yang `return nil, nil`; not-found = sentinel error.
- [ ] Unique/constraint violation dipetakan ke sentinel error (→ HTTP 409), bukan 500.
- [ ] Response 5xx tidak membocorkan `err.Error()` mentah.
