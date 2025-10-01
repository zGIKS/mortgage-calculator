package queryservices

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/queries"
	"finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/domain/services"
)

type MortgageQueryServiceImpl struct {
	repository repositories.MortgageRepository
}

func NewMortgageQueryService(repository repositories.MortgageRepository) services.MortgageQueryService {
	return &MortgageQueryServiceImpl{
		repository: repository,
	}
}

func (s *MortgageQueryServiceImpl) HandleGetByID(
	ctx context.Context,
	query *queries.GetMortgageByIDQuery,
) (*entities.Mortgage, error) {
	return s.repository.FindByID(ctx, query.MortgageID)
}

func (s *MortgageQueryServiceImpl) HandleGetHistory(
	ctx context.Context,
	query *queries.GetMortgageHistoryQuery,
) ([]*entities.Mortgage, error) {
	return s.repository.FindByUserID(ctx, query.UserID, query.Limit, query.Offset)
}
