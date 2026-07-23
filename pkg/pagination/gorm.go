package pagination

import "gorm.io/gorm"

// Query menjalankan Count + Offset/Limit/Find generic di atas *gorm.DB yang
// sudah di-scope caller (mis. sudah ada .Where(...)). Dipakai HANYA di adapter
// layer (repository) — application/domain layer tetap cukup pakai Request/Meta
// (pure, tanpa gorm) di atas.
func Query[T any](db *gorm.DB, req Request) ([]T, Meta, error) {
	req = req.Normalize()

	var total int64
	if err := db.Session(&gorm.Session{}).Model(new(T)).Count(&total).Error; err != nil {
		return nil, Meta{}, err
	}

	var rows []T
	if err := db.Session(&gorm.Session{}).Offset(req.Offset()).Limit(req.Limit).Find(&rows).Error; err != nil {
		return nil, Meta{}, err
	}

	return rows, NewMeta(req, total), nil
}
