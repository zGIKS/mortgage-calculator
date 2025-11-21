package repositories

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

// BankRepository exposes data access for bank configurations.
type BankRepository interface {
	FindByID(ctx context.Context, id valueobjects.BankID) (*entities.Bank, error)
	List(ctx context.Context) ([]*entities.Bank, error)
}
