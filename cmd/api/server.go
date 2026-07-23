package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	authInfra "github.com/bagusyanuar/hris-backend/internal/auth/infrastructure"
	"github.com/bagusyanuar/hris-backend/internal/di"
	"github.com/bagusyanuar/hris-backend/internal/shared/config"
	"github.com/bagusyanuar/hris-backend/internal/shared/middleware"

	"github.com/gofiber/fiber/v3"
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

	server.setupRoutes()
	return server
}

// setupRoutes initializes all routes and their dependencies.
func (s *Server) setupRoutes() {
	// Health Check / Test Route
	s.app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "success",
			"message":     fmt.Sprintf("%s is running", s.cfg.AppName),
			"version":     s.cfg.AppVersion,
			"environment": s.cfg.AppEnv,
		})
	})

	// Initialize Shared Dependencies
	tokenGenerator := authInfra.NewJWTService(s.cfg.JwtSecret, s.cfg.JwtExpiryHour, s.cfg.JwtRefreshExpiryHour)

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
		log.Printf("Starting %s server on port %s...", s.cfg.AppName, port)
		if err := s.app.Listen(":" + port); err != nil {
			log.Printf("Server shut down with error: %v", err)
		}
	}()

	// Menunggu signal masuk (Ctrl+C atau kill command)
	<-shutdownChan
	log.Println("Shutting down server gracefully...")

	// Matikan server Fiber
	if err := s.app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber server: %v", err)
	} else {
		log.Println("Fiber server stopped.")
	}

	// Close database connection
	if s.sqlDB != nil {
		log.Println("Closing database connections...")
		if err := s.sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}

	// Tambahkan jeda opsional kecil untuk membiarkan resource lain bersih-bersih
	time.Sleep(1 * time.Second)
	log.Println("HRIS Backend server gracefully stopped.")
}
