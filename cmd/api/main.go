package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bagusyanuar/hris-backend/internal/application/auth"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/config"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/database"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository"
	"github.com/bagusyanuar/hris-backend/internal/infrastructure/security"
	httpAuth "github.com/bagusyanuar/hris-backend/internal/interfaces/http/auth"
	"github.com/bagusyanuar/hris-backend/internal/interfaces/http/middleware"
	
	applicationOrg "github.com/bagusyanuar/hris-backend/internal/application/organization"
	httpOrg "github.com/bagusyanuar/hris-backend/internal/interfaces/http/organization"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Database
	db, err := database.InitDB(cfg)
	var sqlDB *sql.DB
	if err != nil {
		log.Printf("Warning: Failed to connect database: %v. Server will start but database queries may fail.", err)
	} else {
		// Dapatkan instance sql.DB untuk bisa diclose saat shutdown
		if sDB, err := db.DB(); err == nil {
			sqlDB = sDB
		}
	}

	// Initialize Fiber v3 App
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName + " v" + cfg.AppVersion,
	})

	// Health Check / Test Route
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "success",
			"message":     fmt.Sprintf("%s is running", cfg.AppName),
			"version":     cfg.AppVersion,
			"environment": cfg.AppEnv,
		})
	})

	// Setup Dependencies
	userRepo := repository.NewUserRepository(db)
	tokenGenerator := security.NewJWTService(cfg.JwtSecret, cfg.JwtExpiryHour, cfg.JwtRefreshExpiryHour)
	authService := auth.NewService(userRepo, tokenGenerator)
	authHandler := httpAuth.NewHandler(authService)

	// Organization Dependencies
	orgRepo := repository.NewOrganizationRepository(db)
	orgService := applicationOrg.NewService(orgRepo)
	orgHandler := httpOrg.NewHandler(orgService)

	// Setup Routes
	api := app.Group("/api/v1")
	
	// Auth Routes
	authRoutes := api.Group("/auth")
	authRoutes.Post("/login", authHandler.Login)
	authRoutes.Post("/refresh", authHandler.Refresh)

	// Organization Routes
	httpOrg.RegisterRoutes(api, orgHandler)

	// Protected Example Route
	api.Get("/users/me", middleware.AuthProtected(tokenGenerator), func(c fiber.Ctx) error {
		userID := c.Locals("userID")
		role := c.Locals("role")
		return c.JSON(fiber.Map{
			"user_id": userID,
			"role":    role,
		})
	})

	// Bind server to configured port
	port := cfg.AppPort
	if port == "" {
		port = "8000"
	}

	// Setup channel untuk menerima signal shutdown OS
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Jalankan server di goroutine agar tidak memblokir channel signal
	go func() {
		log.Printf("Starting %s server on port %s...", cfg.AppName, port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server shut down with error: %v", err)
		}
	}()

	// Menunggu signal masuk (Ctrl+C atau kill command)
	<-shutdownChan
	log.Println("Shutting down server gracefully...")

	// Matikan server Fiber
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber server: %v", err)
	} else {
		log.Println("Fiber server stopped.")
	}

	// Close database connection
	if sqlDB != nil {
		log.Println("Closing database connections...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}

	// Tambahkan jeda opsional kecil untuk membiarkan resource lain bersih-bersih
	time.Sleep(1 * time.Second)
	log.Println("HRIS Backend server gracefully stopped.")
}
