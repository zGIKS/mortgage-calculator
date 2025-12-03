package commands

import (
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type UpdateMortgageCommand struct {
	mortgageID           valueobjects.MortgageID
	propertyPrice        *float64
	downPayment          *float64
	loanAmount           *float64
	bonoTechoPropio      *float64
	interestRate         *float64
	rateType             *string
	paymentFrequencyDays *int
	daysInYear           *int
	termMonths           *int
	termYears            *int
	gracePeriodMonths    *int
	gracePeriodType      *string
	currency             *string
	npvDiscountRate      *float64
	administrationFee    *float64
	portes               *float64
	additionalCosts      *float64
	lifeInsuranceRate    *float64
	propertyInsurance    *float64
	evaluationFee        *float64
	disbursementFee      *float64
}

func NewUpdateMortgageCommand(
	mortgageID valueobjects.MortgageID,
	propertyPrice *float64,
	downPayment *float64,
	loanAmount *float64,
	bonoTechoPropio *float64,
	interestRate *float64,
	rateType *string,
	paymentFrequencyDays *int,
	daysInYear *int,
	termMonths *int,
	termYears *int,
	gracePeriodMonths *int,
	gracePeriodType *string,
	currency *string,
	npvDiscountRate *float64,
	administrationFee *float64,
	portes *float64,
	additionalCosts *float64,
	lifeInsuranceRate *float64,
	propertyInsurance *float64,
	evaluationFee *float64,
	disbursementFee *float64,
) (*UpdateMortgageCommand, error) {
	if mortgageID.Value() == 0 {
		return nil, errors.New("mortgage ID is required")
	}

	// Validate if any value is provided
	hasUpdates := propertyPrice != nil || downPayment != nil || loanAmount != nil ||
		bonoTechoPropio != nil || interestRate != nil || rateType != nil ||
		paymentFrequencyDays != nil || daysInYear != nil ||
		termMonths != nil || termYears != nil || gracePeriodMonths != nil || gracePeriodType != nil ||
		currency != nil || npvDiscountRate != nil || administrationFee != nil || portes != nil ||
		additionalCosts != nil || lifeInsuranceRate != nil || propertyInsurance != nil ||
		evaluationFee != nil || disbursementFee != nil

	if !hasUpdates {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Validate values if provided
	if propertyPrice != nil && *propertyPrice <= 0 {
		return nil, errors.New("property price must be greater than zero")
	}
	if downPayment != nil && *downPayment < 0 {
		return nil, errors.New("down payment cannot be negative")
	}
	if loanAmount != nil && *loanAmount <= 0 {
		return nil, errors.New("loan amount must be greater than zero")
	}
	if bonoTechoPropio != nil && *bonoTechoPropio < 0 {
		return nil, errors.New("bono techo propio cannot be negative")
	}
	if interestRate != nil && *interestRate < 0 {
		return nil, errors.New("interest rate cannot be negative")
	}
	if termMonths != nil && *termMonths <= 0 {
		return nil, errors.New("term months must be greater than zero")
	}
	if termYears != nil && *termYears <= 0 {
		return nil, errors.New("term years must be greater than zero")
	}
	if gracePeriodMonths != nil && *gracePeriodMonths < 0 {
		return nil, errors.New("grace period months cannot be negative")
	}
	if gracePeriodMonths != nil && termMonths != nil && *gracePeriodMonths >= *termMonths {
		return nil, errors.New("grace period months must be less than term months")
	}
	if administrationFee != nil && *administrationFee < 0 {
		return nil, errors.New("fees cannot be negative")
	}
	if portes != nil && *portes < 0 {
		return nil, errors.New("portes cannot be negative")
	}
	if additionalCosts != nil && *additionalCosts < 0 {
		return nil, errors.New("additional costs cannot be negative")
	}
	if lifeInsuranceRate != nil && *lifeInsuranceRate < 0 {
		return nil, errors.New("life insurance rate cannot be negative")
	}
	if propertyInsurance != nil && *propertyInsurance < 0 {
		return nil, errors.New("property insurance rate cannot be negative")
	}
	if evaluationFee != nil && *evaluationFee < 0 {
		return nil, errors.New("evaluation fee cannot be negative")
	}
	if disbursementFee != nil && *disbursementFee < 0 {
		return nil, errors.New("disbursement fee cannot be negative")
	}

	// Validate enumerations if provided
	if rateType != nil {
		if _, err := valueobjects.NewRateType(*rateType); err != nil {
			return nil, err
		}
	}
	if gracePeriodType != nil {
		if _, err := valueobjects.NewGracePeriodType(*gracePeriodType); err != nil {
			return nil, err
		}
	}
	if currency != nil {
		if _, err := valueobjects.NewCurrency(*currency); err != nil {
			return nil, err
		}
	}

	if paymentFrequencyDays != nil && *paymentFrequencyDays <= 0 {
		return nil, errors.New("payment frequency days must be greater than zero")
	}
	if daysInYear != nil && *daysInYear <= 0 {
		return nil, errors.New("days in year must be greater than zero")
	}

	return &UpdateMortgageCommand{
		mortgageID:           mortgageID,
		propertyPrice:        propertyPrice,
		downPayment:          downPayment,
		loanAmount:           loanAmount,
		bonoTechoPropio:      bonoTechoPropio,
		interestRate:         interestRate,
		rateType:             rateType,
		paymentFrequencyDays: paymentFrequencyDays,
		daysInYear:           daysInYear,
		termMonths:           termMonths,
		termYears:            termYears,
		gracePeriodMonths:    gracePeriodMonths,
		gracePeriodType:      gracePeriodType,
		currency:             currency,
		npvDiscountRate:      npvDiscountRate,
		administrationFee:    administrationFee,
		portes:               portes,
		additionalCosts:      additionalCosts,
		lifeInsuranceRate:    lifeInsuranceRate,
		propertyInsurance:    propertyInsurance,
		evaluationFee:        evaluationFee,
		disbursementFee:      disbursementFee,
	}, nil
}

// Getters
func (c *UpdateMortgageCommand) MortgageID() valueobjects.MortgageID { return c.mortgageID }
func (c *UpdateMortgageCommand) PropertyPrice() *float64             { return c.propertyPrice }
func (c *UpdateMortgageCommand) DownPayment() *float64               { return c.downPayment }
func (c *UpdateMortgageCommand) LoanAmount() *float64                { return c.loanAmount }
func (c *UpdateMortgageCommand) BonoTechoPropio() *float64           { return c.bonoTechoPropio }
func (c *UpdateMortgageCommand) InterestRate() *float64              { return c.interestRate }
func (c *UpdateMortgageCommand) RateType() *string                   { return c.rateType }
func (c *UpdateMortgageCommand) PaymentFrequencyDays() *int          { return c.paymentFrequencyDays }
func (c *UpdateMortgageCommand) DaysInYear() *int                    { return c.daysInYear }
func (c *UpdateMortgageCommand) TermMonths() *int                    { return c.termMonths }
func (c *UpdateMortgageCommand) TermYears() *int                     { return c.termYears }
func (c *UpdateMortgageCommand) GracePeriodMonths() *int             { return c.gracePeriodMonths }
func (c *UpdateMortgageCommand) GracePeriodType() *string            { return c.gracePeriodType }
func (c *UpdateMortgageCommand) Currency() *string                   { return c.currency }
func (c *UpdateMortgageCommand) NPVDiscountRate() *float64           { return c.npvDiscountRate }
func (c *UpdateMortgageCommand) AdministrationFee() *float64         { return c.administrationFee }
func (c *UpdateMortgageCommand) Portes() *float64                    { return c.portes }
func (c *UpdateMortgageCommand) AdditionalCosts() *float64           { return c.additionalCosts }
func (c *UpdateMortgageCommand) LifeInsuranceRate() *float64         { return c.lifeInsuranceRate }
func (c *UpdateMortgageCommand) PropertyInsurance() *float64         { return c.propertyInsurance }
func (c *UpdateMortgageCommand) EvaluationFee() *float64             { return c.evaluationFee }
func (c *UpdateMortgageCommand) DisbursementFee() *float64           { return c.disbursementFee }
