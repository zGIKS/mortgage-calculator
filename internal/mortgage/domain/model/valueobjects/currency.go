package valueobjects

import "errors"

type Currency string

const (
	CurrencyPEN Currency = "PEN" // Soles
	CurrencyUSD Currency = "USD" // Dólares
)

func NewCurrency(value string) (Currency, error) {
	curr := Currency(value)
	switch curr {
	case CurrencyPEN, CurrencyUSD:
		return curr, nil
	default:
		return "", errors.New("invalid currency, must be PEN or USD")
	}
}

func (c Currency) String() string {
	return string(c)
}
