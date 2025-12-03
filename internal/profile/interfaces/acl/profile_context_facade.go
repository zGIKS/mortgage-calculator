package acl

import "context"

// ProfileContextFacade define el contrato ACL para que otros bounded contexts consulten Profile
// Este facade expone solo las operaciones necesarias para consultas de perfil
type ProfileContextFacade interface {
	// ExistsByDNI verifica si existe un perfil con el DNI especificado
	ExistsByDNI(ctx context.Context, dni string) (bool, error)

	// FindUserIDByDNI obtiene el UserID asociado a un DNI
	FindUserIDByDNI(ctx context.Context, dni string) (string, error)
}
