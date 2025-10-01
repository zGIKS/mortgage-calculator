package acl

import "context"

// IAMContextFacade define el contrato ACL para que otros bounded contexts consulten IAM
// Este facade expone solo las operaciones necesarias para validación de autenticación
type IAMContextFacade interface {
	// ValidateToken valida un token JWT y retorna el UserID si es válido
	ValidateToken(ctx context.Context, token string) (uint64, error)

	// GetUserEmailByID obtiene el email de un usuario por su ID
	GetUserEmailByID(ctx context.Context, userID uint64) (string, error)
}
