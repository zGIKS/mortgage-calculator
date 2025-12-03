package queryservices

import (
	"context"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/queries"
	"finanzas-backend/internal/profile/domain/repositories"
	"finanzas-backend/internal/profile/domain/services"
)

type profileQueryServiceImpl struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileQueryService(profileRepo repositories.ProfileRepository) services.ProfileQueryService {
	return &profileQueryServiceImpl{
		profileRepo: profileRepo,
	}
}

func (s *profileQueryServiceImpl) HandleFindByID(ctx context.Context, query queries.FindProfileByIDQuery) (*entities.Profile, error) {
	return s.profileRepo.FindByID(ctx, query.ProfileID())
}

func (s *profileQueryServiceImpl) HandleFindByUserID(ctx context.Context, query queries.FindProfileByUserIDQuery) (*entities.Profile, error) {
	return s.profileRepo.FindByUserID(ctx, query.UserID())
}

func (s *profileQueryServiceImpl) HandleFindByDNI(ctx context.Context, query queries.FindProfileByDNIQuery) (*entities.Profile, error) {
	return s.profileRepo.FindByDNI(ctx, query.DNI())
}
