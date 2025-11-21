package queryservices

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/queries"
	"finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/domain/services"
)

type bankQueryServiceImpl struct {
	bankRepository repositories.BankRepository
}

func NewBankQueryService(bankRepository repositories.BankRepository) services.BankQueryService {
	return &bankQueryServiceImpl{
		bankRepository: bankRepository,
	}
}

func (s *bankQueryServiceImpl) HandleGetAll(ctx context.Context, query *queries.GetAllBanksQuery) ([]*entities.Bank, error) {
	return s.bankRepository.List(ctx)
}

func (s *bankQueryServiceImpl) HandleGetByID(ctx context.Context, query *queries.GetBankByIDQuery) (*entities.Bank, error) {
	return s.bankRepository.FindByID(ctx, query.BankID)
}
