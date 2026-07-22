# Konvensi Kode Go

1. **Gunakan Context**: Selalu sertakan `context.Context` sebagai argumen pertama pada fungsi-fungsi di layer application, domain repository, dan infrastructure (misal `FindByID(ctx context.Context, id string)`).
2. **Error Handling**:
   - Tangani error sedini mungkin.
   - Jangan abaikan error (`_ = someFunc()`).
   - Gunakan custom domain error (misal `ErrEmployeeNotFound = errors.New("employee not found")`) di layer domain agar interface layer bisa memetakan status HTTP dengan tepat (misal 404 Not Found).
3. **Dependency Injection (Google Wire)**: Proyek ini wajib menggunakan `google/wire` untuk injeksi dependensi secara compile-time. Semua dependensi (Repository, Service, Handler) di-registrasikan ke dalam `wire.ProviderSet` di dalam file `internal/di/wire.go`. File `cmd/api/server.go` akan bersih karena cukup memanggil `di.InitializeAPI(s.db, tokenGenerator)`.
4. **Cross-Domain Communication (Bounded Contexts)**: Untuk saat ini, komunikasi antar modul/domain dilakukan secara langsung (Synchronous) melalui injeksi *Application Service* modul lain menggunakan `google/wire` (bukan menggunakan Message Broker/Event Bus). Hindari injeksi *Repository* modul lain secara langsung ke dalam *Service*; selalu gunakan *Application Service* modul tersebut sebagai jembatan/API internal.
5. **Configuration**: Load konfigurasi dari environment variables atau config file sekali saja di `cmd/api/main.go` menggunakan library seperti `viper` atau `envconfig`, lalu teruskan struct config ke service yang membutuhkan.
6. **Acronym Naming Consistency**: Ikuti standar Go untuk akronim (misal ID, HTTP, URL, API). Jika menggunakan *CamelCase* untuk akronim lokal seperti KTP atau PTKP, pastikan konsisten secara presisi (termasuk huruf besar-kecil) antara *Domain Entity*, *GORM Model*, dan *DTO*. Contoh: gunakan `PtkpStatus` di semua layer, jangan dicampur dengan `PTKPStatus`.
7. **Mandatory Build Check**: AI WAJIB menjalankan perintah `go build ./...` di terminal setiap kali selesai men-generate atau memodifikasi kumpulan file `.go`. Tujuannya untuk menangkap *syntax error*, salah *type*, atau *import* yang hilang sebelum melaporkan pekerjaan selesai kepada user.
