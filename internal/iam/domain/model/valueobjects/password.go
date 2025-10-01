package valueobjects

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedValue string `json:"-" gorm:"column:password_hash"`
}

func NewPassword(plainPassword string) (Password, error) {
	if len(plainPassword) < 6 {
		return Password{}, errors.New("password must be at least 6 characters long")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hashedValue: string(hashedBytes)}, nil
}

func NewPasswordFromHash(hashedValue string) Password {
	return Password{hashedValue: hashedValue}
}

func (p Password) Hash() string {
	return p.hashedValue
}

func (p Password) Matches(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}
