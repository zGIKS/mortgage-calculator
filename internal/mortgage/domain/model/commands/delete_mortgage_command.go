package commands

import (
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type DeleteMortgageCommand struct {
	mortgageID valueobjects.MortgageID
	userID     valueobjects.UserID
}

func NewDeleteMortgageCommand(
	mortgageID valueobjects.MortgageID,
	userID valueobjects.UserID,
) (*DeleteMortgageCommand, error) {
	if mortgageID.Value() == 0 {
		return nil, errors.New("mortgage ID is required")
	}
	if userID.String() == "" {
		return nil, errors.New("user ID is required")
	}

	return &DeleteMortgageCommand{
		mortgageID: mortgageID,
		userID:     userID,
	}, nil
}

// Getters
func (c *DeleteMortgageCommand) MortgageID() valueobjects.MortgageID { return c.mortgageID }
func (c *DeleteMortgageCommand) UserID() valueobjects.UserID         { return c.userID }
