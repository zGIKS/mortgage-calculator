package resources

import (
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"time"
)

// BankResource representa la respuesta de un banco
type BankResource struct {
	ID                   string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                 string    `json:"name" example:"Banco de Crédito del Perú"`
	RateType             string    `json:"rate_type" example:"EFFECTIVE"`
	PaymentFrequencyDays int       `json:"payment_frequency_days" example:"30"`
	DaysInYear           int       `json:"days_in_year" example:"360"`
	IncludesInflation    bool      `json:"includes_inflation" example:"false"`
	CreatedAt            time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// TransformToBankResource transforma una entidad Bank a BankResource
func TransformToBankResource(bank *entities.Bank) BankResource {
	return BankResource{
		ID:                   bank.ID().String(),
		Name:                 bank.Name(),
		RateType:             bank.RateType().String(),
		PaymentFrequencyDays: bank.PaymentFrequencyDays(),
		DaysInYear:           bank.DaysInYear(),
		IncludesInflation:    bank.IncludesInflationRate(),
		CreatedAt:            bank.CreatedAt(),
	}
}

// TransformToBankResources transforma una lista de entidades Bank a BankResource
func TransformToBankResources(banks []*entities.Bank) []BankResource {
	resources := make([]BankResource, 0, len(banks))
	for _, bank := range banks {
		resources = append(resources, TransformToBankResource(bank))
	}
	return resources
}
