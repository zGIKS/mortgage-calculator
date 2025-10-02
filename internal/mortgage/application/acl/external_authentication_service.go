package acl

import (
	"context"
	"errors"

	iam_acl "finanzas-backend/internal/iam/interfaces/acl"
)

// ExternalAuthenticationService - ACL implementation para acceder a IAM desde Mortgage
// Este servicio actúa como Anti-Corruption Layer entre Mortgage e IAM
type ExternalAuthenticationService struct {
	iamFacade iam_acl.IAMContextFacade
}

func NewExternalAuthenticationService(iamFacade iam_acl.IAMContextFacade) *ExternalAuthenticationService {
	return &ExternalAuthenticationService{
		iamFacade: iamFacade,
	}
}

// ValidateTokenAndGetUserID valida un token JWT y retorna el UserID de Mortgage como string
func (s *ExternalAuthenticationService) ValidateTokenAndGetUserID(ctx context.Context, token string) (string, error) {
	// Llamar al facade de IAM para validar el token (retorna UUID string)
	userIDString, err := s.iamFacade.ValidateToken(ctx, token)
	if err != nil {
		return "", errors.New("invalid or expired token")
	}

	if userIDString == "" {
		return "", errors.New("invalid user ID from token")
	}

	return userIDString, nil
}

// GetUserEmail obtiene el email de un usuario por su ID (UUID string)
func (s *ExternalAuthenticationService) GetUserEmail(ctx context.Context, userID string) (string, error) {
	return s.iamFacade.GetUserEmailByID(ctx, userID)
}
