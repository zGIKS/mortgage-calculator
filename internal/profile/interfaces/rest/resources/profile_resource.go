package resources

import "time"

type ProfileResource struct {
	ID             string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID         string    `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	DNI            string    `json:"dni" example:"46027897"`
	FirstName      string    `json:"first_name" example:"ROXANA KARINA"`
	FirstLastName  string    `json:"first_last_name" example:"DELGADO"`
	SecondLastName string    `json:"second_last_name" example:"CUELLAR"`
	FullName       string    `json:"full_name" example:"DELGADO CUELLAR ROXANA KARINA"`
	PhoneNumber    string    `json:"phone_number" example:"987654321"`
	MonthlyIncome  float64   `json:"monthly_income" example:"5000.00"`
	Currency       string    `json:"currency" example:"PEN"`
	MaritalStatus  string    `json:"marital_status" example:"SOLTERO"`
	IsFirstHome    bool      `json:"is_first_home" example:"true"`
	HasOwnLand     bool      `json:"has_own_land" example:"false"`
	CreatedAt      time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

type CreateProfileResource struct {
	DNI           string  `json:"dni" example:"46027897" validate:"required,len=8"`
	PhoneNumber   string  `json:"phone_number" example:"987654321" validate:"required,len=9"`
	MonthlyIncome float64 `json:"monthly_income" example:"5000.00" validate:"required,min=0"`
	Currency      string  `json:"currency" example:"PEN" validate:"required,oneof=PEN USD"`
	MaritalStatus string  `json:"marital_status" example:"SOLTERO" validate:"required,oneof=SOLTERO CASADO DIVORCIADO VIUDO"`
	IsFirstHome   bool    `json:"is_first_home" example:"true"`
	HasOwnLand    bool    `json:"has_own_land" example:"false"`
}

type UpdateProfileResource struct {
	PhoneNumber   *string  `json:"phone_number,omitempty" example:"987654321" validate:"omitempty,len=9"`
	MonthlyIncome *float64 `json:"monthly_income,omitempty" example:"6000.00" validate:"omitempty,min=0"`
	Currency      *string  `json:"currency,omitempty" example:"USD" validate:"omitempty,oneof=PEN USD"`
	MaritalStatus *string  `json:"marital_status,omitempty" example:"CASADO" validate:"omitempty,oneof=SOLTERO CASADO DIVORCIADO VIUDO"`
	IsFirstHome   *bool    `json:"is_first_home,omitempty" example:"false"`
	HasOwnLand    *bool    `json:"has_own_land,omitempty" example:"true"`
}

type ReniecDataResource struct {
	DNI            string `json:"dni" example:"46027897"`
	FirstName      string `json:"first_name" example:"ROXANA KARINA"`
	FirstLastName  string `json:"first_last_name" example:"DELGADO"`
	SecondLastName string `json:"second_last_name" example:"CUELLAR"`
	FullName       string `json:"full_name" example:"DELGADO CUELLAR ROXANA KARINA"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
