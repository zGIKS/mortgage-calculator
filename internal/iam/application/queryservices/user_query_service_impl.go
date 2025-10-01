package queryservices

import (
	"context"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/queries"
	"finanzas-backend/internal/iam/domain/repositories"
	"finanzas-backend/internal/iam/domain/services"
)

type userQueryServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserQueryService(userRepo repositories.UserRepository) services.UserQueryService {
	return &userQueryServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userQueryServiceImpl) HandleFindByEmail(ctx context.Context, query queries.FindUserByEmailQuery) (*entities.User, error) {
	return s.userRepo.FindByEmail(ctx, query.Email())
}

func (s *userQueryServiceImpl) HandleFindByID(ctx context.Context, query queries.FindUserByIDQuery) (*entities.User, error) {
	return s.userRepo.FindByID(ctx, query.UserID())
}
