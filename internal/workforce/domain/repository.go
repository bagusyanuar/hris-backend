package domain

import "context"

// DepartmentRepository — FindAll WAJIB filter company_id dari scope.FromContext(ctx)
// begitu RBAC landing (scoping-convention.md §3). Untuk sekarang scope selalu kosong
// (owner-mode, tanpa filter tambahan).
type DepartmentRepository interface {
	Create(ctx context.Context, d *Department) error
	FindByID(ctx context.Context, id string) (*Department, error) // not-found => ErrDepartmentNotFound
	// FindByCompanyAndCode dipakai cek duplikasi code dalam satu company sebelum insert.
	// Mengembalikan ErrDepartmentNotFound (bukan nil,nil) kalau kosong.
	FindByCompanyAndCode(ctx context.Context, companyID, code string) (*Department, error)
	// FindParentID mengembalikan ParentID milik department id — dipakai DetectCycle.
	FindParentID(ctx context.Context, id string) (*string, error)
	// sort/order mentah dari client — whitelist kolom sortable dilakukan di adapter (pagination.SortMap).
	// search kosong = tanpa filter; non-kosong = match code ATAU name (ILIKE substring, case-insensitive).
	FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*Department, int64, error)
	// FindAllTree — TANPA pagination, dipakai GET /departments/tree buat FE Tabel (nested row,
	// expand/collapse) & Bagan (tree diagram). Mengembalikan SEMUA department aktif dalam scope;
	// dataset per-company diasumsikan kecil.
	FindAllTree(ctx context.Context) ([]*Department, error)
	// FindNamesByIDs — batch id->name, dipakai embed department_name di JobPositionResponse
	// (hindari N+1). ids kosong = map kosong, bukan error.
	FindNamesByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Update(ctx context.Context, d *Department) error
	Delete(ctx context.Context, id string) error
}

// JobTitleRepository — FindAll WAJIB filter company_id dari scope.FromContext(ctx)
// begitu RBAC landing. Untuk sekarang scope selalu kosong (owner-mode).
type JobTitleRepository interface {
	Create(ctx context.Context, jt *JobTitle) error
	FindByID(ctx context.Context, id string) (*JobTitle, error) // not-found => ErrJobTitleNotFound
	// FindByCompanyAndCode dipakai cek duplikasi code dalam satu company sebelum insert.
	FindByCompanyAndCode(ctx context.Context, companyID, code string) (*JobTitle, error)
	// search kosong = tanpa filter; non-kosong = match code ATAU name (ILIKE substring, case-insensitive).
	FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*JobTitle, int64, error)
	// FindNamesByIDs — batch id->name, dipakai embed job_title_name di JobPositionResponse
	// (hindari N+1). ids kosong = map kosong, bukan error.
	FindNamesByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Update(ctx context.Context, jt *JobTitle) error
	Delete(ctx context.Context, id string) error
}

// JobPositionRepository — FindAll/FindAllChart WAJIB filter company_id dari
// scope.FromContext(ctx) begitu RBAC landing. Untuk sekarang scope selalu kosong
// (owner-mode, tanpa filter tambahan).
type JobPositionRepository interface {
	Create(ctx context.Context, jp *JobPosition) error
	FindByID(ctx context.Context, id string) (*JobPosition, error) // not-found => ErrJobPositionNotFound
	// FindParentID mengembalikan ReportsToID milik job position id — dipakai DetectCycle.
	FindParentID(ctx context.Context, id string) (*string, error)
	// sort/order mentah dari client — whitelist kolom sortable dilakukan di adapter (pagination.SortMap).
	// search kosong = tanpa filter; non-kosong = match name (ILIKE substring, case-insensitive) — Job Position gak punya kolom code.
	FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*JobPosition, int64, error)
	// FindAllChart — TANPA pagination, dipakai GET /job-positions/chart (decision-log.md ADR-004).
	// Mengembalikan SEMUA job position aktif dalam scope; dataset per-company diasumsikan kecil.
	FindAllChart(ctx context.Context) ([]*JobPosition, error)
	Update(ctx context.Context, jp *JobPosition) error
	Delete(ctx context.Context, id string) error
}
