# UUID Generation (Primary Key)

Semua entitas yang menggunakan UUID sebagai Primary Key wajib mengimplementasikan pola *auto-generate* UUID pada dua *layer* berikut:
1. **Domain Layer (`entity.go`)**: Di dalam *constructor function* (`NewEntityName(...)`), pastikan ada pengecekan jika ID kosong, maka diisi dengan UUID baru (`if id == "" { id = uuid.New().String() }`).
2. **Infrastructure Layer (`models.go`)**: Tambahkan *hook* GORM `BeforeCreate` pada model yang bersangkutan untuk mengisi `m.ID` dengan UUID baru jika masih kosong.
