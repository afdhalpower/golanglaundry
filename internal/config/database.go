package config

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afdhalpower/golanglaundry/internal/models"
)

var DB *gorm.DB

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.App.Debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	DB = db

	slog.Info("database connected successfully",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"name", cfg.Database.Name,
	)

	return db, nil
}

func RunAutoMigration(db *gorm.DB) error {
	slog.Info("running auto migration")

	err := db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Service{},
		&models.Order{},
		&models.Payment{},
		&models.Expense{},
		&models.OrderDetail{},
		&models.OrderTracking{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	slog.Info("auto migration completed")
	return nil
}

func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
