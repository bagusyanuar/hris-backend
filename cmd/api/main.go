package main

import (
	"log"

	"github.com/bagusyanuar/hris-backend/internal/shared/config"
	"github.com/bagusyanuar/hris-backend/internal/shared/database"
	"github.com/bagusyanuar/hris-backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize global structured logger
	zapLogger := logger.Init(cfg.AppEnv, cfg.AppDebug)
	// Sync flushing stdout/stderr commonly returns a harmless "invalid argument"
	// error on some platforms (known zap/Go issue) — safe to discard here.
	defer func() { _ = logger.Sync() }()

	// Initialize Database
	db, err := database.InitDB(cfg)
	if err != nil {
		zapLogger.Fatal("failed to connect database", zap.Error(err))
	}

	// Initialize and setup the server
	server := NewServer(cfg, db)

	// Start the server (this will block until shutdown signal is received)
	server.Start()
}
