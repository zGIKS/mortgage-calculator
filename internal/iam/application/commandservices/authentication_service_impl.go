package commandservices

import (
	"context"
	"errors"
	"finanzas-backend/internal/iam/domain/model/commands"
	"finanzas-backend/internal/iam/domain/repositories"
	"finanzas-backend/internal/iam/domain/services"
	"finanzas-backend/internal/iam/infrastructure/security"
)

type authenticationServiceImpl struct {
	userRepo   repositories.UserRepository
	jwtService *security.JWTService
}

func NewAuthenticationService(
	userRepo repositories.UserRepository,
	jwtService *security.JWTService,
) services.AuthenticationService {
	return &authenticationServiceImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authenticationServiceImpl) HandleLogin(ctx context.Context, cmd commands.LoginCommand) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, cmd.Email())
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Verify password
	if !user.VerifyPassword(cmd.Password()) {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID().Value(), user.Email().Value())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
