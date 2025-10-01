package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string `json:"value" gorm:"column:email"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(value) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{value: strings.ToLower(value)}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}
