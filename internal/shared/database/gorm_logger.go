package database

import (
	"context"
	"errors"
	"time"

	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// zapGormLogger meneruskan log query GORM ke pkg/logger (zap), menggantikan
// logger bawaan GORM yang print polos ke stdout.
type zapGormLogger struct {
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
}

// newGormLogger membangun GORM logger.Interface: level Info (log semua query)
// di non-production, level Warn (hanya slow-query + error) di production.
func newGormLogger(env string, slowThresholdMs int) gormlogger.Interface {
	level := gormlogger.Warn
	if env != "production" {
		level = gormlogger.Info
	}
	return &zapGormLogger{
		logLevel:      level,
		slowThreshold: time.Duration(slowThresholdMs) * time.Millisecond,
	}
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.logLevel = level
	return &clone
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		logger.FromContext(ctx).Sugar().Infof(msg, data...)
	}
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		logger.FromContext(ctx).Sugar().Warnf(msg, data...)
	}
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		logger.FromContext(ctx).Sugar().Errorf(msg, data...)
	}
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	log := logger.FromContext(ctx)
	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	switch {
	case err != nil && l.logLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		log.Error("gorm query error", append(fields, zap.Error(err))...)
	case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.logLevel >= gormlogger.Warn:
		log.Warn("gorm slow query", fields...)
	case l.logLevel >= gormlogger.Info:
		log.Debug("gorm query", fields...)
	}
}
