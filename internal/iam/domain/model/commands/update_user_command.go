package commands

import (
	"errors"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
)

type UpdateUserCommand struct {
	userID   valueobjects.UserID
	email    *string
	password *string
	fullName *string
}

func NewUpdateUserCommand(
	userID valueobjects.UserID,
	email *string,
	password *string,
	fullName *string,
) (*UpdateUserCommand, error) {
	if userID.IsZero() {
		return nil, errors.New("user ID is required")
	}

	// Validate if any value is provided
	if email == nil && password == nil && fullName == nil {
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

	// Validate full name if provided
	if fullName != nil {
		if len(*fullName) < 2 {
			return nil, errors.New("full name must be at least 2 characters")
		}
	}

	return &UpdateUserCommand{
		userID:   userID,
		email:    email,
		password: password,
		fullName: fullName,
	}, nil
}

// Getters
func (c *UpdateUserCommand) UserID() valueobjects.UserID { return c.userID }
func (c *UpdateUserCommand) Email() *string              { return c.email }
func (c *UpdateUserCommand) Password() *string           { return c.password }
func (c *UpdateUserCommand) FullName() *string           { return c.fullName }
