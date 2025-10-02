package repositories

import (
	"context"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id valueobjects.UserID) (*entities.User, error)
	FindByIDValue(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
