package database

import (
	"fmt"
	"log"

	"github.com/bagusyanuar/hris-backend/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB menginisialisasi koneksi PostgreSQL menggunakan GORM
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort, cfg.DbSslMode, cfg.DbTz,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Printf("Connected to database %s on %s:%s successfully", cfg.DbName, cfg.DbHost, cfg.DbPort)
	return db, nil
}
