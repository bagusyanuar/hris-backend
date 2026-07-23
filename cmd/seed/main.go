package main

import (
	"errors"
	"log"

	"github.com/bagusyanuar/hris-backend/internal/shared/config"
	"github.com/bagusyanuar/hris-backend/internal/shared/database"
	"github.com/bagusyanuar/hris-backend/internal/user/adapter/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Connect database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Seeding database...")

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("@Admin1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	// 4. Prepare seeds
	seeds := []models.UserModel{
		{
			ID:       uuid.New(),
			Email:    "admin@hris.local",
			Password: string(hashedPassword),
			Status:   "active",
		},
		{
			ID:       uuid.New(),
			Email:    "employee@hris.local",
			Password: string(hashedPassword),
			Status:   "active",
		},
	}

	// 5. Insert or update users
	for _, user := range seeds {
		var existing models.UserModel
		err := db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&existing).Error
		if err == nil {
			// Update existing user
			existing.Password = user.Password
			existing.Status = user.Status
			if err := db.Save(&existing).Error; err != nil {
				log.Fatalf("failed to update user %s: %v", user.Email, err)
			}
			log.Printf("Successfully updated user: %s", user.Email)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Insert new user
			if err := db.Create(&user).Error; err != nil {
				log.Fatalf("failed to seed user %s: %v", user.Email, err)
			}
			log.Printf("Successfully seeded user: %s", user.Email)
		} else {
			log.Fatalf("failed checking existing user %s: %v", user.Email, err)
		}
	}

	log.Println("Database seeding completed successfully.")
}
