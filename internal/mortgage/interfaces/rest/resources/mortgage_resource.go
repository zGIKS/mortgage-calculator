package resources

import (
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"time"
)

// CalculateMortgageRequest representa la solicitud para calcular un crédito hipotecario
type CalculateMortgageRequest struct {
	PropertyPrice     float64 `json:"property_price" binding:"required,gt=0"`
	DownPayment       float64 `json:"down_payment" binding:"required,gte=0"`
	LoanAmount        float64 `json:"loan_amount" binding:"required,gt=0"`
	BonoTechoPropio   float64 `json:"bono_techo_propio" binding:"gte=0"`
	InterestRate      float64 `json:"interest_rate" binding:"required,gte=0"`
	BankID            string  `json:"bank_id" binding:"required"`
	TermMonths        int     `json:"term_months" binding:"required,gt=0"`
	GracePeriodMonths int     `json:"grace_period_months" binding:"gte=0"`
	GracePeriodType   string  `json:"grace_period_type" binding:"required,oneof=NONE TOTAL PARTIAL"`
	Currency          string  `json:"currency" binding:"required,oneof=PEN USD"`
	NPVDiscountRate   float64 `json:"npv_discount_rate" binding:"gte=0"`
}

// UpdateMortgageRequest representa la solicitud para actualizar un crédito hipotecario
type UpdateMortgageRequest struct {
	PropertyPrice        *float64 `json:"property_price,omitempty" binding:"omitempty,gt=0"`
	DownPayment          *float64 `json:"down_payment,omitempty" binding:"omitempty,gte=0"`
	LoanAmount           *float64 `json:"loan_amount,omitempty" binding:"omitempty,gt=0"`
	BonoTechoPropio      *float64 `json:"bono_techo_propio,omitempty" binding:"omitempty,gte=0"`
	InterestRate         *float64 `json:"interest_rate,omitempty" binding:"omitempty,gte=0"`
	RateType             *string  `json:"rate_type,omitempty" binding:"omitempty,oneof=NOMINAL EFFECTIVE"`
	BankID               *string  `json:"bank_id,omitempty"`
	PaymentFrequencyDays *int     `json:"payment_frequency_days,omitempty" binding:"omitempty,gt=0"`
	DaysInYear           *int     `json:"days_in_year,omitempty" binding:"omitempty,gt=0"`
	TermMonths           *int     `json:"term_months,omitempty" binding:"omitempty,gt=0"`
	GracePeriodMonths    *int     `json:"grace_period_months,omitempty" binding:"omitempty,gte=0"`
	GracePeriodType      *string  `json:"grace_period_type,omitempty" binding:"omitempty,oneof=NONE TOTAL PARTIAL"`
	Currency             *string  `json:"currency,omitempty" binding:"omitempty,oneof=PEN USD"`
	NPVDiscountRate      *float64 `json:"npv_discount_rate,omitempty" binding:"omitempty,gte=0"`
}

// PaymentScheduleItemResource representa un item del cronograma
type PaymentScheduleItemResource struct {
	Period           int     `json:"period"`
	Installment      float64 `json:"installment"`
	Interest         float64 `json:"interest"`
	Amortization     float64 `json:"amortization"`
	RemainingBalance float64 `json:"remaining_balance"`
	IsGracePeriod    bool    `json:"is_grace_period"`
}

// MortgageResponse representa la respuesta completa con todos los cálculos
type MortgageResponse struct {
	ID                   uint64  `json:"id"`
	UserID               string  `json:"user_id"`
	PropertyPrice        float64 `json:"property_price"`
	DownPayment          float64 `json:"down_payment"`
	LoanAmount           float64 `json:"loan_amount"`
	BonoTechoPropio      float64 `json:"bono_techo_propio"`
	InterestRate         float64 `json:"interest_rate"`
	RateType             string  `json:"rate_type"`
	BankID               *string `json:"bank_id,omitempty"`
	BankName             string  `json:"bank_name,omitempty"`
	TermMonths           int     `json:"term_months"`
	GracePeriodMonths    int     `json:"grace_period_months"`
	GracePeriodType      string  `json:"grace_period_type"`
	Currency             string  `json:"currency"`
	PaymentFrequencyDays int     `json:"payment_frequency_days"`
	DaysInYear           int     `json:"days_in_year"`

	// Resultados calculados
	PrincipalFinanced float64                       `json:"principal_financed"`
	PeriodicRate      float64                       `json:"periodic_rate"`
	FixedInstallment  float64                       `json:"fixed_installment"`
	PaymentSchedule   []PaymentScheduleItemResource `json:"payment_schedule"`
	TotalInterestPaid float64                       `json:"total_interest_paid"`
	TotalPaid         float64                       `json:"total_paid"`
	NPV               float64                       `json:"npv"`
	IRR               float64                       `json:"irr"`
	TCEA              float64                       `json:"tcea"`

	CreatedAt time.Time `json:"created_at"`
}

