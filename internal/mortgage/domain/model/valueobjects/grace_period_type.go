package valueobjects

import "errors"

type GracePeriodType string

const (
	GracePeriodNone    GracePeriodType = "NONE"    // Sin periodo de gracia
	GracePeriodTotal   GracePeriodType = "TOTAL"   // Gracia total (no paga ni capital ni intereses)
	GracePeriodPartial GracePeriodType = "PARTIAL" // Gracia parcial (solo paga intereses)
)

func NewGracePeriodType(value string) (GracePeriodType, error) {
	gpt := GracePeriodType(value)
	switch gpt {
	case GracePeriodNone, GracePeriodTotal, GracePeriodPartial:
		return gpt, nil
	default:
		return "", errors.New("invalid grace period type, must be NONE, TOTAL or PARTIAL")
	}
}

func (g GracePeriodType) String() string {
	return string(g)
}
