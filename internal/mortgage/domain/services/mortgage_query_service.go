package services

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/queries"
)

type MortgageQueryService interface {
	HandleGetByID(ctx context.Context, query *queries.GetMortgageByIDQuery) (*entities.Mortgage, error)
	HandleGetHistory(ctx context.Context, query *queries.GetMortgageHistoryQuery) ([]*entities.Mortgage, error)
}
