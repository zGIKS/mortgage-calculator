package acl

import (
	"context"

	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"finanzas-backend/internal/profile/domain/repositories"
	"finanzas-backend/internal/profile/interfaces/acl"
)

type profileContextFacadeImpl struct {
	profileRepo repositories.ProfileRepository
}

// NewProfileContextFacade crea una nueva instancia del facade ACL de Profile
func NewProfileContextFacade(
	profileRepo repositories.ProfileRepository,
) acl.ProfileContextFacade {
	return &profileContextFacadeImpl{
		profileRepo: profileRepo,
	}
}

// ExistsByDNI verifica si existe un perfil con el DNI especificado
func (f *profileContextFacadeImpl) ExistsByDNI(ctx context.Context, dni string) (bool, error) {
	return f.profileRepo.ExistsByDNI(ctx, dni)
}

// FindUserIDByDNI obtiene el UserID asociado a un DNI
func (f *profileContextFacadeImpl) FindUserIDByDNI(ctx context.Context, dni string) (string, error) {
	profile, err := f.profileRepo.FindByDNI(ctx, dni)
	if err != nil {
		return "", err
	}
	if profile == nil {
		return "", nil
	}

	return profile.UserID().String(), nil
}

// CreateProfile crea un perfil automáticamente con datos de RENIEC
func (f *profileContextFacadeImpl) CreateProfile(ctx context.Context, userID, dni, firstName, firstLastName, secondLastName string) error {
	// Create value objects
	userIDVO, err := valueobjects.NewUserIDFromString(userID)
	if err != nil {
		return err
	}

	dniVO, err := valueobjects.NewDNI(dni)
	if err != nil {
		return err
	}

	// Create profile with RENIEC data and empty optional fields
	profile, err := entities.NewProfile(
		userIDVO,
		dniVO,
		firstName,
		firstLastName,
		secondLastName,
		valueobjects.EmptyPhoneNumber(),
		valueobjects.EmptyMonthlyIncome(),
		valueobjects.EmptyMaritalStatus(),
		false, // isFirstHome
		false, // hasOwnLand
	)
	if err != nil {
		return err
	}

	// Save profile
	return f.profileRepo.Save(ctx, profile)
}
