# HRIS Documentation Hub

Selamat datang di pusat dokumentasi HRIS Backend. Repositori ini menerapkan konsep **Docs-as-Code**, di mana dokumentasi arsitektur, kebutuhan bisnis, dan spesifikasi teknis (*tech specs*) hidup berdampingan secara permanen dengan *source code*.

---

## 📁 Struktur Direktori Utama

Dokumentasi dipisahkan secara tegas menjadi tiga pilar utama:

### 1. `docs/requirement/` (Product Requirements Document - PRD)
Fokus pada **WHAT** dan **WHY**. Berisi kebutuhan bisnis, aturan main, dan *scope* dari sudut pandang *Product Management*.
- **Tujuan:** Menjadi *Single Source of Truth* bagi tim Bisnis, QA, Frontend, dan Backend.
- **Aturan Pembuatan:** Setiap PRD **wajib** mengikuti pedoman **6 Pilar Enterprise**:
  1. *The "Why"* (Tujuan & Dampak)
  2. *Scope & Out-of-Scope* (Batasan)
  3. *User Roles & Permissions* (Hak Akses)
  4. *Acceptance Criteria* (Skenario Given-When-Then)
  5. *Technical & Architectural Constraints* (Aturan Engineering)
  6. *Dependencies* (Ketergantungan Modul)

### 2. `docs/technical/` (Technical Specs & RFC)
Fokus pada **HOW**. Berisi cetak biru arsitektur dari sudut pandang *Engineering*.
Setiap modul/domain memiliki sub-foldernya sendiri (misal: `docs/technical/employee/`) yang berisi pecahan dokumen:
- `tech-spec.md` : *Request for Comments* (RFC) berisi rancangan arsitektur, API *contracts*, dan skema database (ERD).
- `user-stories.md` : Alur logika sistem (menggunakan *Mermaid Sequence diagram*).
- `decision-log.md` : Architecture Decision Records (ADR) untuk mencatat **mengapa** suatu keputusan teknis diambil dan riwayat perubahannya.
- (Opsional) `data-dictionary.md`, `infrastructure.md`, `test-plan.md`.

### 3. `docs/databases/` (Database Schema & DBML)
Berisi rancangan struktur *database* fisik dalam format DBML (*Database Markup Language*).
- **Tujuan:** Memudahkan *generate* ERD visual menggunakan *tools* seperti dbdiagram.io.
- Setiap modul/domain wajib memiliki file `.dbml` tersendiri (misal: `docs/databases/employee.dbml`).

---

## 🤖 Cara Menggunakan Automasi AI (Scaffolding Docs)

Proyek ini telah dilengkapi dengan *Workflow* AI pintar untuk membantu *developer* merancang arsitektur. 
Jika Anda ingin merancang fitur atau modul baru, gunakan **Slash Command** berikut di dalam panel Chat IDE Anda:

**`/scaffold-docs`**

Contoh penggunaan dengan *prompt* langsung:
> `/scaffold-docs Tolong buatkan dokumentasi untuk modul Attendance (Absensi). Fitur utamanya adalah clock-in/out pakai koordinat GPS.`

**Alur Kerja AI (`/scaffold-docs`):**
1. **Fase Bisnis:** AI akan membuatkan **PRD** di folder `docs/requirement/` dan berhenti sementara untuk meminta persetujuan (*approval*) Anda.
2. **Fase Engineering:** Setelah PRD di-klik *approve*, AI akan otomatis menerjemahkan logika bisnis tersebut menjadi file-file arsitektur teknis lengkap di dalam `docs/technical/<domain>/`.

---

## 🔄 Panduan Versioning (Jika Ada Perubahan)

Semua dokumen di sini adalah *living documents* (dokumen hidup yang terus *update*).
- **Update di Tempat:** Jika ada perubahan aturan bisnis di kemudian hari, **jangan** membuat file baru (misal: `employee_v2.md`). Selalu perbarui file aslinya.
- **Catat Sejarahnya:** Setiap kali Anda merubah sistem/alur yang sudah berjalan, catat **alasan perubahannya** di dalam file `decision-log.md` (*ADR*) di masing-masing domain.
- **Changelog:** Biasakan untuk menyisipkan tabel *Changelog* kecil di dokumen PRD untuk melacak siapa yang mengubah dan kapan.
