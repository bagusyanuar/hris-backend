---
name: scaffold-api-docs
description: Scaffold API documentation using Swagger YAML and Bruno collections per domain.
---

# Scaffolding API Documentation

Gunakan skill ini ketika selesai membuat endpoint API baru atau ada request untuk membuat dokumentasi API.

## Investigasi Awal (Wajib!)
Sebelum menulis file dokumentasi, AI **WAJIB**:
1. Membaca `handler.go` yang terkait untuk memetakan secara *exhaustive* SEMUA *error code* yang di-*return* (misal: 200, 201, 400, 404, 409, 422, 500).
2. Mengecek rute yang sudah ada menggunakan command `ls` atau alat pencari (`grep`) di folder Swagger/Bruno agar **TIDAK MENDUPLIKASI** endpoint (misal membuat `GetEmployee.bru` dan `GetEmployeeDetail.bru` secara bersamaan).
3. Memastikan variabel parameter URL konsisten antara Swagger (`{id}`) dan Bruno (`:id`).

## Aturan Dokumentasi API
Setiap dokumentasi API harus dibagi berdasarkan **Domain** (bukan *monolithic* di satu file) untuk kemudahan tim Frontend. 
Ada 2 jenis dokumentasi yang wajib dibuat secara berbarengan:

20: ### 1. Swagger OpenAPI (YAML)
21: Lokasi: `docs/api/swagger/<domain>.yaml`
22: - Buat file `.yaml` dengan nama domain yang bersangkutan (contoh: `auth.yaml`, `employee.yaml`).
23: - Di dalam file, definisikan `openapi: 3.0.3` dan `info`.
24: - Definisikan tipe otorisasi di `components.securitySchemes`.
25: - Setiap endpoint di `paths` harus memiliki:
26:   - `summary` & `tags` (tag diisi nama domain).
27:   - `requestBody` jika ada, pastikan menandai field mana saja yang `required`.
28:   - Berikan `example` untuk setiap `properties`.
29:   - Daftarkan **seluruh kemungkinan respons** (`responses`) yang ditemukan saat fase investigasi di `handler.go`. Pastikan JSON example mengikuti standar `pkg/response`.
30:   - **Aturan 422 Validation Error**: Khusus untuk status 422, wajib menyertakan contoh *array format error* dari Gofiber/Validator.
31:   - **Hapus Rute Usang**: Pastikan untuk menghapus rute *dummy/monolith* jika sistem sudah menerapkan pola modular/Progressive Save.
32: - **Registrasi Docker Compose**: Ketika membuat file spesifikasi YAML baru (domain baru), AI **WAJIB** mendaftarkannya ke dalam environment variable `URLS` di file [docker-compose.yaml](file:///Users/dystopia/go/hris-backend/docs/api/swagger/docker-compose.yaml) agar langsung dapat diakses lewat dropdown Swagger UI lokal.
33: - **Versioning (Penting!)**: Jika terjadi perubahan, perbaikan, atau penambahan endpoint pada API contract, AI **WAJIB** menaikkan versi API contract menggunakan format **Semantic Versioning (SemVer, e.g., MAJOR.MINOR.PATCH)** pada field `info.version` di file Swagger YAML, serta menyesuaikan info versi pada bagian dokumentasi Bruno jika relevan.




### 2. Bruno Collection
Lokasi: `docs/api/bruno/<Domain>/<EndpointName>.bru`
- Buat folder dengan awalan huruf kapital, lalu di dalamnya buat file `.bru`.
- Tambahkan blok request (`get`, `post`, dsb.) lengkap dengan URL parameter dan body (jika ada). Pastikan penamaan variabel selaras (menggunakan `:id`).
- Gunakan `vars` seperti `{{baseUrl}}` agar fleksibel.
- Wajib menambahkan blok `docs { ... }` (menggunakan format Markdown) di dalam file `.bru` yang minimal memuat:
  - Deskripsi singkat endpoint.
  - Penjelasan field *Request Body* beserta status `required`/`optional`.
  - **Expected Responses**: Contoh block JSON aktual untuk SEMUA status HTTP yang ditemukan saat investigasi `handler.go`. Harus mencakup sukses (200/201) dan seluruh error (termasuk contoh format validator 422).
