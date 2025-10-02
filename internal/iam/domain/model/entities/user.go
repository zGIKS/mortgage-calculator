package entities

import (
	"finanzas-backend/internal/iam/domain/model/valueobjects"
	"time"
)

type User struct {
	id        valueobjects.UserID
	email     valueobjects.Email
	password  valueobjects.Password
	fullName  string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email valueobjects.Email, password valueobjects.Password, fullName string) (*User, error) {
	return &User{
		email:     email,
		password:  password,
		fullName:  fullName,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func ReconstructUser(id valueobjects.UserID, email valueobjects.Email, password valueobjects.Password, fullName string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		password:  password,
		fullName:  fullName,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() valueobjects.UserID           { return u.id }
func (u *User) Email() valueobjects.Email         { return u.email }
func (u *User) Password() valueobjects.Password   { return u.password }
func (u *User) FullName() string                  { return u.fullName }
func (u *User) CreatedAt() time.Time              { return u.createdAt }
func (u *User) UpdatedAt() time.Time              { return u.updatedAt }

func (u *User) SetID(id valueobjects.UserID) {
	u.id = id
}

func (u *User) VerifyPassword(plainPassword string) bool {
	return u.password.Matches(plainPassword)
}

func (u *User) UpdateEmail(email valueobjects.Email) {
	u.email = email
	u.updatedAt = time.Now()
}

func (u *User) UpdatePassword(password valueobjects.Password) {
	u.password = password
	u.updatedAt = time.Now()
}
