package domain

import "context"

// CompanyRepository — semua FindXxx/FindAll scope-aware sesuai kontrak scope.FromContext,
// tapi Company sendiri adalah legal root (dia scope-nya sendiri), bukan dikelilingi filter
// company_id (scoping-convention.md §1).
type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	FindByID(ctx context.Context, id string) (*Company, error) // not-found => ErrCompanyNotFound
	// sort/order mentah dari client — whitelist kolom sortable dilakukan di adapter (pagination.SortMap).
	FindAll(ctx context.Context, page, limit int, sort, order string) ([]*Company, int64, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id string) error
}

// BranchRepository — FindAllByCompany WAJIB filter company_id dari scope.FromContext(ctx)
// begitu RBAC landing (scoping-convention.md §3). Untuk sekarang scope selalu kosong (owner-mode).
type BranchRepository interface {
	Create(ctx context.Context, branch *Branch) error
	FindByID(ctx context.Context, id string) (*Branch, error) // not-found => ErrBranchNotFound
	// FindByCompanyAndCode mengembalikan ErrBranchNotFound (bukan nil,nil) kalau kosong,
	// dipakai untuk cek duplikasi code dalam satu company sebelum insert.
	FindByCompanyAndCode(ctx context.Context, companyID, code string) (*Branch, error)
	// sort/order mentah dari client — whitelist kolom sortable dilakukan di adapter (pagination.SortMap).
	FindAllByCompany(ctx context.Context, companyID string, page, limit int, sort, order string) ([]*Branch, int64, error)
	// DemoteMainBranch set is_main=false untuk main branch lama di company yang sama
	// (decision-log.md ADR-004) — dipanggil di dalam TxManager.Do sebelum insert/update main baru.
	DemoteMainBranch(ctx context.Context, companyID string) error
	Update(ctx context.Context, branch *Branch) error
	Delete(ctx context.Context, id string) error
}
