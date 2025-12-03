package repositories

import (
	"context"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
)

type ProfileRepository interface {
	Save(ctx context.Context, profile *entities.Profile) error
	Update(ctx context.Context, profile *entities.Profile) error
	FindByID(ctx context.Context, id valueobjects.ProfileID) (*entities.Profile, error)
	FindByUserID(ctx context.Context, userID valueobjects.UserID) (*entities.Profile, error)
	FindByDNI(ctx context.Context, dni string) (*entities.Profile, error)
	ExistsByUserID(ctx context.Context, userID valueobjects.UserID) (bool, error)
	ExistsByDNI(ctx context.Context, dni string) (bool, error)
}
