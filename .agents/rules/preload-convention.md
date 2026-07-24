# Konvensi Preload Relasi di Response DTO

Aturan ini mengatur **cara menampilkan nama entity terkait** (foreign key) di response API, supaya FE gak kepaksa nampilin raw UUID di tabel/form atau nembak request tambahan cuma buat resolve nama. Lahir dari pola yang berulang independen di dua modul ([organization ADR-006](../../docs/technical/organization/decision-log.md#adr-006-nested-branches-di-get-companies--batch-query-manual-bukan-gorm-preload), [workforce-structure ADR-006](../../docs/technical/workforce-structure/decision-log.md#adr-006-jobpositionresponse-embed-departmentjob_title-sebagai-objek-id-name--batch-lookup-bukan-raw-id)) — cukup berulang buat dijadiin rule, bukan keputusan ad-hoc per modul.

Aturan ini bersifat **WAJIB (STRICT)**.

---

## 1. Nested Object `{id, name}`, BUKAN Flat `<relasi>_name`

Kalau response perlu nampilin label/nama dari entity yang direferensi FK (many-to-one, mis. `JobPosition` → `Department`), **WAJIB** bentuk nested object minimal `{id, name}` — **DILARANG** flatten jadi field terpisah kayak `department_name`.

**DILARANG:**
```go
type JobPositionResponse struct {
    DepartmentID   string `json:"department_id"`
    DepartmentName string `json:"department_name"` // ❌ flat, terpisah dari ID-nya
    ...
}
```

**WAJIB:**
```go
// Ref minimal {id, name} — reusable, dipakai berapa pun relasi yang butuh preload serupa.
type JobPositionRef struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type JobPositionResponse struct {
    Department JobPositionRef `json:"department"`
    JobTitle   JobPositionRef `json:"job_title"`
    ...
}
```

**Kenapa nested, bukan flat:**
- ID & nama entity yang sama secara logis satu unit — nested jaga mereka tetap sepasang, gak kececeran di dua field independen yang gampang out-of-sync kalau salah satu di-update parsial.
- Ekstensibel: kalau nanti butuh field lain dari relasi itu (mis. `code`), tinggal nambah ke struct `Ref`-nya — gak perlu nambah field flat baru per kebutuhan (`department_code`, `department_slug`, dst menumpuk di root object).
- FE dapet `id` di dalam object yang sama buat submit balik ke form edit — gak perlu nyimpen `department_id` terpisah dari `department_name` di state.

## 2. Nama Field = Nama Relasi (Singular), Bukan Suffix `_id`/`_name`

Field nested dinamain sesuai peran relasinya (`department`, `job_title`, `reports_to`), **bukan** `department_ref`/`department_info` atau semacamnya. Kalau relasinya nullable (mis. `reports_to` opsional), objeknya nullable — bukan dipaksa object kosong `{id: "", name: ""}`.

## 3. Batch Lookup — Wajib Hindari N+1

Isi nested object **DILARANG** di-resolve per-row (loop query satu-satu). Pola wajib:

1. **Repository** relasi yang di-preload nyediain method batch: `FindNamesByIDs(ctx context.Context, ids []string) (map[string]string, error)` — `SELECT id, name WHERE id IN (...)`, ids kosong = map kosong (bukan error).
2. **Application Layer** kumpulin seluruh ID relasi dari hasil query utama (List/Chart), panggil `FindNamesByIDs` **sekali per relasi** (bukan per row), lalu map hasilnya ke tiap item.
3. Untuk endpoint single-record (Get/Create/Update), boleh langsung `FindByID` biasa kalau entity relasinya emang udah kepanggil buat validasi lain di flow yang sama (gak perlu query tambahan).

**DILARANG** — repository A inject repository B lalu JOIN cross-domain di adapter layer (nyimpang dari [coding-convention.md](coding-convention.md) §4 — cross-domain HARUS lewat Application Service, bukan repository langsung, DAN nyimpang dari batas aggregate/bounded context per [architecture.md](architecture.md)). Batch lookup dilakukan di Application Layer masing-masing modul, hasilnya di-compose di situ — bukan di-JOIN di database layer lintas bounded context.

## 4. Scope: Many-to-One, Bukan Pengganti Pola Nested Array (Many)

Rule ini spesifik buat relasi **many-to-one** (satu `JobPosition` py satu `Department`). Kalau relasinya **one-to-many** (mis. `Company` embed banyak `Branch`), itu pola berbeda (list nested penuh, bukan `Ref` minimal) — tetap ikutin precedent [organization ADR-006](../../docs/technical/organization/decision-log.md#adr-006-nested-branches-di-get-companies--batch-query-manual-bukan-gorm-preload) (`FindAllByCompanyIDs` batch, group manual di Application Layer, embed full list child response, bukan `Ref` `{id,name}`).

## 5. Dokumentasi API (selaras api-documentation.md)

- Definisikan `Ref` struct sebagai `components.schemas` reusable di Swagger (mis. `JobPositionRef`), `$ref` dari field yang makenya — jangan `example` mentah tanpa `schema`.
- Field nested object masuk `required` di schema kalau relasinya non-nullable (selalu ada); nested-nya sendiri boleh `nullable: true` kalau relasinya opsional (mis. `reports_to`).
- Bruno: baris "Field yang SELALU ada" sebut eksplisit field nested-nya (`department` & `job_title`, bukan `department_id`/`department_name` terpisah).

## 6. Checklist Review

- [ ] Response yang nampilin nama entity FK pakai nested `{id, name}` (atau `Ref` struct reusable), bukan flat `<relasi>_name`.
- [ ] Field nested dinamain sesuai peran relasi (singular), bukan generic `_ref`/`_info`.
- [ ] Resolusi nama di List/Chart pakai batch lookup (`FindNamesByIDs`, satu query per relasi) — bukan N+1 per row.
- [ ] Gak ada repository lintas-domain di-JOIN langsung di adapter — batch lookup tetap lewat repository masing-masing domain, di-compose di Application Layer.
- [ ] Payload **request** (create/update) tetap pakai ID flat biasa (`department_id`) — rule ini cuma soal bentuk **response**.
- [ ] Relasi one-to-many (list child) pakai pola nested-array (organization ADR-006), bukan `Ref` `{id,name}`.
- [ ] Swagger: `Ref` struct sebagai `components.schemas` reusable, bukan `example` mentah.
