package main

import (
	"log"

	"github.com/bagusyanuar/hris-backend/internal/shared/config"
	"github.com/bagusyanuar/hris-backend/internal/shared/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect database: %v. Server will start but database queries may fail.", err)
	}

	// Initialize and setup the server
	server := NewServer(cfg, db)

	// Start the server (this will block until shutdown signal is received)
	server.Start()
}
