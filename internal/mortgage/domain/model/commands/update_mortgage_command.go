package commands

import (
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type UpdateMortgageCommand struct {
	mortgageID        valueobjects.MortgageID
	propertyPrice     *float64
	downPayment       *float64
	loanAmount        *float64
	bonoTechoPropio   *float64
	interestRate      *float64
	rateType          *string
	termMonths        *int
	gracePeriodMonths *int
	gracePeriodType   *string
	currency          *string
	npvDiscountRate   *float64
}

func NewUpdateMortgageCommand(
	mortgageID valueobjects.MortgageID,
	propertyPrice *float64,
	downPayment *float64,
	loanAmount *float64,
	bonoTechoPropio *float64,
	interestRate *float64,
	rateType *string,
	termMonths *int,
	gracePeriodMonths *int,
	gracePeriodType *string,
	currency *string,
	npvDiscountRate *float64,
) (*UpdateMortgageCommand, error) {
	if mortgageID.Value() == 0 {
		return nil, errors.New("mortgage ID is required")
	}

	// Validate if any value is provided
	hasUpdates := propertyPrice != nil || downPayment != nil || loanAmount != nil ||
		bonoTechoPropio != nil || interestRate != nil || rateType != nil ||
		termMonths != nil || gracePeriodMonths != nil || gracePeriodType != nil ||
		currency != nil || npvDiscountRate != nil

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
	if gracePeriodMonths != nil && *gracePeriodMonths < 0 {
		return nil, errors.New("grace period months cannot be negative")
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

	return &UpdateMortgageCommand{
		mortgageID:        mortgageID,
		propertyPrice:     propertyPrice,
		downPayment:       downPayment,
		loanAmount:        loanAmount,
		bonoTechoPropio:   bonoTechoPropio,
		interestRate:      interestRate,
		rateType:          rateType,
		termMonths:        termMonths,
		gracePeriodMonths: gracePeriodMonths,
		gracePeriodType:   gracePeriodType,
		currency:          currency,
		npvDiscountRate:   npvDiscountRate,
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
func (c *UpdateMortgageCommand) TermMonths() *int                    { return c.termMonths }
func (c *UpdateMortgageCommand) GracePeriodMonths() *int             { return c.gracePeriodMonths }
func (c *UpdateMortgageCommand) GracePeriodType() *string            { return c.gracePeriodType }
func (c *UpdateMortgageCommand) Currency() *string                   { return c.currency }
func (c *UpdateMortgageCommand) NPVDiscountRate() *float64           { return c.npvDiscountRate }
