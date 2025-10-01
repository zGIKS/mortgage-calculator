package resources

import "time"

type UserResource struct {
	ID        uint64    `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	FullName  string    `json:"full_name" example:"John Doe"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

type RegisterUserResource struct {
	Email    string `json:"email" example:"user@example.com" validate:"required,email"`
	Password string `json:"password" example:"mypassword123" validate:"required,min=6"`
	FullName string `json:"full_name" example:"John Doe" validate:"required,min=2"`
}

type LoginResource struct {
	Email    string `json:"email" example:"user@example.com" validate:"required,email"`
	Password string `json:"password" example:"mypassword123" validate:"required"`
}

type LoginResponseResource struct {
	Token string       `json:"token" example:"user-1-token"`
	User  UserResource `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
