---
name: scaffold-api-docs
description: Scaffold API documentation using Swagger YAML and Bruno collections per domain.
---

# Scaffolding API Documentation

Gunakan skill ini ketika selesai membuat endpoint API baru atau ada request untuk membuat dokumentasi API.

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
  - Daftarkan **seluruh kemungkinan respons** (`responses`) dengan contoh JSON, seperti `200`, `201`, `400` (Bad Request), `401` (Unauthorized), dsb.

### 2. Bruno Collection
Lokasi: `docs/api/bruno/<Domain>/<EndpointName>.bru`
- Buat folder dengan awalan huruf kapital, lalu di dalamnya buat file `.bru`.
- Tambahkan blok request (`get`, `post`, dsb.) lengkap dengan URL parameter dan body (jika ada).
- Gunakan `vars` seperti `{{baseUrl}}` agar fleksibel.
- Wajib menambahkan blok `docs { ... }` (menggunakan format Markdown) di dalam file `.bru` yang minimal memuat:
  - Deskripsi singkat endpoint.
  - Penjelasan field *Request Body* beserta status `required`/`optional`.
  - **Expected Responses**: Contoh block JSON aktual untuk status code `200`, `400`, `401`, dan seterusnya sesuai dengan yang didefinisikan di backend.
