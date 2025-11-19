package commands

import (
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type CalculateMortgageCommand struct {
	UserID               string
	PropertyPrice        float64
	DownPayment          float64
	LoanAmount           float64
	BonoTechoPropio      float64
	InterestRate         float64
	RateType             *string // "NOMINAL" o "EFFECTIVE"
	BankID               *string
	PaymentFrequencyDays *int
	DaysInYear           *int
	TermMonths           int
	GracePeriodMonths    int
	GracePeriodType      string  // "NONE", "TOTAL", "PARTIAL"
	Currency             string  // "PEN" o "USD"
	NPVDiscountRate      float64 // Tasa de descuento para calcular VAN (opcional)
}

func NewCalculateMortgageCommand(
	userID string,
	propertyPrice float64,
	downPayment float64,
	loanAmount float64,
	bonoTechoPropio float64,
	interestRate float64,
	rateType *string,
	bankID *string,
	paymentFrequencyDays *int,
	daysInYear *int,
	termMonths int,
	gracePeriodMonths int,
	gracePeriodType string,
	currency string,
	npvDiscountRate float64,
) (*CalculateMortgageCommand, error) {
	// Validaciones básicas
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if propertyPrice <= 0 {
		return nil, errors.New("property price must be greater than zero")
	}
	if loanAmount <= 0 {
		return nil, errors.New("loan amount must be greater than zero")
	}
	if interestRate < 0 {
		return nil, errors.New("interest rate cannot be negative")
	}
	if termMonths <= 0 {
		return nil, errors.New("term months must be greater than zero")
	}
	if gracePeriodMonths < 0 {
		return nil, errors.New("grace period months cannot be negative")
	}
	if gracePeriodMonths >= termMonths {
		return nil, errors.New("grace period months must be less than term months")
	}

	// Validar tipos de enumeraciones
	if rateType != nil {
		if _, err := valueobjects.NewRateType(*rateType); err != nil {
			return nil, err
		}
	}
	if _, err := valueobjects.NewGracePeriodType(gracePeriodType); err != nil {
		return nil, err
	}
	if _, err := valueobjects.NewCurrency(currency); err != nil {
		return nil, err
	}
	if bankID != nil {
		if _, err := valueobjects.NewBankID(*bankID); err != nil {
			return nil, err
		}
	}
	if rateType == nil && bankID == nil {
		return nil, errors.New("either bank_id or rate_type must be provided")
	}
	if paymentFrequencyDays != nil && *paymentFrequencyDays <= 0 {
		return nil, errors.New("payment frequency days must be greater than zero")
	}
	if daysInYear != nil && *daysInYear <= 0 {
		return nil, errors.New("days in year must be greater than zero")
	}

	return &CalculateMortgageCommand{
		UserID:               userID,
		PropertyPrice:        propertyPrice,
		DownPayment:          downPayment,
		LoanAmount:           loanAmount,
		BonoTechoPropio:      bonoTechoPropio,
		InterestRate:         interestRate,
		RateType:             rateType,
		BankID:               bankID,
		PaymentFrequencyDays: paymentFrequencyDays,
		DaysInYear:           daysInYear,
		TermMonths:           termMonths,
		GracePeriodMonths:    gracePeriodMonths,
		GracePeriodType:      gracePeriodType,
		Currency:             currency,
		NPVDiscountRate:      npvDiscountRate,
	}, nil
}
