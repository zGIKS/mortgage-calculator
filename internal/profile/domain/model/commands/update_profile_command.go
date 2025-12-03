package commands

import (
	"errors"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
)

type UpdateProfileCommand struct {
	profileID     valueobjects.ProfileID
	phoneNumber   *string
	monthlyIncome *float64
	currency      *string
	maritalStatus *string
	isFirstHome   *bool
	hasOwnLand    *bool
}

func NewUpdateProfileCommand(
	profileID valueobjects.ProfileID,
	phoneNumber *string,
	monthlyIncome *float64,
	currency *string,
	maritalStatus *string,
	isFirstHome *bool,
	hasOwnLand *bool,
) (*UpdateProfileCommand, error) {
	if profileID.IsZero() {
		return nil, errors.New("profile ID is required")
	}

	// Treat empty strings as nil
	if phoneNumber != nil && *phoneNumber == "" {
		phoneNumber = nil
	}
	if currency != nil && *currency == "" {
		currency = nil
	}
	if maritalStatus != nil && *maritalStatus == "" {
		maritalStatus = nil
	}

	// At least one field must be provided
	if phoneNumber == nil && monthlyIncome == nil && currency == nil &&
		maritalStatus == nil && isFirstHome == nil && hasOwnLand == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Validate monthly income if provided
	if monthlyIncome != nil && *monthlyIncome < 0 {
		return nil, errors.New("monthly income cannot be negative")
	}

	return &UpdateProfileCommand{
		profileID:     profileID,
		phoneNumber:   phoneNumber,
		monthlyIncome: monthlyIncome,
		currency:      currency,
		maritalStatus: maritalStatus,
		isFirstHome:   isFirstHome,
		hasOwnLand:    hasOwnLand,
	}, nil
}

func (c *UpdateProfileCommand) ProfileID() valueobjects.ProfileID { return c.profileID }
func (c *UpdateProfileCommand) PhoneNumber() *string              { return c.phoneNumber }
func (c *UpdateProfileCommand) MonthlyIncome() *float64           { return c.monthlyIncome }
func (c *UpdateProfileCommand) Currency() *string                 { return c.currency }
func (c *UpdateProfileCommand) MaritalStatus() *string            { return c.maritalStatus }
func (c *UpdateProfileCommand) IsFirstHome() *bool                { return c.isFirstHome }
func (c *UpdateProfileCommand) HasOwnLand() *bool                 { return c.hasOwnLand }
