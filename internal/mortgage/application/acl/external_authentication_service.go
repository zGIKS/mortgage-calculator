package acl

import (
	"context"
	"errors"

	iam_acl "finanzas-backend/internal/iam/interfaces/acl"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
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

// ValidateTokenAndGetUserID valida un token JWT y retorna el UserID de Mortgage
func (s *ExternalAuthenticationService) ValidateTokenAndGetUserID(ctx context.Context, token string) (*valueobjects.UserID, error) {
	// Llamar al facade de IAM para validar el token
	userIDValue, err := s.iamFacade.ValidateToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	if userIDValue == 0 {
		return nil, errors.New("invalid user ID from token")
	}

	// Convertir a value object del contexto Mortgage
	userID, err := valueobjects.NewUserID(userIDValue)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

// GetUserEmail obtiene el email de un usuario por su ID
func (s *ExternalAuthenticationService) GetUserEmail(ctx context.Context, userID valueobjects.UserID) (string, error) {
	return s.iamFacade.GetUserEmailByID(ctx, userID.Value())
}
