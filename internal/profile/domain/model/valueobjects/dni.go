package valueobjects

import (
	"errors"
	"regexp"
)

type DNI struct {
	value string
}

var dniRegex = regexp.MustCompile(`^\d{8}$`)

func NewDNI(value string) (DNI, error) {
	if value == "" {
		return DNI{}, errors.New("DNI cannot be empty")
	}
	if !dniRegex.MatchString(value) {
		return DNI{}, errors.New("DNI must be exactly 8 digits")
	}
	return DNI{value: value}, nil
}

func (d DNI) Value() string {
	return d.value
}
