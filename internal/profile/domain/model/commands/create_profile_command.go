package commands

import (
	"errors"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
)

type CreateProfileCommand struct {
	userID            valueobjects.UserID
	dni               string
	phoneNumber       string
	monthlyIncome     float64
	currency          string
	maritalStatus     string
	isFirstHome       bool
	hasOwnLand        bool
}

func NewCreateProfileCommand(
	userID valueobjects.UserID,
	dni string,
	phoneNumber string,
	monthlyIncome float64,
	currency string,
	maritalStatus string,
	isFirstHome bool,
	hasOwnLand bool,
) (CreateProfileCommand, error) {
	if userID.IsZero() {
		return CreateProfileCommand{}, errors.New("user ID is required")
	}
	if dni == "" {
		return CreateProfileCommand{}, errors.New("DNI is required")
	}
	if phoneNumber == "" {
		return CreateProfileCommand{}, errors.New("phone number is required")
	}
	if monthlyIncome < 0 {
		return CreateProfileCommand{}, errors.New("monthly income cannot be negative")
	}
	if currency == "" {
		return CreateProfileCommand{}, errors.New("currency is required")
	}
	if maritalStatus == "" {
		return CreateProfileCommand{}, errors.New("marital status is required")
	}

	return CreateProfileCommand{
		userID:        userID,
		dni:           dni,
		phoneNumber:   phoneNumber,
		monthlyIncome: monthlyIncome,
		currency:      currency,
		maritalStatus: maritalStatus,
		isFirstHome:   isFirstHome,
		hasOwnLand:    hasOwnLand,
	}, nil
}

func (c CreateProfileCommand) UserID() valueobjects.UserID { return c.userID }
func (c CreateProfileCommand) DNI() string                 { return c.dni }
func (c CreateProfileCommand) PhoneNumber() string         { return c.phoneNumber }
func (c CreateProfileCommand) MonthlyIncome() float64      { return c.monthlyIncome }
func (c CreateProfileCommand) Currency() string            { return c.currency }
func (c CreateProfileCommand) MaritalStatus() string       { return c.maritalStatus }
func (c CreateProfileCommand) IsFirstHome() bool           { return c.isFirstHome }
func (c CreateProfileCommand) HasOwnLand() bool            { return c.hasOwnLand }
