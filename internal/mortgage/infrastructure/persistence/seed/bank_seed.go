package seed

import (
	"time"

	"finanzas-backend/internal/mortgage/infrastructure/persistence/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedBanks ensures core bank profiles exist (starting with Interbank).
func SeedBanks(db *gorm.DB) error {
	banks := []struct {
		Name                 string
		RateType             string
		PaymentFrequencyDays int
		DaysInYear           int
		IncludesInflation    bool
	}{
		{
			Name:                 "Interbank",
			RateType:             "EFFECTIVE",
			PaymentFrequencyDays: 30,  // cuotas vencidas mensuales por defecto
			DaysInYear:           360, // convención bancaria peruana
			IncludesInflation:    false,
		},
	}

	for _, bankData := range banks {
		// Check if bank already exists by name
		var existingBank models.BankModel
		result := db.Where("name = ?", bankData.Name).First(&existingBank)

		// Only create if not found
		if result.Error == gorm.ErrRecordNotFound {
			now := time.Now()
			newBank := models.BankModel{
				ID:                   uuid.New(),
				Name:                 bankData.Name,
				RateType:             bankData.RateType,
				PaymentFrequencyDays: bankData.PaymentFrequencyDays,
				DaysInYear:           bankData.DaysInYear,
				IncludesInflation:    bankData.IncludesInflation,
				CreatedAt:            now,
				UpdatedAt:            now,
			}

			if err := db.Create(&newBank).Error; err != nil {
				return err
			}
		} else if result.Error != nil {
			// If error is not "record not found", return it
			return result.Error
		}
		// If bank exists, do nothing
	}

	return nil
}
