# Automation & Code Generation Planning

Dokumen ini berisi hasil *brainstorming* mengenai teknik dan tools otomasi yang direncanakan untuk diimplementasikan pada project HRIS Backend ini. Otomasi ini bertujuan untuk meningkatkan produktivitas, menjaga kualitas kode, dan mengurangi pekerjaan repetitif (terutama dalam arsitektur Domain-Driven Design).

## 1. Code Generation
Karena DDD banyak menggunakan *Interface* (terutama di layer repository), *code generation* akan sangat mempercepat proses pengembangan.

*   **Dependency Injection Automation:** [google/wire](https://github.com/google/wire)
    *   **Tujuan:** Mengotomatisasi *wiring* dependencies (menyambungkan Repository ke Service, Service ke Handler) saat inisialisasi aplikasi.
    *   **Benefit:** Mencegah `main.go` atau file bootstrap menjadi terlalu panjang dan kompleks. Menghindari *runtime error* akibat dependency yang lupa di-inject, karena pengecekan dilakukan saat kompilasi.

*   **Mock Generator:** [vektra/mockery](https://github.com/vektra/mockery) atau [uber-go/mock](https://github.com/uber-go/mock)
    *   **Tujuan:** Men-generate file *mock* untuk *Interface* secara otomatis.
    *   **Benefit:** Mempercepat penulisan *Unit Test* pada layer Application (Service), karena developer tidak perlu menulis *mock repository* secara manual.

## 2. Testing Automation
Fokus pada keandalan testing selain *Unit Test* standar.

*   **Integration Test dengan Real Database:** [testcontainers-go](https://github.com/testcontainers/testcontainers-go)
    *   **Tujuan:** Melakukan testing query database (GORM) menggunakan real database, bukan sekadar *mock*.
    *   **Benefit:** Saat *integration test* dijalankan, sistem akan otomatis melakukan *spin up* container PostgreSQL (via Docker), menjalankan migrasi, mengeksekusi test, dan menghancurkan container setelah selesai. Hasil test lebih valid.

*   **Test Coverage Enforcer**
    *   **Tujuan:** Memastikan *code coverage* tidak turun di bawah standar (misal: 80%).
    *   **Benefit:** Mendorong tim untuk selalu menulis *unit test* setiap membuat fitur baru.

## 3. API Documentation Automation

*   **Swagger/OpenAPI Generator:** [swaggo/swag](https://github.com/swaggo/swag)
    *   **Tujuan:** Men-generate dokumentasi Swagger otomatis dari komentar kode (meskipun aturan awal meminta penulisan YAML manual, ini dapat dipertimbangkan jika dokumentasi manual terasa memberatkan).
    *   **Benefit:** Mengurangi risiko dokumentasi API (Swagger/Bruno) tidak sinkron dengan implementasi kode aktual.

## 4. Security & Quality Scanning

*   **Go Vulnerability Scanner:** [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
    *   **Tujuan:** Memeriksa CVE (*Common Vulnerabilities and Exposures*) pada dependensi di `go.mod`.
    *   **Benefit:** Mengamankan aplikasi dari *library* pihak ketiga yang rentan secara otomatis.

## 5. CI/CD & Workflow Automation

*   **GitHub Actions / GitLab CI**
    *   **Tujuan:** Membuat *pipeline* terotomatisasi yang berjalan setiap ada *Pull Request* (PR) atau dorongan kode ke *branch* utama.
    *   **Benefit:** Otomatis menjalankan `make lint`, `go test`, `govulncheck`, dan proses *build* untuk memastikan kode yang akan di-merge dalam kondisi sehat.

*   **Commitlint & Semantic Release**
    *   **Tujuan:** Mengotomatisasi penulisan *Changelog* dan penentuan versi rilis.
    *   **Benefit:** Melanjutkan praktik *Conventional Commits* dengan secara otomatis menghasilkan catatan rilis (`CHANGELOG.md`) dan *tagging* versi (v1.0.0, dll) ketika kode digabungkan.

---

## 6. AI-Driven Development Flow
Untuk mempercepat pengembangan modul baru, kita mengadopsi alur kerja berbasis AI (AI-Driven Development). Alur ini memisahkan fase perancangan (*design*) dan penulisan kode (*coding*).

### Alur Kerja (Workflow):
1. **Brainstorming Ideation:** Diskusi awal mengenai konsep modul baru (misal: modul Cuti/Leave).
2. **Technical Document Generation:** AI Agent akan membuat *Product Requirements Document* (PRD) / *Technical Document* terstruktur (di dalam folder `docs/technical/`) berdasarkan hasil diskusi. Dokumen mencakup skema DB, entitas, aturan bisnis, dan kontrak API.
3. **Review & Approval:** Developer (Man-in-the-loop) mereview dan menyetujui dokumen teknis tersebut.
4. **Automated Code Generation:** AI Agent secara otomatis mengeksekusi pembuatan kode (menggunakan skill `scaffold-domain` & `scaffold-api-docs`) berdasarkan dokumen yang telah disetujui, mencakup pembuatan *Entity, Repository, Service, Handler*, dan dokumentasi API.

---

## 7. Aturan Dokumentasi (Living Documentation)
Semua *Technical Document* (PRD) wajib disimpan di dalam direktori `docs/technical/`.

Apabila di kemudian hari terjadi perubahan konsep bisnis atau arsitektur pada suatu modul, **JANGAN** membuat file versi baru (misal: `employee-v2.md`). 

### Konsep Living Document:
* **Single Source of Truth (SSOT):** Langsung ubah/update isi file dokumen lama agar dokumen tersebut selalu mencerminkan *current state* dari kode sumber.
* **Git History:** Riwayat perubahan secara detil (per baris) akan otomatis dilacak oleh sistem *version control* (Git).
* **Revision History Section:** Sebagai catatan yang ramah-baca, tambahkan *section* `## Revision History` di bagian paling bawah setiap dokumen teknis untuk mencatat secara ringkas tanggal dan alasan perubahan konsep.

**Contoh Format Revision History di bawah Dokumen:**
```markdown
## Revision History
* **v1.0.0 (08 Juli 2026):** Inisialisasi awal konsep modul.
* **v1.1.0 (15 Agustus 2026):** Penambahan field `BankName` atas request tim Finance.
```
