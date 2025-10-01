package persistence

import (
	"fmt"
	"log"

	"finanzas-backend/internal/shared/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Import models for auto-migration
	iamModels "finanzas-backend/internal/iam/infrastructure/persistence/models"
	mortgageModels "finanzas-backend/internal/mortgage/infrastructure/persistence/models"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&iamModels.UserModel{},
		&mortgageModels.MortgageModel{},
	)
}
