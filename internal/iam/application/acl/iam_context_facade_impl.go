package acl

import (
	"context"
	"errors"

	"finanzas-backend/internal/iam/domain/repositories"
	"finanzas-backend/internal/iam/infrastructure/security"
	"finanzas-backend/internal/iam/interfaces/acl"
)

type iamContextFacadeImpl struct {
	jwtService *security.JWTService
	userRepo   repositories.UserRepository
}

// NewIAMContextFacade crea una nueva instancia del facade ACL de IAM
func NewIAMContextFacade(
	jwtService *security.JWTService,
	userRepo repositories.UserRepository,
) acl.IAMContextFacade {
	return &iamContextFacadeImpl{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// ValidateToken valida un token JWT y retorna el UserID si es válido
func (f *iamContextFacadeImpl) ValidateToken(ctx context.Context, token string) (uint64, error) {
	claims, err := f.jwtService.ValidateToken(token)
	if err != nil {
		return 0, errors.New("invalid or expired token")
	}

	// Verificar que el usuario aún existe en la base de datos
	user, err := f.userRepo.FindByIDValue(ctx, claims.UserID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("user not found")
	}

	return claims.UserID, nil
}

// GetUserEmailByID obtiene el email de un usuario por su ID
func (f *iamContextFacadeImpl) GetUserEmailByID(ctx context.Context, userID uint64) (string, error) {
	user, err := f.userRepo.FindByIDValue(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	return user.Email().Value(), nil
}
