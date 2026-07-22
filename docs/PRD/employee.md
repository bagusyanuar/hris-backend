# Product Requirements: Employee Module

## 1. Tujuan & Dampak (The "Why")
Mengelola data induk (Core Data) karyawan secara terpusat untuk menggantikan pencatatan manual/Excel. Tujuannya adalah untuk mendigitalisasi arsip kepegawaian sehingga pembaruan data menjadi *real-time*, mempercepat pencarian profil karyawan dari hitungan menit menjadi hitungan detik, serta memberikan sumber data tunggal (Single Source of Truth) yang solid bagi modul operasional lain seperti Absensi dan Payroll.

## 2. Scope & Out-of-Scope (Batasan Tegas)

**In-Scope (Dikerjakan):**
- Pencatatan Biodata, Alamat, Kontak, dan Kontak Darurat Karyawan.
- Penempatan *Job Position* karyawan (relasi ke modul Organization).
- Penyimpanan data Rekening Bank karyawan untuk keperluan transfer gaji.
- Pencatatan riwayat pendidikan (Education) dan pengalaman kerja (Experience).
- Manajemen tautan dokumen digital (URL upload Ijazah, KTP, KK, dll).
- Status kepegawaian (Contoh: Active, Inactive, Resign).

**Out-of-Scope (TIDAK Dikerjakan di modul ini):**
- Perhitungan Gaji, Tunjangan, PPh 21, atau BPJS (Ini masuk ranah modul Payroll).
- Pencatatan *Clock-in/out* Absensi atau pengajuan Cuti (Masuk modul Attendance/Leave).
- Registrasi akun (Sign up), reset password, atau manajemen hak akses level sistem (Masuk modul Auth).
- Fitur Rekrutmen (Applicant Tracking System / ATS).

## 3. User Roles & Permissions
- **Superadmin / HR Manager:** Bisa melihat (Read), membuat (Create), mengubah (Update), dan me-nonaktifkan (Offboard) semua data karyawan di perusahaan.
- **Karyawan (Employee):** Hanya bisa melihat (Read-only) profil dirinya sendiri melalui portal ESS (*Employee Self Service*). 
- *Catatan:* Untuk MVP tahap 1, karyawan tidak diperbolehkan *update* data sendiri secara langsung. Semua perubahan data harus melalui permintaan manual ke pihak HR. Sistem *Self-Service Approval Workflow* belum masuk di tahap ini.

## 4. Kriteria Penerimaan (Acceptance Criteria)

**Skenario 1: Validasi Duplikasi KTP**
- **Given** HR Admin sedang mengisi form tambah karyawan baru.
- **When** Admin memasukkan NIK KTP (`ktp_number`) yang sudah ada di *database*.
- **Then** Sistem menolak proses penyimpanan dan menampilkan error *"NIK KTP sudah terdaftar pada karyawan lain"*.

**Skenario 2: Offboarding (Resign / PHK)**
- **Given** Profil seorang karyawan aktif.
- **When** HR Admin menekan tombol "Nonaktifkan Karyawan" atau memproses Resign.
- **Then** Sistem mengubah status karyawan menjadi `INACTIVE`, mengisi data `resign_date`, memblokir hak akses login-nya di modul Auth, **NAMUN** tidak menghapus data karyawan secara fisik dari database.

**Skenario 3: Aturan Rekening Bank**
- **Given** HR Admin sedang melengkapi data Bank Karyawan.
- **When** Admin menyimpan data.
- **Then** Sistem memastikan harus ada tepat 1 rekening bank yang ditandai sebagai Rekening Utama (`is_primary = true`) per karyawan.

## 5. Technical & Architectural Constraints
- **Domain-Driven Design (DDD):** Kode harus terisolasi di domainnya sendiri (`internal/domain/employee`). Modul ini dilarang mengambil data dari DB langsung (*bypass*) ke tabel modul lain, harus melalui interface/API internal.
- **Data Deletion:** Dilarang menggunakan aksi `HARD DELETE` (`DELETE FROM ...`) untuk tabel karyawan demi integritas sejarah penggajian. Wajib menggunakan `SOFT DELETE` (mengisi `deleted_at`).
- **UI Constraints (Frontend):** Form penambahan karyawan (Create) sangat panjang, sehingga **wajib** dirender dengan gaya **Wizard / Multi-step Form**. Setiap perpindahan step harus melakukan penyimpanan progresif (Progressive Save) ke Backend melalui endpoint yang berbeda (misal: Step 1 Simpan Core Data, Step 2 Simpan Personal, dst) untuk mencegah hilangnya input jika terjadi kegagalan sistem atau browser ter-refresh.

## 6. Dependencies (Ketergantungan)
- **Modul Auth:** Pembuatan Karyawan Baru berstatus **Blokir** (Dependent) pada keberhasilan pembuatan akun di tabel `users`. Karyawan wajib memiliki `user_id`.
- **Modul Organization:** Memerlukan *endpoint* (misal: `GET /organization/job-positions`) agar Frontend bisa me-*render* dropdown pemilihan Jabatan saat menambahkan Karyawan.

---

## 7. Data Schema & Business Rules (Database Map)

Berikut adalah breakdown logika tabel untuk mempermudah perancangan Entity (BE) dan Mocking/UI (FE).

### 7.1. Employee (Data Utama Pekerjaan)
Menyimpan posisi karyawan di perusahaan saat ini.
- **Aturan Bisnis:** Harus terhubung ke `Job Position`. Kolom `join_date` mutlak wajib diisi untuk acuan THR dan senioritas.

| id | user_id | employee_code (NIK) | job_position_id | employment_status | join_date | status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `emp-1` | `usr-1` | HR-001 | `pos-1` (CEO) | PERMANENT | 2020-01-01 | ACTIVE |

### 7.2. Personal Data (Biodata Diri)
Menyimpan data sipil karyawan.
- **Aturan Bisnis:** `ptkp_status` wajib diisi sebagai variabel penting untuk kalkulator PPh21 di modul masa depan.

| employee_id | full_name | ktp_number | gender | marital_status | ptkp_status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `emp-1` | Budi Santoso | 317100000001 | MALE | MARRIED | K/1 |

### 7.3. Contact & Bank Account
- **Aturan Bisnis:** `is_primary` digunakan oleh sistem Payroll untuk mendeteksi kemana gaji bulan tersebut harus ditransfer.

| employee_id | phone_number | bank_name | account_number | account_holder_name | is_primary |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `emp-1` | 08123456789 | BCA | 1234567890 | Budi Santoso | true |

### 7.4. History & Documents (Rekam Jejak & Arsip)
- **Aturan Bisnis:** Bersifat One-to-Many. Bisa ditambahkan menyusul (opsional di awal pendaftaran).

**Tabel `employee_educations` (Contoh):**
| level | institution_name | major | start_year | end_year | score |
| :--- | :--- | :--- | :--- | :--- | :--- |
| S1 | Universitas Indonesia | Sistem Informasi | 2010 | 2014 | 3.50 |

**Tabel `employee_documents` (Contoh):**
| document_type | document_name | document_url |
| :--- | :--- | :--- |
| KTP | scan_ktp.pdf | `https://s3.../ktp_budi.pdf` |
