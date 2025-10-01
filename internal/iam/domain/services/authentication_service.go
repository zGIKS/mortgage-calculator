package services

import (
	"context"
	"finanzas-backend/internal/iam/domain/model/commands"
)

type AuthenticationService interface {
	HandleLogin(ctx context.Context, cmd commands.LoginCommand) (string, error)
}
