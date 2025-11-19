package valueobjects

import (
	"errors"
	"strings"
)

// BankID represents the identifier of a financial institution that provides mortgage rates.
type BankID struct {
	value string `json:"value" gorm:"column:bank_id"`
}

// NewBankID validates and creates a new BankID.
func NewBankID(value string) (BankID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return BankID{}, errors.New("bank ID cannot be empty")
	}
	if len(value) > 50 {
		return BankID{}, errors.New("bank ID cannot exceed 50 characters")
	}
	return BankID{value: strings.ToUpper(value)}, nil
}

// Value returns the raw identifier.
func (b BankID) Value() string {
	return b.value
}

// String satisfies fmt.Stringer.
func (b BankID) String() string {
	return b.value
}
