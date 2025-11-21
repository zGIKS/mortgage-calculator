package services

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/queries"
)

type BankQueryService interface {
	HandleGetAll(ctx context.Context, query *queries.GetAllBanksQuery) ([]*entities.Bank, error)
	HandleGetByID(ctx context.Context, query *queries.GetBankByIDQuery) (*entities.Bank, error)
}
