package repositories

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type MortgageRepository interface {
	Save(ctx context.Context, mortgage *entities.Mortgage) error
	Update(ctx context.Context, mortgage *entities.Mortgage) error
	Delete(ctx context.Context, id valueobjects.MortgageID) error
	FindByID(ctx context.Context, id valueobjects.MortgageID) (*entities.Mortgage, error)
	FindByUserID(ctx context.Context, userID valueobjects.UserID, limit, offset int) ([]*entities.Mortgage, error)
}
