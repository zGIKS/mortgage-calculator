package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID `json:"value" gorm:"type:uuid;column:user_id"`
}

func NewUserID(value uuid.UUID) (UserID, error) {
	if value == uuid.Nil {
		return UserID{}, errors.New("user ID cannot be nil")
	}
	return UserID{value: value}, nil
}

func NewUserIDFromString(value string) (UserID, error) {
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, errors.New("invalid UUID format")
	}
	return NewUserID(parsedUUID)
}

func GenerateUserID() UserID {
	return UserID{value: uuid.New()}
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
