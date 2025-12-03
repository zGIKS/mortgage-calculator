package resources

import "time"

type UserResource struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string    `json:"email" example:"user@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

type RegisterUserResource struct {
	DNI      string `json:"dni" example:"12345678" validate:"required,len=8,numeric"`
	Email    string `json:"email" example:"user@example.com" validate:"required,email"`
	Password string `json:"password" example:"mypassword123" validate:"required,min=6"`
}

type LoginResource struct {
	Email    string `json:"email" example:"user@example.com" validate:"required,email"`
	Password string `json:"password" example:"mypassword123" validate:"required"`
}

type LoginResponseResource struct {
	Token string       `json:"token" example:"user-1-token"`
	User  UserResource `json:"user"`
}

type UpdateUserResource struct {
	Password *string `json:"password" example:"newpassword123" validate:"required,min=6"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
