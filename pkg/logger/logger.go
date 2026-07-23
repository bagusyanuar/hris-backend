// Package logger menyediakan wrapper zap.Logger sebagai satu-satunya
// sumber logging terstruktur di seluruh layer (domain, application, adapter, transport).
package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var (
	base *zap.Logger
	once sync.Once
)

// Init menginisialisasi global logger sesuai environment aplikasi.
// env "production" -> JSON encoder, level Info. Selain itu -> console encoder, level Debug.
func Init(env string, debug bool) *zap.Logger {
	once.Do(func() {
		var cfg zap.Config
		if env == "production" {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		if debug {
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}

		l, err := cfg.Build()
		if err != nil {
			// Fallback: logger tidak boleh bikin aplikasi gagal start.
			l = zap.NewNop()
		}
		base = l
	})
	return base
}

// L mengembalikan global logger. Jika Init belum dipanggil (mis. di unit test),
// kembalikan Nop logger supaya caller tidak nil-panic.
func L() *zap.Logger {
	if base == nil {
		return zap.NewNop()
	}
	return base
}

// Sync melakukan flush buffer log. Panggil sebelum aplikasi exit (graceful shutdown).
func Sync() error {
	if base == nil {
		return nil
	}
	return base.Sync()
}

// WithContext menyisipkan logger (biasanya sudah dibubuhi field request_id)
// ke dalam context, supaya application/adapter layer bisa memakai logger
// yang sama tanpa perlu inject manual di tiap fungsi.
func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext mengambil logger dari context. Kalau tidak ada, fallback ke global logger.
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok && l != nil {
		return l
	}
	return L()
}
