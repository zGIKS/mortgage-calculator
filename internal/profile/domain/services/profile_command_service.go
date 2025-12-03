package services

import (
	"context"
	"finanzas-backend/internal/profile/domain/model/commands"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
)

type ProfileCommandService interface {
	HandleCreate(ctx context.Context, cmd commands.CreateProfileCommand) (*valueobjects.ProfileID, error)
	HandleUpdate(ctx context.Context, cmd *commands.UpdateProfileCommand) error
}
