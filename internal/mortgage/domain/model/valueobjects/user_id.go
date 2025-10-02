package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID
}

func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, errors.New("invalid UUID format")
	}
	return UserID{value: parsedUUID}, nil
}

func NewUserIDFromUUID(value uuid.UUID) (UserID, error) {
	if value == uuid.Nil {
		return UserID{}, errors.New("user ID cannot be nil")
	}
	return UserID{value: value}, nil
}

func (u UserID) Value() uuid.UUID {
	return u.value
}

func (u UserID) String() string {
	return u.value.String()
}
