package valueobjects

import (
	"errors"
	"regexp"
)

type PhoneNumber struct {
	value string
}

// Peruvian phone numbers are 9 digits
var phoneRegex = regexp.MustCompile(`^\d{9}$`)

func NewPhoneNumber(value string) (PhoneNumber, error) {
	if value == "" {
		return PhoneNumber{}, errors.New("phone number cannot be empty")
	}
	if !phoneRegex.MatchString(value) {
		return PhoneNumber{}, errors.New("phone number must be exactly 9 digits")
	}
	return PhoneNumber{value: value}, nil
}

func (p PhoneNumber) Value() string {
	return p.value
}
