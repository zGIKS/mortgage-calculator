package commandservices

import (
	"context"
	"errors"
	"finanzas-backend/internal/iam/application/outboundservices/acl"
	"finanzas-backend/internal/iam/domain/model/commands"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
	"finanzas-backend/internal/iam/domain/repositories"
	"finanzas-backend/internal/iam/domain/services"
	"finanzas-backend/internal/iam/infrastructure/external"
)

type userCommandServiceImpl struct {
	userRepo              repositories.UserRepository
	reniecService         *external.ReniecService
	externalProfileService *acl.ExternalProfileService
}

func NewUserCommandService(
	userRepo repositories.UserRepository,
	reniecService *external.ReniecService,
	externalProfileService *acl.ExternalProfileService,
) services.UserCommandService {
	return &userCommandServiceImpl{
		userRepo:              userRepo,
		reniecService:         reniecService,
		externalProfileService: externalProfileService,
	}
}

func (s *userCommandServiceImpl) HandleRegister(ctx context.Context, cmd commands.RegisterUserCommand) (*valueobjects.UserID, error) {
	// Step 1: Validate DNI with RENIEC and get person data
	personData, err := s.reniecService.GetPersonData(ctx, cmd.DNI())
	if err != nil {
		return nil, errors.New("DNI validation failed: " + err.Error())
	}

	// Step 2: Check if DNI is already registered in Profile context
	existsByDNI, err := s.externalProfileService.DNIExists(ctx, cmd.DNI())
	if err != nil {
		return nil, err
	}
	if existsByDNI {
		return nil, errors.New("user with this DNI already exists")
	}

	// Step 3: Check if email is already registered
	exists, err := s.userRepo.ExistsByEmail(ctx, cmd.Email())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Step 4: Create value objects (email and password only, NO DNI)
	email, err := valueobjects.NewEmail(cmd.Email())
	if err != nil {
		return nil, err
	}

	password, err := valueobjects.NewPassword(cmd.Password())
	if err != nil {
		return nil, err
	}

	// Step 5: Create entity (without DNI)
	user, err := entities.NewUser(email, password)
	if err != nil {
		return nil, err
	}

	// Step 6: Persist (only email and password)
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// Step 7: Create Profile automatically with RENIEC data
	userID := user.ID()
	if err := s.externalProfileService.CreateProfile(
		ctx,
		userID.String(),
		cmd.DNI(),
		personData.FirstName,
		personData.FirstLastName,
		personData.SecondLastName,
	); err != nil {
		// If profile creation fails, we should consider rolling back the user
		// For now, we'll just return the error
		return nil, errors.New("failed to create profile: " + err.Error())
	}

	return &userID, nil
}

func (s *userCommandServiceImpl) HandleUpdate(ctx context.Context, cmd *commands.UpdateUserCommand) error {
	// Find existing user
	user, err := s.userRepo.FindByID(ctx, cmd.UserID())
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Update password
	password, err := valueobjects.NewPassword(*cmd.Password())
	if err != nil {
		return err
	}
	user.UpdatePassword(password)

	// Update in repository
	return s.userRepo.Update(ctx, user)
}
