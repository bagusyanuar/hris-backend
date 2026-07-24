package pagination

import "strings"

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// Request adalah parameter halaman + sort + search mentah dari client (query
// param page/limit/sort/order/search). Sort adalah logical key dari client
// (mis. "name"), BUKAN nama kolom DB asli — resolusinya lewat OrderClause +
// SortMap supaya gak ada raw client string yang nyampe ke GORM Order() (SQL
// injection). Search sama prinsipnya lewat SearchClause + whitelist kolom.
type Request struct {
	Page   int
	Limit  int
	Sort   string
	Order  string // "asc" | "desc"
	Search string // opsional; kosong = tanpa filter search
}

// Normalize mengisi default (page=1, limit=20, order=asc) dan meng-cap limit
// maksimum, supaya caller (application/handler) gak perlu ulang logic ini
// tiap domain. Sort TIDAK divalidasi di sini (whitelist-nya beda tiap entity)
// — itu tanggung jawab OrderClause di adapter layer.
func (r Request) Normalize() Request {
	if r.Page < 1 {
		r.Page = DefaultPage
	}
	if r.Limit < 1 || r.Limit > MaxLimit {
		r.Limit = DefaultLimit
	}
	if !strings.EqualFold(r.Order, "desc") {
		r.Order = "asc"
	} else {
		r.Order = "desc"
	}
	return r
}

func (r Request) Offset() int {
	return (r.Page - 1) * r.Limit
}

// SortMap memetakan logical sort key yang boleh dipakai client ke nama kolom
// DB asli. WAJIB whitelist eksplisit per entity — jangan pernah teruskan
// Request.Sort mentah ke GORM Order().
type SortMap map[string]string

// OrderClause mengembalikan klausa "kolom ASC/DESC" yang aman buat GORM
// .Order(...). Kalau r.Sort gak ada di whitelist allowed, fallback ke
// defaultSort (juga harus ada di allowed).
func (r Request) OrderClause(allowed SortMap, defaultSort string) string {
	col, ok := allowed[r.Sort]
	if !ok {
		col = allowed[defaultSort]
	}
	dir := "ASC"
	if strings.EqualFold(r.Order, "desc") {
		dir = "DESC"
	}
	return col + " " + dir
}

// SearchClause mengembalikan klausa "(col1 ILIKE ? OR col2 ILIKE ? ...)" aman
// buat GORM .Where(clause, args...) beserta args-nya (satu "%search%" per
// kolom), atau ("", nil) kalau Search kosong atau columns kosong. columns
// WAJIB whitelist eksplisit oleh caller (adapter) — alasan sama seperti
// SortMap: jangan pernah teruskan nama kolom mentah dari client (SQL
// injection lewat Order()/Where() string mentah, pagination-convention.md §3).
//
// Dipakai buat kasus search same-table sederhana (mis. WHERE code ILIKE ?
// OR name ILIKE ?). Kasus search lintas-table (mis. Company match nama
// Branch anaknya via EXISTS subquery — lihat organization tech-spec.md §6.1)
// TETAP boleh ditulis manual di adapter, di luar helper ini — bukan
// pelanggaran, itu pengecualian terdokumentasi (pagination-convention.md §3).
func (r Request) SearchClause(columns ...string) (string, []any) {
	if r.Search == "" || len(columns) == 0 {
		return "", nil
	}
	like := "%" + r.Search + "%"
	parts := make([]string, len(columns))
	args := make([]any, len(columns))
	for i, c := range columns {
		parts[i] = c + " ILIKE ?"
		args[i] = like
	}
	return "(" + strings.Join(parts, " OR ") + ")", args
}

// Meta adalah metadata halaman yang dikembalikan ke client bersama Items.
type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func NewMeta(req Request, total int64) Meta {
	return Meta{Page: req.Page, Limit: req.Limit, Total: total}
}
