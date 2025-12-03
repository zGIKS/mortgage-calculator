package commands

import (
	"errors"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
)

type UpdateUserCommand struct {
	userID   valueobjects.UserID
	password *string
}

func NewUpdateUserCommand(
	userID valueobjects.UserID,
	password *string,
) (*UpdateUserCommand, error) {
	if userID.IsZero() {
		return nil, errors.New("user ID is required")
	}

	// Validate if password is provided
	if password == nil {
		return nil, errors.New("password must be provided for update")
	}

	// Validate password
	if _, err := valueobjects.NewPassword(*password); err != nil {
		return nil, err
	}

	return &UpdateUserCommand{
		userID:   userID,
		password: password,
	}, nil
}

// Getters
func (c *UpdateUserCommand) UserID() valueobjects.UserID { return c.userID }
func (c *UpdateUserCommand) Password() *string           { return c.password }
