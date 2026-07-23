package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bagusyanuar/hris-backend/internal/shared/config"
	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Default fallback kalau env var terkait tidak di-set (nilai 0/kosong dari Viper).
const (
	defaultMaxOpenConns         = 25
	defaultMaxIdleConns         = 10
	defaultConnMaxLifetimeMin   = 30
	defaultConnMaxIdleTimeMin   = 5
	defaultConnectTimeoutSec    = 10
	defaultConnectRetryAttempts = 5
	defaultConnectRetryDelaySec = 2
	defaultSlowQueryThresholdMs = 200
)

func orDefault(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

// InitDB menginisialisasi koneksi PostgreSQL menggunakan GORM: retry dengan
// backoff saat startup, connection pool eksplisit, dan GORM logger yang
// diarahkan ke pkg/logger (zap) alih-alih print polos ke stdout.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=%s",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbSslMode, cfg.DbTz,
	)

	maxAttempts := orDefault(cfg.DbConnectRetryAttempts, defaultConnectRetryAttempts)
	retryDelay := time.Duration(orDefault(cfg.DbConnectRetryDelaySec, defaultConnectRetryDelaySec)) * time.Second
	connectTimeout := time.Duration(orDefault(cfg.DbConnectTimeoutSec, defaultConnectTimeoutSec)) * time.Second

	gormConfig := &gorm.Config{
		Logger:                 newGormLogger(cfg.AppEnv, orDefault(cfg.DbSlowQueryThresholdMs, defaultSlowQueryThresholdMs)),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	}

	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = connectOnce(dsn, gormConfig, connectTimeout)
		if err == nil {
			break
		}
		logger.L().Warn("failed to connect database, retrying",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxAttempts),
			zap.Error(err),
		)
		if attempt < maxAttempts {
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect database after %d attempts: %w", maxAttempts, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic sql.DB handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(orDefault(cfg.DbMaxOpenConns, defaultMaxOpenConns))
	sqlDB.SetMaxIdleConns(orDefault(cfg.DbMaxIdleConns, defaultMaxIdleConns))
	sqlDB.SetConnMaxLifetime(time.Duration(orDefault(cfg.DbConnMaxLifetimeMin, defaultConnMaxLifetimeMin)) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(orDefault(cfg.DbConnMaxIdleTimeMin, defaultConnMaxIdleTimeMin)) * time.Minute)

	logger.L().Info("connected to database successfully",
		zap.String("db_name", cfg.DbName),
		zap.String("host", cfg.DbHost),
		zap.String("port", cfg.DbPort),
	)
	return db, nil
}

// connectOnce membuka koneksi GORM lalu memverifikasinya dengan Ping ber-timeout,
// supaya DB yang belum listen (mis. container Postgres belum ready) terdeteksi
// sebagai kegagalan alih-alih hang tanpa batas waktu.
func connectOnce(dsn string, gormConfig *gorm.Config, timeout time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
