package acl

import "context"

// IAMContextFacade define el contrato ACL para que otros bounded contexts consulten IAM
// Este facade expone solo las operaciones necesarias para validación de autenticación
type IAMContextFacade interface {
	// ValidateToken valida un token JWT y retorna el UserID como string si es válido
	ValidateToken(ctx context.Context, token string) (string, error)

	// GetUserEmailByID obtiene el email de un usuario por su ID (UUID string)
	GetUserEmailByID(ctx context.Context, userID string) (string, error)
}
