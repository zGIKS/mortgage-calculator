package seed

import (
	"time"

	"finanzas-backend/internal/mortgage/infrastructure/persistence/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedBanks ensures core bank profiles exist.
// Datos basados en información de la SBS (Superintendencia de Banca, Seguros y AFP)
// Fuente: https://www.sbs.gob.pe/app/pp/EstadisticasSAEEPortal/Paginas/TIActivaTipoCreditoEmpresa.aspx?tip=B
// Regulación: Resolución S.B.S. N° 3274-2017 (año comercial de 360 días)
func SeedBanks(db *gorm.DB) error {
	banks := []struct {
		Name                 string
		RateType             string
		PaymentFrequencyDays int
		DaysInYear           int
		IncludesInflation    bool
	}{
		{
			// Interbank - TEA hipotecaria: 7.72% (Nov 2025)
			// Fuente: https://interbank.pe/tasas-tarifas
			Name:                 "Interbank",
			RateType:             "EFFECTIVE",
			PaymentFrequencyDays: 30,  // cuotas mensuales
			DaysInYear:           360, // convención bancaria peruana (SBS)
			IncludesInflation:    false,
		},
		{
			// BBVA Perú - TEA hipotecaria: 7.65% (Nov 2025)
			// Fuente: https://www.bbva.pe/personas/productos/prestamos/credito-hipotecario.html
			Name:                 "BBVA",
			RateType:             "EFFECTIVE",
			PaymentFrequencyDays: 30,  // cuotas mensuales
			DaysInYear:           360, // convención bancaria peruana (SBS)
			IncludesInflation:    false,
		},
		{
			// BCP (Banco de Crédito del Perú) - TEA hipotecaria: 8.02% (Nov 2025)
			// Fuente: https://www.viabcp.com/tasasytarifas
			Name:                 "BCP",
			RateType:             "EFFECTIVE",
			PaymentFrequencyDays: 30,  // cuotas mensuales
			DaysInYear:           360, // convención bancaria peruana (SBS)
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
