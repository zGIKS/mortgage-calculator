package services

import (
	"context"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/queries"
)

type UserQueryService interface {
	HandleFindByEmail(ctx context.Context, query queries.FindUserByEmailQuery) (*entities.User, error)
	HandleFindByID(ctx context.Context, query queries.FindUserByIDQuery) (*entities.User, error)
}
