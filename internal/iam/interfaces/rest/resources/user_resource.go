package resources

import "time"

type UserResource struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
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

type UpdateUserResource struct {
	Email    *string `json:"email,omitempty" example:"newemail@example.com" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" example:"newpassword123" validate:"omitempty,min=6"`
	FullName *string `json:"full_name,omitempty" example:"Jane Doe" validate:"omitempty,min=2"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
