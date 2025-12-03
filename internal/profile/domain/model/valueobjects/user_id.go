package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// UserID represents User ID from IAM context
type UserID struct {
	value uuid.UUID
}

func NewUserID(value uuid.UUID) (UserID, error) {
	if value == uuid.Nil {
		return UserID{}, errors.New("user ID cannot be nil")
	}
	return UserID{value: value}, nil
}

func NewUserIDFromString(value string) (UserID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, errors.New("invalid user ID format")
	}
	return NewUserID(id)
}

func (u UserID) Value() uuid.UUID {
	return u.value
}

func (u UserID) String() string {
	return u.value.String()
}

func (u UserID) IsZero() bool {
	return u.value == uuid.Nil
}
