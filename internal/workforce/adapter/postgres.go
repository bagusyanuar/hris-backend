package adapter

import (
	"context"
	"errors"

	"github.com/bagusyanuar/hris-backend/internal/workforce/adapter/models"
	"github.com/bagusyanuar/hris-backend/internal/workforce/domain"
	"github.com/bagusyanuar/hris-backend/pkg/pagination"
	"gorm.io/gorm"
)

// departmentSortMap, jobTitleSortMap, jobPositionSortMap = whitelist logical sort key
// -> kolom DB asli (pagination.md: jangan pernah teruskan Request.Sort mentah ke GORM Order()).
var departmentSortMap = pagination.SortMap{
	"code":       "code",
	"name":       "name",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var jobTitleSortMap = pagination.SortMap{
	"code":        "code",
	"name":        "name",
	"grade_level": "grade_level",
	"created_at":  "created_at",
	"updated_at":  "updated_at",
}

var jobPositionSortMap = pagination.SortMap{
	"name":            "name",
	"headcount_quota": "headcount_quota",
	"created_at":      "created_at",
	"updated_at":      "updated_at",
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

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{db: db}
}

// Create = INSERT baru => WAJIB Create(), BUKAN Save() (persistence-convention.md §1).
func (r *departmentRepository) Create(ctx context.Context, d *domain.Department) error {
	model, err := models.DepartmentFromDomain(d)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDepartmentCodeDuplicate
		}
		return err
	}
	return nil
}

func (r *departmentRepository) FindByID(ctx context.Context, id string) (*domain.Department, error) {
	var model models.DepartmentModel
	if err := dbFromContext(ctx, r.db).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *departmentRepository) FindByCompanyAndCode(ctx context.Context, companyID, code string) (*domain.Department, error) {
	var model models.DepartmentModel
	if err := dbFromContext(ctx, r.db).
		First(&model, "company_id = ? AND code = ?", companyID, code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindParentID dipakai domain.DetectCycle — SELECT kolom parent_id saja, bukan seluruh row.
func (r *departmentRepository) FindParentID(ctx context.Context, id string) (*string, error) {
	var model models.DepartmentModel
	if err := dbFromContext(ctx, r.db).Select("id", "parent_id").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}
		return nil, err
	}
	return model.ToDomain().ParentID, nil
}

func (r *departmentRepository) FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*domain.Department, int64, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order, Search: search}
	db := dbFromContext(ctx, r.db).Order(req.OrderClause(departmentSortMap, "created_at"))
	if clause, args := req.SearchClause("code", "name"); clause != "" {
		db = db.Where(clause, args...)
	}
	rows, meta, err := pagination.Query[models.DepartmentModel](db, req)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*domain.Department, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, meta.Total, nil
}

// FindAllTree — TANPA pagination, dipakai GET /departments/tree (FE Tabel nested row + Bagan).
func (r *departmentRepository) FindAllTree(ctx context.Context) ([]*domain.Department, error) {
	var rows []models.DepartmentModel
	if err := dbFromContext(ctx, r.db).
		Where("is_active = ?", true).
		Order("created_at").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Department, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, nil
}

// FindNamesByIDs — batch SELECT id, name, dipakai embed department_name di JobPositionResponse
// (hindari N+1, lihat decision-log.md ADR-006).
func (r *departmentRepository) FindNamesByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}
	var rows []models.DepartmentModel
	if err := dbFromContext(ctx, r.db).Select("id", "name").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ID.String()] = row.Name
	}
	return result, nil
}

// Update = record PK sudah pasti ada => Save() aman untuk semantik update.
func (r *departmentRepository) Update(ctx context.Context, d *domain.Department) error {
	model, err := models.DepartmentFromDomain(d)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Save(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDepartmentCodeDuplicate
		}
		return err
	}
	return nil
}

func (r *departmentRepository) Delete(ctx context.Context, id string) error {
	return dbFromContext(ctx, r.db).Delete(&models.DepartmentModel{}, "id = ?", id).Error
}

type jobTitleRepository struct {
	db *gorm.DB
}

func NewJobTitleRepository(db *gorm.DB) domain.JobTitleRepository {
	return &jobTitleRepository{db: db}
}

// Create = INSERT baru => WAJIB Create(), BUKAN Save() (persistence-convention.md §1).
func (r *jobTitleRepository) Create(ctx context.Context, jt *domain.JobTitle) error {
	model, err := models.JobTitleFromDomain(jt)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrJobTitleCodeDuplicate
		}
		return err
	}
	return nil
}

