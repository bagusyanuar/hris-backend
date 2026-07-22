# Dokumentasi API (Swagger & Bruno)

Setiap endpoint API yang dibuat harus didokumentasikan di dua tempat dengan aturan "Split by Domain" (per file untuk setiap domain, bukan monolithic):

**ATURAN WAJIB (STRICT RULES):**
- **Anti-Duplikasi:** AI WAJIB mengecek eksistensi rute/file yang sudah ada (menggunakan fitur *search* atau `ls`) sebelum meng-generate endpoint baru untuk menghindari duplikasi.
- **Konsistensi Variabel URL:** Parameter di URL harus konsisten. Pada dokumentasi Swagger gunakan `{id}`, pada Bruno gunakan `:id`. Jangan gunakan penamaan *custom* yang bisa membingungkan FE (seperti `:employee_id` atau `{{employee_id}}`).
- **Exhaustive Error Responses:** Dokumentasi *Response* tidak boleh asal-asalan. AI WAJIB membaca file `handler.go` untuk mendaftar SEMUA HTTP Error Code yang mungkin terjadi (`200`, `201`, `400`, `404`, `409`, `422`, `500`). Khusus untuk `422 Unprocessable Entity`, wajib menyertakan contoh array error dari validator.
- **API Contract Versioning:** Setiap kali terjadi perubahan, perbaikan, atau penambahan endpoint baru pada API contract, AI **WAJIB** menaikkan versi API contract menggunakan format **Semantic Versioning (SemVer, e.g., MAJOR.MINOR.PATCH)** pada field `info.version` di spesifikasi Swagger YAML serta memperbarui info versi di koleksi Bruno jika dibutuhkan.

1. **Swagger OpenAPI (YAML)** di `docs/api/swagger/<domain>.yaml`. Wajib mendeskripsikan `requestBody`, `responses` super komplit, `required` fields, dan `example`. Pastikan untuk menghapus rute *dummy/monolith* jika sistem sudah menerapkan pola modular/Progressive Save. AI **WAJIB** mendaftarkan berkas YAML baru ke dalam environment `URLS` pada [docker-compose.yaml](file:///Users/dystopia/go/hris-backend/docs/api/swagger/docker-compose.yaml) agar langsung tampil di Swagger UI lokal.
2. **Bruno Collection** di `docs/api/bruno/<Domain>/<EndpointName>.bru`. Wajib menyertakan blok `docs { ... }` (menggunakan format Markdown) yang menjelaskan endpoint, `required` properties, dan seluruh variasi *Expected Responses* persis seperti yang tercatat di Swagger.
