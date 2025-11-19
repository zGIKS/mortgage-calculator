package entities

import (
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"strings"
	"time"
)

// Bank encapsulates the configuration required to interpret a bank mortgage offer.
type Bank struct {
	id                    valueobjects.BankID
	name                  string
	rateType              valueobjects.RateType
	paymentFrequencyDays  int
	daysInYear            int
	includesInflationRate bool
	createdAt             time.Time
	updatedAt             time.Time
}

// NewBank creates a bank aggregate ensuring required invariants.
func NewBank(
	id valueobjects.BankID,
	name string,
	rateType valueobjects.RateType,
	paymentFrequencyDays int,
	daysInYear int,
	includesInflation bool,
) (*Bank, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return nil, errors.New("bank name cannot be empty")
	}
	if paymentFrequencyDays <= 0 {
		return nil, errors.New("payment frequency days must be greater than zero")
	}
	if daysInYear <= 0 {
		return nil, errors.New("days in year must be greater than zero")
	}

	return &Bank{
		id:                    id,
		name:                  name,
		rateType:              rateType,
		paymentFrequencyDays:  paymentFrequencyDays,
		daysInYear:            daysInYear,
		includesInflationRate: includesInflation,
		createdAt:             time.Now(),
		updatedAt:             time.Now(),
	}, nil
}

// ReconstructBank rebuilds the aggregate from persistence.
func ReconstructBank(
	id valueobjects.BankID,
	name string,
	rateType valueobjects.RateType,
	paymentFrequencyDays int,
	daysInYear int,
	includesInflation bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Bank {
	return &Bank{
		id:                    id,
		name:                  name,
		rateType:              rateType,
		paymentFrequencyDays:  paymentFrequencyDays,
		daysInYear:            daysInYear,
		includesInflationRate: includesInflation,
		createdAt:             createdAt,
		updatedAt:             updatedAt,
	}
}

// Getters
func (b *Bank) ID() valueobjects.BankID         { return b.id }
func (b *Bank) Name() string                    { return b.name }
func (b *Bank) RateType() valueobjects.RateType { return b.rateType }
func (b *Bank) PaymentFrequencyDays() int       { return b.paymentFrequencyDays }
func (b *Bank) DaysInYear() int                 { return b.daysInYear }
func (b *Bank) IncludesInflationRate() bool     { return b.includesInflationRate }
func (b *Bank) CreatedAt() time.Time            { return b.createdAt }
func (b *Bank) UpdatedAt() time.Time            { return b.updatedAt }
