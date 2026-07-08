package database

import (
	"fmt"
	"log"

	"github.com/bagusyanuar/hris-backend/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB menginisialisasi koneksi PostgreSQL menggunakan GORM
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=%s",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbSslMode, cfg.DbTz,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Printf("Connected to database %s on %s:%s successfully", cfg.DbName, cfg.DbHost, cfg.DbPort)
	return db, nil
}
