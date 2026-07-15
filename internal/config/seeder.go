package config

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/afdhalpower/golanglaundry/internal/models"
)

func RunSeeder(db *gorm.DB) error {
	slog.Info("running database seeder")

	// Check if admin already exists
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		slog.Info("seeder skipped: admin user already exists")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Name:         "Administrator",
		Email:        "admin@laundry.com",
		PasswordHash: string(hash),
		Role:         "admin",
		IsActive:     true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	slog.Info("default admin user created", "email", admin.Email)
	return nil
}
