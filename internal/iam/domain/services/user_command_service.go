package services

import (
	"context"
	"finanzas-backend/internal/iam/domain/model/commands"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
)

type UserCommandService interface {
	HandleRegister(ctx context.Context, cmd commands.RegisterUserCommand) (*valueobjects.UserID, error)
	HandleUpdate(ctx context.Context, cmd *commands.UpdateUserCommand) error
}