// MortgageSummaryResource representa un resumen de hipoteca (para listas)
type MortgageSummaryResource struct {
	ID               uint64    `json:"id"`
	UserID           string    `json:"user_id"`
	PropertyPrice    float64   `json:"property_price"`
	LoanAmount       float64   `json:"loan_amount"`
	Currency         string    `json:"currency"`
	TermMonths       int       `json:"term_months"`
	FixedInstallment float64   `json:"fixed_installment"`
	TCEA             float64   `json:"tcea"`
	BankID           *string   `json:"bank_id,omitempty"`
	BankName         string    `json:"bank_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// TransformToMortgageResponse transforma una entidad Mortgage a MortgageResponse
func TransformToMortgageResponse(mortgage *entities.Mortgage) MortgageResponse {
	scheduleItems := make([]PaymentScheduleItemResource, 0)
	if mortgage.PaymentSchedule() != nil {
		for _, item := range mortgage.PaymentSchedule().GetItems() {
			scheduleItems = append(scheduleItems, PaymentScheduleItemResource{
				Period:           item.Period,
				Installment:      item.Installment,
				Interest:         item.Interest,
				Amortization:     item.Amortization,
				RemainingBalance: item.RemainingBalance,
				IsGracePeriod:    item.IsGracePeriod,
			})
		}
	}

	var bankID *string
	if mortgage.BankID() != nil {
		value := mortgage.BankID().String()
		bankID = &value
	}

	return MortgageResponse{
		ID:                   mortgage.ID().Value(),
		UserID:               mortgage.UserID().String(),
		PropertyPrice:        mortgage.PropertyPrice(),
		DownPayment:          mortgage.DownPayment(),
		LoanAmount:           mortgage.LoanAmount(),
		BonoTechoPropio:      mortgage.BonoTechoPropio(),
		InterestRate:         mortgage.InterestRate(),
		RateType:             mortgage.RateType().String(),
		BankID:               bankID,
		BankName:             mortgage.BankName(),
		TermMonths:           mortgage.TermMonths(),
		GracePeriodMonths:    mortgage.GracePeriodMonths(),
		GracePeriodType:      mortgage.GracePeriodType().String(),
		Currency:             mortgage.Currency().String(),
		PaymentFrequencyDays: mortgage.PaymentFrequencyDays(),
		DaysInYear:           mortgage.DaysInYear(),
		PrincipalFinanced:    mortgage.PrincipalFinanced(),
		PeriodicRate:         mortgage.PeriodicRate(),
		FixedInstallment:     mortgage.FixedInstallment(),
		PaymentSchedule:      scheduleItems,
		TotalInterestPaid:    mortgage.TotalInterestPaid(),
		TotalPaid:            mortgage.TotalPaid(),
		NPV:                  mortgage.NPV(),
		IRR:                  mortgage.IRR(),
		TCEA:                 mortgage.TCEA(),
		CreatedAt:            mortgage.CreatedAt(),
	}
}

// TransformToMortgageSummary transforma una entidad Mortgage a MortgageSummaryResource
func TransformToMortgageSummary(mortgage *entities.Mortgage) MortgageSummaryResource {
	var bankID *string
	if mortgage.BankID() != nil {
		value := mortgage.BankID().String()
		bankID = &value
	}

	return MortgageSummaryResource{
		ID:               mortgage.ID().Value(),
		UserID:           mortgage.UserID().String(),
		PropertyPrice:    mortgage.PropertyPrice(),
		LoanAmount:       mortgage.LoanAmount(),
		Currency:         mortgage.Currency().String(),
		TermMonths:       mortgage.TermMonths(),
		FixedInstallment: mortgage.FixedInstallment(),
		TCEA:             mortgage.TCEA(),
		BankID:           bankID,
		BankName:         mortgage.BankName(),
		CreatedAt:        mortgage.CreatedAt(),
	}
}
