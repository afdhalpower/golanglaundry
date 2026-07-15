package config

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/afdhalpower/golanglaundry/internal/models"
)

func RunSeeder(db *gorm.DB) error {
	slog.Info("running database seeder")

	// Seed admin user
	if err := seedAdmin(db); err != nil {
		return err
	}

	// Seed expense categories
	if err := seedExpenseCategories(db); err != nil {
		return err
	}

	return nil
}

func seedAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
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

func seedExpenseCategories(db *gorm.DB) error {
	categories := []string{
		"Sabun",
		"Pewangi",
		"Plastik",
		"Air",
		"Listrik",
		"Gaji",
		"Transportasi",
		"Lainnya",
	}

	for _, name := range categories {
		var count int64
		db.Model(&models.ExpenseCategory{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			if err := db.Create(&models.ExpenseCategory{Name: name}).Error; err != nil {
				return err
			}
			slog.Info("expense category created", "name", name)
		}
	}

	return nil
}
