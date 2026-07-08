# Product Requirements: Organization Module

Dokumen ini menjelaskan fungsionalitas, konsep, dan contoh data untuk modul **Organization** di dalam sistem HRIS. Tujuannya adalah untuk menyamakan pemahaman antara tim Backend dan Frontend agar tidak ada kebingungan saat membangun UI/UX.

## Konsep Dasar (The 3-Pillars)

Sistem organisasi perusahaan di dalam HRIS ini dibangun di atas 3 pilar utama:
1. **Department**: Struktur unit kerja perusahaan. Bisa berupa divisi utama, departemen, atau sub-departemen. Memiliki relasi *hierarki* (Parent-Child).
2. **Job Title**: Master data pangkat, golongan, atau jenjang karir baku yang berlaku secara umum di perusahaan.
3. **Job Position**: Slot jabatan aktual / definitif di dalam struktur organisasi yang merupakan kombinasi antara **Department** dan **Job Title**. Jabatan ini mendefinisikan "Siapa lapor ke siapa" (*reporting line*) serta jumlah batasan kuota pegawai (*headcount quota*).

### Mengapa Menggunakan Konsep 3 Pilar? (Position-Based vs Person-Based)
Di banyak sistem HR sederhana (UMKM), biasanya karyawan langsung ditempelkan nama jabatannya (*Person-Based*). Namun, untuk sistem *Enterprise-Grade*, kita menggunakan pendekatan *Position-Based* (3 Pilar). Artinya, struktur organisasi dan "kursi" jabatannya dibentuk terlebih dahulu, baru kemudian karyawan menduduki kursi tersebut.

**Pros (Kelebihan):**
- **Struktur Independen:** Jika seorang manajer *resign*, struktur pelaporan (*reporting line*) di bawahnya tidak rusak karena bawahan melapor ke *Posisi* manajer, bukan ke *Orang*-nya.
- **Manajemen Kuota & Budget (Headcount):** Memudahkan finance dan HR untuk membatasi jumlah pegawai (contoh: Posisi "Staf IT" hanya boleh diisi maksimal 5 orang).
- **Standarisasi Gaji (Grade):** Pemisahan `Job Title` (Pangkat) memastikan standarisasi gaji/fasilitas yang adil di lintas departemen (contoh: Manajer IT dan Manajer HR berada di *grade* yang sama).

**Cons (Kekurangan):**
- **Kompleksitas Awal (Setup):** Butuh waktu ekstra di awal untuk mengatur master data. Admin HR tidak bisa langsung menambah karyawan, mereka harus membuat Departemen, lalu Job Title, dan menyatukannya jadi Job Position terlebih dahulu.
- **Kurang Cocok untuk Start-Up Kecil:** Perusahaan dengan struktur yang sangat cair (pegawai merangkap banyak peran abstrak) mungkin merasa sistem ini terlalu kaku.

---

## 1. Department (Unit Kerja)
Menyimpan struktur divisi atau departemen. Relasi bersifat *Tree* atau hierarkis menggunakan `parent_id`.

### Aturan Bisnis:
- Jika `parent_id` adalah `null`, berarti ini adalah departemen level tertinggi (Root).
- Di frontend, tampilan ini biasanya direpresentasikan sebagai **Tree View** atau nested list.

### Sample Data (Tabel `departments`):

| id | code | name | parent_id | is_active |
| :--- | :--- | :--- | :--- | :--- |
| `dept-1` | DIR | Direksi | `null` | true |
| `dept-2` | TI | Divisi Teknologi Informasi | `dept-1` | true |
| `dept-3` | DEV | Departemen Pengembangan (Engineering) | `dept-2` | true |
| `dept-4` | OPR | Divisi Operasional | `dept-1` | true |
| `dept-5` | SDM | Divisi Sumber Daya Manusia (HR) | `dept-1` | true |

---

## 2. Job Title (Pangkat / Grade)
Master data standarisasi jabatan atau jenjang karir. Biasanya digunakan untuk menentukan standar gaji (*Salary Band*) atau fasilitas (*Benefit*).

### Aturan Bisnis:
- `grade_level` menentukan tinggi/rendahnya pangkat secara angka (misalnya makin tinggi angkanya, makin tinggi pangkatnya).
- Tidak terkait dengan departemen tertentu (independen).

### Sample Data (Tabel `job_titles`):

| id | code | name | grade_level | is_active |
| :--- | :--- | :--- | :--- | :--- |
| `title-1` | DIR | Direktur | 10 | true |
| `title-2` | KDV | Kepala Divisi / GM | 9 | true |
| `title-3` | MGR | Manajer | 7 | true |
| `title-4` | SPV | Supervisor | 5 | true |
| `title-5` | STF | Staf | 3 | true |

---

## 3. Job Position (Jabatan Aktif / Posisi)
Ini adalah "kursi" aktual yang diduduki oleh pegawai di dalam struktur organisasi.

### Aturan Bisnis:
- **Kombinasi**: Setiap Job Position harus menempel pada satu `Department` dan satu `Job Title`.
- **Reporting Line**: `reports_to_id` menunjuk ke ID Job Position lain sebagai atasannya, membentuk **Organization Chart** (Bagan Struktur Organisasi).
- **Headcount Quota**: Menentukan batas maksimal pegawai yang bisa menduduki jabatan ini (misalnya, CEO kuotanya 1, tapi Software Engineer kuotanya bisa 10).

### Sample Data (Tabel `job_positions`):

| id | name (Posisi) | department_id | job_title_id | reports_to_id (Atasan) | headcount_quota |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `pos-1` | Direktur Utama | `dept-1` (DIR) | `title-1` (DIR) | `null` | 1 |
| `pos-2` | Direktur Teknologi (CTO) | `dept-2` (TI) | `title-1` (DIR) | `pos-1` (Dirut) | 1 |
| `pos-3` | Kepala Divisi TI | `dept-2` (TI) | `title-2` (KDV) | `pos-2` (Dirtek) | 1 |
| `pos-4` | Manajer Pengembangan | `dept-3` (DEV) | `title-3` (MGR) | `pos-3` (Kadiv TI) | 3 |
| `pos-5` | Supervisor Backend | `dept-3` (DEV) | `title-4` (SPV) | `pos-4` (Mgr DEV) | 5 |
| `pos-6` | Staf Programmer Backend | `dept-3` (DEV) | `title-5` (STF) | `pos-5` (Spv BE) | 10 |

---

## Catatan Khusus untuk Frontend (FE)

1. **Pembuatan Form `Job Position`**:
   - Di form "Create Job Position", FE membutuhkan dropdown untuk memilih `Department` dan `Job Title`. Oleh karena itu, FE harus memanggil API `GET /organization/departments` dan `GET /organization/job-titles` terlebih dahulu untuk mengisi *dropdown option*.
   - Input `Reports To` adalah *autocomplete dropdown* yang mengambil data dari `GET /organization/job-positions`.
   
2. **Organization Chart (Bagan Organisasi)**:
   - Data `job_positions` yang saling terkait lewat `reports_to_id` bisa dirender menjadi **Organization Chart** visual.
   - Root (puncak) dari chart adalah posisi dengan `reports_to_id` bernilai `null` (seperti contoh `pos-1` CEO di atas).

3. **Status `is_active`**:
   - Secara *default*, data yang dikembalikan oleh API adalah yang aktif (jika tidak difilter). Jika suatu departemen/posisi di-nonaktifkan, UI dapat menampilkannya dengan warna abu-abu (greyed out) atau disembunyikan.
