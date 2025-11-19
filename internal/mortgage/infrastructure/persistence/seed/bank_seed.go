package seed

import (
	"time"

	"finanzas-backend/internal/mortgage/infrastructure/persistence/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedBanks ensures core bank profiles exist (starting with Interbank).
func SeedBanks(db *gorm.DB) error {
	now := time.Now()
	records := []models.BankModel{
		{
			ID:                   "INTERBANK",
			Name:                 "Interbank",
			RateType:             "EFFECTIVE",
			PaymentFrequencyDays: 30,  // cuotas vencidas mensuales por defecto
			DaysInYear:           360, // convención bancaria peruana
			IncludesInflation:    false,
			CreatedAt:            now,
			UpdatedAt:            now,
		},
	}

	for _, record := range records {
		entity := record
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&entity).Error
		if err != nil {
			return err
		}
	}

	return nil
}
