package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	authAdapter "github.com/bagusyanuar/hris-backend/internal/auth/adapter"
	"github.com/bagusyanuar/hris-backend/internal/di"
	"github.com/bagusyanuar/hris-backend/internal/shared/config"
	"github.com/bagusyanuar/hris-backend/internal/shared/middleware"
	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"github.com/bagusyanuar/hris-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the HTTP server and its dependencies.
type Server struct {
	app   *fiber.App
	cfg   *config.Config
	db    *gorm.DB
	sqlDB *sql.DB
}

// NewServer initializes a new Server instance.
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	// Dapatkan instance sql.DB untuk bisa diclose saat shutdown
	var sqlDB *sql.DB
	if db != nil {
		if sDB, err := db.DB(); err == nil {
			sqlDB = sDB
		}
	}

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName + " v" + cfg.AppVersion,
	})

	server := &Server{
		app:   app,
		cfg:   cfg,
		db:    db,
		sqlDB: sqlDB,
	}

	app.Use(middleware.RequestLogger())

	server.setupRoutes()
	return server
}

// setupRoutes initializes all routes and their dependencies.
func (s *Server) setupRoutes() {
	// Root Info Route
	s.app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "success",
			"message":     fmt.Sprintf("%s is running", s.cfg.AppName),
			"version":     s.cfg.AppVersion,
			"environment": s.cfg.AppEnv,
		})
	})

	// Health Check Route — dipakai k8s liveness/readiness probe
	s.app.Get("/health", s.healthCheck)

	// Initialize Shared Dependencies
	tokenGenerator := authAdapter.NewJWTService(s.cfg.JwtSecret, s.cfg.JwtExpiryHour, s.cfg.JwtRefreshExpiryHour)

	// Setup API Routes
	api := s.app.Group("/api/v1")

	// Initialize API Handlers via Google Wire
	handlers, err := di.InitializeAPI(s.db, tokenGenerator)
	if err != nil {
		panic("failed to initialize DI: " + err.Error())
	}

	// Register all routes
	handlers.RegisterRoutes(api)

	// Protected Example Route
	api.Get("/users/me", middleware.AuthProtected(tokenGenerator), func(c fiber.Ctx) error {
		userID := c.Locals("userID")
		role := c.Locals("role")
		return c.JSON(fiber.Map{
			"user_id": userID,
			"role":    role,
		})
	})
}

// healthCheck memverifikasi koneksi database masih hidup, dipakai k8s
// liveness/readiness probe atau load balancer health check.
func (s *Server) healthCheck(c fiber.Ctx) error {
	if s.sqlDB == nil {
		return response.Error(c, fiber.StatusServiceUnavailable, "database not initialized", nil)
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := s.sqlDB.PingContext(ctx); err != nil {
		logger.FromContext(c.Context()).Error("health check: database ping failed", zap.Error(err))
		return response.Error(c, fiber.StatusServiceUnavailable, "database unreachable", nil)
	}

	return response.Success(c, fiber.StatusOK, "ok", nil)
}

// Start runs the HTTP server and listens for OS signals for graceful shutdown.
func (s *Server) Start() {
	port := s.cfg.AppPort
	if port == "" {
		port = "8000"
	}

	// Setup channel untuk menerima signal shutdown OS
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Jalankan server di goroutine agar tidak memblokir channel signal
	go func() {
		logger.L().Info("starting server", zap.String("app", s.cfg.AppName), zap.String("port", port))
		if err := s.app.Listen(":" + port); err != nil {
			logger.L().Error("server shut down with error", zap.Error(err))
		}
	}()

	// Menunggu signal masuk (Ctrl+C atau kill command)
	<-shutdownChan
	logger.L().Info("shutting down server gracefully...")

	// Matikan server Fiber
	if err := s.app.Shutdown(); err != nil {
		logger.L().Error("error shutting down fiber server", zap.Error(err))
	} else {
		logger.L().Info("fiber server stopped")
	}

	// Close database connection
	if s.sqlDB != nil {
		logger.L().Info("closing database connections...")
		if err := s.sqlDB.Close(); err != nil {
			logger.L().Error("error closing database connection", zap.Error(err))
		} else {
			logger.L().Info("database connection closed successfully")
		}
	}

	// Tambahkan jeda opsional kecil untuk membiarkan resource lain bersih-bersih
	time.Sleep(1 * time.Second)
	logger.L().Info("HRIS backend server gracefully stopped")
}
