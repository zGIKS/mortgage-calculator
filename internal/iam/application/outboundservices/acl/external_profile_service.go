package acl

import (
	"context"
	profile_acl "finanzas-backend/internal/profile/interfaces/acl"
)

// ExternalProfileService - ACL implementation for accessing Profile from IAM context
type ExternalProfileService struct {
	profileFacade profile_acl.ProfileContextFacade
}

func NewExternalProfileService(profileFacade profile_acl.ProfileContextFacade) *ExternalProfileService {
	return &ExternalProfileService{
		profileFacade: profileFacade,
	}
}

// DNIExists verifies if a DNI is already registered in the Profile context
func (s *ExternalProfileService) DNIExists(ctx context.Context, dni string) (bool, error) {
	return s.profileFacade.ExistsByDNI(ctx, dni)
}

// CreateProfile creates a profile automatically with RENIEC data
func (s *ExternalProfileService) CreateProfile(ctx context.Context, userID, dni, firstName, firstLastName, secondLastName string) error {
	return s.profileFacade.CreateProfile(ctx, userID, dni, firstName, firstLastName, secondLastName)
}
