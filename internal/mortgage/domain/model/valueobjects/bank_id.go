package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// BankID represents the identifier of a financial institution that provides mortgage rates.
type BankID struct {
	value uuid.UUID `gorm:"column:bank_id;type:uuid"`
}

// NewBankID validates and creates a new BankID.
func NewBankID(value uuid.UUID) (BankID, error) {
	if value == uuid.Nil {
		return BankID{}, errors.New("bank ID cannot be nil")
	}
	return BankID{value: value}, nil
}

// NewBankIDFromString creates a BankID from a string UUID.
func NewBankIDFromString(value string) (BankID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return BankID{}, errors.New("invalid UUID format")
	}
	return NewBankID(parsed)
}

// Value returns the raw identifier.
func (b BankID) Value() uuid.UUID {
	return b.value
}

// String satisfies fmt.Stringer.
func (b BankID) String() string {
	return b.value.String()
}

// IsZero returns true if the BankID is the zero value.
func (b BankID) IsZero() bool {
	return b.value == uuid.Nil
}
