package commandservices

import (
	"context"
	"errors"
	"finanzas-backend/internal/profile/domain/model/commands"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"finanzas-backend/internal/profile/domain/repositories"
	"finanzas-backend/internal/profile/domain/services"
	"finanzas-backend/internal/profile/infrastructure/external"
)

type profileCommandServiceImpl struct {
	profileRepo   repositories.ProfileRepository
	reniecService *external.ReniecService
}

func NewProfileCommandService(
	profileRepo repositories.ProfileRepository,
	reniecService *external.ReniecService,
) services.ProfileCommandService {
	return &profileCommandServiceImpl{
		profileRepo:   profileRepo,
		reniecService: reniecService,
	}
}

func (s *profileCommandServiceImpl) HandleCreate(ctx context.Context, cmd commands.CreateProfileCommand) (*valueobjects.ProfileID, error) {
	// Step 1: Check if profile already exists for this user
	exists, err := s.profileRepo.ExistsByUserID(ctx, cmd.UserID())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("profile already exists for this user")
	}

	// Step 2: Validate DNI with RENIEC
	personData, err := s.reniecService.GetPersonDataByDNI(ctx, cmd.DNI())
	if err != nil {
		return nil, errors.New("DNI validation failed: " + err.Error())
	}

	// Step 3: Check if DNI is already registered by another user
	existsByDNI, err := s.profileRepo.ExistsByDNI(ctx, cmd.DNI())
	if err != nil {
		return nil, err
	}
	if existsByDNI {
		return nil, errors.New("DNI already registered by another user")
	}

	// Step 4: Create value objects
	dni, err := valueobjects.NewDNI(cmd.DNI())
	if err != nil {
		return nil, err
	}

	phoneNumber, err := valueobjects.NewPhoneNumber(cmd.PhoneNumber())
	if err != nil {
		return nil, err
	}

	monthlyIncome, err := valueobjects.NewMonthlyIncome(cmd.MonthlyIncome(), valueobjects.Currency(cmd.Currency()))
	if err != nil {
		return nil, err
	}

	maritalStatus, err := valueobjects.NewMaritalStatus(cmd.MaritalStatus())
	if err != nil {
		return nil, err
	}

	// Step 5: Create entity with RENIEC data
	profile, err := entities.NewProfile(
		cmd.UserID(),
		dni,
		personData.FirstName,
		personData.FirstLastName,
		personData.SecondLastName,
		phoneNumber,
		monthlyIncome,
		maritalStatus,
		cmd.IsFirstHome(),
		cmd.HasOwnLand(),
	)
	if err != nil {
		return nil, err
	}

	// Step 6: Persist
	if err := s.profileRepo.Save(ctx, profile); err != nil {
		return nil, err
	}

	profileID := profile.ID()
	return &profileID, nil
}

func (s *profileCommandServiceImpl) HandleUpdate(ctx context.Context, cmd *commands.UpdateProfileCommand) error {
	// Find existing profile
	profile, err := s.profileRepo.FindByID(ctx, cmd.ProfileID())
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	// Update phone number if provided
	if cmd.PhoneNumber() != nil {
		phoneNumber, err := valueobjects.NewPhoneNumber(*cmd.PhoneNumber())
		if err != nil {
			return err
		}
		profile.UpdatePhoneNumber(phoneNumber)
	}

	// Update monthly income if provided
	if cmd.MonthlyIncome() != nil && cmd.Currency() != nil {
		monthlyIncome, err := valueobjects.NewMonthlyIncome(*cmd.MonthlyIncome(), valueobjects.Currency(*cmd.Currency()))
		if err != nil {
			return err
		}
		profile.UpdateMonthlyIncome(monthlyIncome)
	}

	// Update marital status if provided
	if cmd.MaritalStatus() != nil {
		maritalStatus, err := valueobjects.NewMaritalStatus(*cmd.MaritalStatus())
		if err != nil {
			return err
		}
		profile.UpdateMaritalStatus(maritalStatus)
	}

	// Update is first home if provided
	if cmd.IsFirstHome() != nil {
		profile.UpdateIsFirstHome(*cmd.IsFirstHome())
	}

	// Update has own land if provided
	if cmd.HasOwnLand() != nil {
		profile.UpdateHasOwnLand(*cmd.HasOwnLand())
	}

	// Update in repository
	return s.profileRepo.Update(ctx, profile)
}