func (r *jobTitleRepository) FindByID(ctx context.Context, id string) (*domain.JobTitle, error) {
	var model models.JobTitleModel
	if err := dbFromContext(ctx, r.db).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrJobTitleNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *jobTitleRepository) FindByCompanyAndCode(ctx context.Context, companyID, code string) (*domain.JobTitle, error) {
	var model models.JobTitleModel
	if err := dbFromContext(ctx, r.db).
		First(&model, "company_id = ? AND code = ?", companyID, code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrJobTitleNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *jobTitleRepository) FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*domain.JobTitle, int64, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order, Search: search}
	db := dbFromContext(ctx, r.db).Order(req.OrderClause(jobTitleSortMap, "created_at"))
	if clause, args := req.SearchClause("code", "name"); clause != "" {
		db = db.Where(clause, args...)
	}
	rows, meta, err := pagination.Query[models.JobTitleModel](db, req)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*domain.JobTitle, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, meta.Total, nil
}

// FindNamesByIDs — batch SELECT id, name, dipakai embed job_title_name di JobPositionResponse
// (hindari N+1, lihat decision-log.md ADR-006).
func (r *jobTitleRepository) FindNamesByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}
	var rows []models.JobTitleModel
	if err := dbFromContext(ctx, r.db).Select("id", "name").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ID.String()] = row.Name
	}
	return result, nil
}

// Update = record PK sudah pasti ada => Save() aman untuk semantik update.
func (r *jobTitleRepository) Update(ctx context.Context, jt *domain.JobTitle) error {
	model, err := models.JobTitleFromDomain(jt)
	if err != nil {
		return err
	}
	if err := dbFromContext(ctx, r.db).Save(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrJobTitleCodeDuplicate
		}
		return err
	}
	return nil
}

func (r *jobTitleRepository) Delete(ctx context.Context, id string) error {
	return dbFromContext(ctx, r.db).Delete(&models.JobTitleModel{}, "id = ?", id).Error
}

type jobPositionRepository struct {
	db *gorm.DB
}

func NewJobPositionRepository(db *gorm.DB) domain.JobPositionRepository {
	return &jobPositionRepository{db: db}
}

// Create = INSERT baru => WAJIB Create(), BUKAN Save() (persistence-convention.md §1).
func (r *jobPositionRepository) Create(ctx context.Context, jp *domain.JobPosition) error {
	model, err := models.JobPositionFromDomain(jp)
	if err != nil {
		return err
	}
	return dbFromContext(ctx, r.db).Create(model).Error
}

func (r *jobPositionRepository) FindByID(ctx context.Context, id string) (*domain.JobPosition, error) {
	var model models.JobPositionModel
	if err := dbFromContext(ctx, r.db).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrJobPositionNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindParentID dipakai domain.DetectCycle — SELECT kolom reports_to_id saja, bukan seluruh row.
func (r *jobPositionRepository) FindParentID(ctx context.Context, id string) (*string, error) {
	var model models.JobPositionModel
	if err := dbFromContext(ctx, r.db).Select("id", "reports_to_id").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrJobPositionNotFound
		}
		return nil, err
	}
	return model.ToDomain().ReportsToID, nil
}

func (r *jobPositionRepository) FindAll(ctx context.Context, page, limit int, sort, order, search string) ([]*domain.JobPosition, int64, error) {
	req := pagination.Request{Page: page, Limit: limit, Sort: sort, Order: order, Search: search}
	db := dbFromContext(ctx, r.db).Order(req.OrderClause(jobPositionSortMap, "created_at"))
	if clause, args := req.SearchClause("name"); clause != "" {
		db = db.Where(clause, args...)
	}
	rows, meta, err := pagination.Query[models.JobPositionModel](db, req)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*domain.JobPosition, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, meta.Total, nil
}

// FindAllChart — TANPA pagination (decision-log.md ADR-004). Scope-aware via
// scope.FromContext(ctx) begitu RBAC landing (staged) — untuk sekarang tanpa filter
// (owner-mode), dataset per-company diasumsikan kecil.
func (r *jobPositionRepository) FindAllChart(ctx context.Context) ([]*domain.JobPosition, error) {
	var rows []models.JobPositionModel
	if err := dbFromContext(ctx, r.db).
		Where("is_active = ?", true).
		Order("created_at").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.JobPosition, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToDomain())
	}
	return result, nil
}

// Update = record PK sudah pasti ada => Save() aman untuk semantik update.
func (r *jobPositionRepository) Update(ctx context.Context, jp *domain.JobPosition) error {
	model, err := models.JobPositionFromDomain(jp)
	if err != nil {
		return err
	}
	return dbFromContext(ctx, r.db).Save(model).Error
}

func (r *jobPositionRepository) Delete(ctx context.Context, id string) error {
	return dbFromContext(ctx, r.db).Delete(&models.JobPositionModel{}, "id = ?", id).Error
}
