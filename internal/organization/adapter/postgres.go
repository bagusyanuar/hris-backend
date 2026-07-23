package adapter

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/organization/adapter/models"
	"github.com/bagusyanuar/hris-backend/internal/organization/domain"
	"github.com/bagusyanuar/hris-backend/pkg/pagination"
	"gorm.io/gorm"
)

// companySortMap & branchSortMap = whitelist logical sort key -> kolom DB asli
// (pagination.md: jangan pernah teruskan Request.Sort mentah ke GORM Order()).
var companySortMap = pagination.SortMap{
	"code":       "code",
	"legal_name": "legal_name",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var branchSortMap = pagination.SortMap{
	"code":       "code",
	"name":       "name",
	"is_main":    "is_main",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

type txKey struct{}

// GormTxManager implementasi domain.TxManager (persistence-convention.md §2).
// Application layer memanggil Do(...); repository di dalam fn membaca koneksi
// transaksi yang sama lewat dbFromContext, bukan koneksi db dasar.
type GormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) domain.TxManager {
	return &GormTxManager{db: db}
}

func (m *GormTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

// dbFromContext mengembalikan handle transaksi aktif (kalau dipanggil di dalam
// TxManager.Do) atau koneksi db dasar sebagai fallback.
func dbFromContext(ctx context.Context, base *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return base.WithContext(ctx)
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) domain.CompanyRepository {
	return &companyRepository{db: db}
}

// Create = INSERT baru => WAJIB Create(), BUKAN Save() (persistence-convention.md §1).
func (r *companyRepository) Create(ctx context.Context, company *domain.Company) error {
	model, err := models.CompanyFromDomain(company)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrCompanyNPWPDuplicate
		}
		return err
	}
	return nil
}

func (r *companyRepository) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	var model models.CompanyModel
	if err := dbFromContext(ctx, r.db).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCompanyNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *companyRepository) FindAll(ctx context.Context, page, limit int, sort, order string) ([]*domain.Company, int64, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}
	db := dbFromContext(ctx, r.db).Order(req.OrderClause(companySortMap, "created_at"))
	rows, meta, err := pagination.Query[models.CompanyModel](db, req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*domain.Company, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, meta.Total, nil
}

// Update = record PK sudah pasti ada => Save() aman untuk semantik update.
func (r *companyRepository) Update(ctx context.Context, company *domain.Company) error {
	model, err := models.CompanyFromDomain(company)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Save(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrCompanyNPWPDuplicate
		}
		return err
	}
	return nil
}

func (r *companyRepository) Delete(ctx context.Context, id string) error {
	return dbFromContext(ctx, r.db).Delete(&models.CompanyModel{}, "id = ?", id).Error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) domain.BranchRepository {
	return &branchRepository{db: db}
}

// Create = INSERT baru => WAJIB Create(), BUKAN Save() (persistence-convention.md §1).
func (r *branchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	model, err := models.BranchFromDomain(branch)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrBranchCodeDuplicate
		}
		return err
	}
	return nil
}

func (r *branchRepository) FindByID(ctx context.Context, id string) (*domain.Branch, error) {
	var model models.BranchModel
	if err := dbFromContext(ctx, r.db).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBranchNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *branchRepository) FindByCompanyAndCode(ctx context.Context, companyID, code string) (*domain.Branch, error) {
	var model models.BranchModel
	if err := dbFromContext(ctx, r.db).
		First(&model, "company_id = ? AND code = ?", companyID, code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBranchNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *branchRepository) FindAllByCompany(ctx context.Context, companyID string, page, limit int, sort, order string) ([]*domain.Branch, int64, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order}
	db := dbFromContext(ctx, r.db).Where("company_id = ?", companyID).Order(req.OrderClause(branchSortMap, "created_at"))
	rows, meta, err := pagination.Query[models.BranchModel](db, req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*domain.Branch, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, meta.Total, nil
}

// DemoteMainBranch — dipanggil di dalam TxManager.Do sebelum insert/update main baru
// (decision-log.md ADR-004). Update langsung by filter, bukan fetch-then-save.
func (r *branchRepository) DemoteMainBranch(ctx context.Context, companyID string) error {
	return dbFromContext(ctx, r.db).
		Model(&models.BranchModel{}).
		Where("company_id = ? AND is_main = ?", companyID, true).
		Update("is_main", false).Error
}

// Update = record PK sudah pasti ada => Save() aman untuk semantik update.
func (r *branchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	model, err := models.BranchFromDomain(branch)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Save(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrBranchCodeDuplicate
		}
		return err
	}
	return nil
}

func (r *branchRepository) Delete(ctx context.Context, id string) error {
	return dbFromContext(ctx, r.db).Delete(&models.BranchModel{}, "id = ?", id).Error
}
