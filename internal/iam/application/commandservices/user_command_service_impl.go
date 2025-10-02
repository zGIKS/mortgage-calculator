package commandservices

import (
	"context"
	"errors"
	"finanzas-backend/internal/iam/domain/model/commands"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
	"finanzas-backend/internal/iam/domain/repositories"
	"finanzas-backend/internal/iam/domain/services"
)

type userCommandServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserCommandService(userRepo repositories.UserRepository) services.UserCommandService {
	return &userCommandServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userCommandServiceImpl) HandleRegister(ctx context.Context, cmd commands.RegisterUserCommand) (*valueobjects.UserID, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, cmd.Email())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Create value objects
	email, err := valueobjects.NewEmail(cmd.Email())
	if err != nil {
		return nil, err
	}

	password, err := valueobjects.NewPassword(cmd.Password())
	if err != nil {
		return nil, err
	}

	// Create entity
	user, err := entities.NewUser(email, password, cmd.FullName())
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	userID := user.ID()
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

	// Update email if provided
	if cmd.Email() != nil {
		// Check if new email is already in use by another user
		existingUser, err := s.userRepo.FindByEmail(ctx, *cmd.Email())
		if err != nil {
			return err
		}
		if existingUser != nil && existingUser.ID().Value() != user.ID().Value() {
			return errors.New("email already in use by another user")
		}

		email, err := valueobjects.NewEmail(*cmd.Email())
		if err != nil {
			return err
		}
		user.UpdateEmail(email)
	}

	// Update password if provided
	if cmd.Password() != nil {
		password, err := valueobjects.NewPassword(*cmd.Password())
		if err != nil {
			return err
		}
		user.UpdatePassword(password)
	}

	// Update in repository
	return s.userRepo.Update(ctx, user)
}
