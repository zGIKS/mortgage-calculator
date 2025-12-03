package services

import (
	"context"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/queries"
)

type ProfileQueryService interface {
	HandleFindByID(ctx context.Context, query queries.FindProfileByIDQuery) (*entities.Profile, error)
	HandleFindByUserID(ctx context.Context, query queries.FindProfileByUserIDQuery) (*entities.Profile, error)
	HandleFindByDNI(ctx context.Context, query queries.FindProfileByDNIQuery) (*entities.Profile, error)
}
