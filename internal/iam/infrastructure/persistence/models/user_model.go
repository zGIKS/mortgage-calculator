package models

import (
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
	"time"
)

type UserModel struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Email        string    `gorm:"uniqueIndex;not null;column:email"`
	PasswordHash string    `gorm:"not null;column:password_hash"`
	FullName     string    `gorm:"not null;column:full_name"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) ToEntity() (*entities.User, error) {
	userID, err := valueobjects.NewUserID(m.ID)
	if err != nil {
		return nil, err
	}

	email, err := valueobjects.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	password := valueobjects.NewPasswordFromHash(m.PasswordHash)

	return entities.ReconstructUser(userID, email, password, m.FullName, m.CreatedAt, m.UpdatedAt), nil
}

func FromEntity(user *entities.User) *UserModel {
	return &UserModel{
		ID:           user.ID().Value(),
		Email:        user.Email().Value(),
		PasswordHash: user.Password().Hash(),
		FullName:     user.FullName(),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	}
}
