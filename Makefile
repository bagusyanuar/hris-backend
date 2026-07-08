# Makefile for HRIS Backend

# Load environment variables from .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Database Connection String for golang-migrate
# Constructed from individual DB variables in .env
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)&timezone=$(DB_TZ)

MIGRATION_PATH=./migrations
MIGRATE=migrate

.PHONY: help run dev build lint migrate-create migrate-up migrate-down migrate-force migrate-status migrate-drop migrate-fresh

help: ## Show this help menu
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the application locally
	@echo "Starting application..."
	@go run cmd/api/main.go

dev: ## Run the application with hot reload (Air)
	@echo "Starting application with Air (Hot Reload)..."
	@air

build: ## Build the application binary
	@echo "Building application binary..."
	@go build -o build/api cmd/api/main.go

migrate-create: ## Create a new migration file. Usage: make migrate-create name=init_schema
	@echo "Creating migration: $(name)"
	@$(MIGRATE) create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

migrate-up: ## Run all up migrations
	@echo "Running up migrations..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-down: ## Rollback the last migration
	@echo "Running down migrations..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" down 1

migrate-force: ## Force migration to a specific version. Usage: make migrate-force version=N
	@echo "Force-setting version to $(version)..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" force $(version)

migrate-status: ## Show migration status
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" version

migrate-drop: ## Drop all tables in the database (CAUTION)
	@echo "Dropping all tables..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" drop -f

migrate-fresh: migrate-drop migrate-up db-seed ## Drop all tables, run migrations and seed

db-seed: ## Seed the database with initial data
	@echo "Seeding database..."
	@go run cmd/seed/main.go

lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run
