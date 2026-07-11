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

### 1. Swagger OpenAPI (YAML)
Lokasi: `docs/api/swagger/<domain>.yaml`
- Buat file `.yaml` dengan nama domain yang bersangkutan (contoh: `auth.yaml`, `employee.yaml`).
- Di dalam file, definisikan `openapi: 3.0.3` dan `info`.
- Definisikan tipe otorisasi di `components.securitySchemes`.
- Setiap endpoint di `paths` harus memiliki:
  - `summary` & `tags` (tag diisi nama domain).
  - `requestBody` jika ada, pastikan menandai field mana saja yang `required`.
  - Berikan `example` untuk setiap `properties`.
  - Daftarkan **seluruh kemungkinan respons** (`responses`) yang ditemukan saat fase investigasi di `handler.go`. Pastikan JSON example mengikuti standar `pkg/response`.
  - **Aturan 422 Validation Error**: Khusus untuk status 422, wajib menyertakan contoh *array format error* dari Gofiber/Validator.
  - **Hapus Rute Usang**: Pastikan untuk menghapus rute *dummy/monolith* jika sistem sudah menerapkan pola modular/Progressive Save.

### 2. Bruno Collection
Lokasi: `docs/api/bruno/<Domain>/<EndpointName>.bru`
- Buat folder dengan awalan huruf kapital, lalu di dalamnya buat file `.bru`.
- Tambahkan blok request (`get`, `post`, dsb.) lengkap dengan URL parameter dan body (jika ada). Pastikan penamaan variabel selaras (menggunakan `:id`).
- Gunakan `vars` seperti `{{baseUrl}}` agar fleksibel.
- Wajib menambahkan blok `docs { ... }` (menggunakan format Markdown) di dalam file `.bru` yang minimal memuat:
  - Deskripsi singkat endpoint.
  - Penjelasan field *Request Body* beserta status `required`/`optional`.
  - **Expected Responses**: Contoh block JSON aktual untuk SEMUA status HTTP yang ditemukan saat investigasi `handler.go`. Harus mencakup sukses (200/201) dan seluruh error (termasuk contoh format validator 422).
