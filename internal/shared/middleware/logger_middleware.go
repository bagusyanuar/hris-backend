package middleware

import (
	"time"

	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLogger mencatat tiap HTTP request (method, path, status, latency, request_id)
// dan menyisipkan logger ber-request_id ke context supaya application/adapter layer
// bisa memakai logger.FromContext(ctx) dengan korelasi yang sama.
func RequestLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("X-Request-Id", requestID)

		reqLogger := logger.L().With(zap.String("request_id", requestID))
		c.SetContext(logger.WithContext(c.Context(), reqLogger))

		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		}

		switch {
		case status >= 500:
			reqLogger.Error("request completed", fields...)
		case status >= 400:
			reqLogger.Warn("request completed", fields...)
		default:
			reqLogger.Info("request completed", fields...)
		}

		return err
	}
}
