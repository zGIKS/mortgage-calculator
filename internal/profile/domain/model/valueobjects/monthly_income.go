package valueobjects

import (
	"errors"
	"fmt"
)

type Currency string

const (
	CurrencyPEN Currency = "PEN" // Peruvian Soles
	CurrencyUSD Currency = "USD" // US Dollars
)

type MonthlyIncome struct {
	amount   float64
	currency Currency
}

func NewMonthlyIncome(amount float64, currency Currency) (MonthlyIncome, error) {
	if amount < 0 {
		return MonthlyIncome{}, errors.New("monthly income cannot be negative")
	}
	if currency != CurrencyPEN && currency != CurrencyUSD {
		return MonthlyIncome{}, errors.New("currency must be PEN or USD")
	}
	return MonthlyIncome{amount: amount, currency: currency}, nil
}

func (m MonthlyIncome) Amount() float64 {
	return m.amount
}

func (m MonthlyIncome) Currency() Currency {
	return m.currency
}

func (m MonthlyIncome) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}
