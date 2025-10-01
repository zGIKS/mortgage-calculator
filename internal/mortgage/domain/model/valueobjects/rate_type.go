package valueobjects

import "errors"

type RateType string

const (
	RateTypeNominal   RateType = "NOMINAL"   // Tasa Nominal Anual (TNA)
	RateTypeEffective RateType = "EFFECTIVE" // Tasa Efectiva Anual (TEA)
)

func NewRateType(value string) (RateType, error) {
	rt := RateType(value)
	switch rt {
	case RateTypeNominal, RateTypeEffective:
		return rt, nil
	default:
		return "", errors.New("invalid rate type, must be NOMINAL or EFFECTIVE")
	}
}

func (r RateType) String() string {
	return string(r)
}
