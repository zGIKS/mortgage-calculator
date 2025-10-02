package commands

import (
	"errors"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
)

type UpdateUserCommand struct {
	userID   valueobjects.UserID
	email    *string
	password *string
}

func NewUpdateUserCommand(
	userID valueobjects.UserID,
	email *string,
	password *string,
) (*UpdateUserCommand, error) {
	if userID.IsZero() {
		return nil, errors.New("user ID is required")
	}

	// Validate if any value is provided
	if email == nil && password == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Validate email if provided
	if email != nil {
		if _, err := valueobjects.NewEmail(*email); err != nil {
			return nil, err
		}
	}

	// Validate password if provided
	if password != nil {
		if _, err := valueobjects.NewPassword(*password); err != nil {
			return nil, err
		}
	}

	return &UpdateUserCommand{
		userID:   userID,
		email:    email,
		password: password,
	}, nil
}

// Getters
func (c *UpdateUserCommand) UserID() valueobjects.UserID { return c.userID }
func (c *UpdateUserCommand) Email() *string              { return c.email }
func (c *UpdateUserCommand) Password() *string           { return c.password }
