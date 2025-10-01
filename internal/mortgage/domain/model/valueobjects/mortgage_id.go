package valueobjects

import "errors"

type MortgageID struct {
	value uint64
}

func NewMortgageID(value uint64) (MortgageID, error) {
	if value == 0 {
		return MortgageID{}, errors.New("mortgage ID cannot be zero")
	}
	return MortgageID{value: value}, nil
}

func (m MortgageID) Value() uint64 {
	return m.value
}
