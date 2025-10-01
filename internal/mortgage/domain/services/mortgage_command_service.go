package services

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/commands"
	"finanzas-backend/internal/mortgage/domain/model/entities"
)

type MortgageCommandService interface {
	HandleCalculateMortgage(ctx context.Context, cmd *commands.CalculateMortgageCommand) (*entities.Mortgage, error)
}
