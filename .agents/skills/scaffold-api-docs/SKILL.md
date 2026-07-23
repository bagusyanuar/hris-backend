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
4. Membaca `application/dto.go` domain terkait untuk memetakan tipe tiap field response DTO — field pointer (`*string`, `*bool`, dst.) = **nullable**, field value (`string`, `bool`, `int`, slice) = **selalu ada meski zero-value**. Ini sumber kebenaran buat langkah nullable-marking di bawah, jangan nebak dari nama field.

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
- **Nullable vs Always-Present Fields (WAJIB)**: Response sukses (`200`/`201`) TIDAK BOLEH cuma punya `example` mentah — definisikan `components.schemas.<Entity>Response` (dan `<Entity>ListResponse` buat endpoint List, reuse `items` + `meta`) lalu `$ref` dari tiap `responses.<code>.content.application/json.schema`, sejajar `example` yang udah ada. Di dalam schema itu:
  - Field pointer di DTO (§Investigasi Awal poin 4) → `nullable: true` + deskripsi singkat kapan null-nya.
  - Field value/slice → masuk `required: [...]`, TIDAK nullable. Slice WAJIB dijelaskan eksplisit "selalu array, gak pernah `null`" di `description` (selaras `pkg/response` — larangan `data` slice `null`).
- **Endpoint Purpose Description (WAJIB)**: Tiap operation WAJIB punya `description` (terpisah dari `summary`) + `operationId` unik camelCase (`listCompanies`, `createBranch`, dst):
  - `description` = *purpose* (kenapa endpoint ini ada) + *side-effect non-obvious* (mis. auto-demote, cascade behavior) + *kapan pakai vs endpoint mirip lainnya*. 2-4 kalimat. JANGAN restate `parameters`/`schema` yang udah ada.
  - `operationId` buat konsumen agentic (OpenAPI-to-tool/MCP-style) — biar AI gak nebak dari `summary`.
  - Kalau behavior endpoint berubah nanti (ubah handler/service logic), `description` WAJIB direview ulang — satu paket sama kewajiban bump SemVer di bawah, jangan biarin drift dari behavior asli.
- **Registrasi Docker Compose**: Ketika membuat file spesifikasi YAML baru (domain baru), AI **WAJIB** mendaftarkannya ke dalam environment variable `URLS` di file [docker-compose.yaml](file:///Users/dystopia/go/hris-backend/docs/api/swagger/docker-compose.yaml) agar langsung dapat diakses lewat dropdown Swagger UI lokal.
- **Versioning (Penting!)**: Jika terjadi perubahan, perbaikan, atau penambahan endpoint pada API contract, AI **WAJIB** menaikkan versi API contract menggunakan format **Semantic Versioning (SemVer, e.g., MAJOR.MINOR.PATCH)** pada field `info.version` di file Swagger YAML, serta menyesuaikan info versi pada bagian dokumentasi Bruno jika relevan.

### 2. Bruno Collection
Lokasi: `docs/api/bruno/<Domain>/<EndpointName>.bru`
- Buat folder dengan awalan huruf kapital, lalu di dalamnya buat file `.bru`.
- Tambahkan blok request (`get`, `post`, dsb.) lengkap dengan URL parameter dan body (jika ada). Pastikan penamaan variabel selaras (menggunakan `:id`).
- Gunakan `vars` seperti `{{baseUrl}}` agar fleksibel.
- Wajib menambahkan blok `docs { ... }` (menggunakan format Markdown) di dalam file `.bru` yang minimal memuat:
  - Baris pembuka "Tujuan: ..." — sejajar `description` Swagger endpoint yang sama (purpose + side-effect non-obvious).
  - Deskripsi singkat endpoint.
  - Penjelasan field *Request Body* beserta status `required`/`optional`.
  - **Expected Responses**: Contoh block JSON aktual untuk SEMUA status HTTP yang ditemukan saat investigasi `handler.go`. Harus mencakup sukses (200/201) dan seluruh error (termasuk contoh format validator 422).
  - **Nullable vs Always-Present Fields**: Persis sejajar Swagger — di bawah tiap contoh JSON sukses, tambah 2 baris ringkas: "Field nullable (bisa `null`): ..." dan "Field yang SELALU ada meski zero-value: ...".
