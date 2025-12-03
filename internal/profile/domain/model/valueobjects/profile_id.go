package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type ProfileID struct {
	value uuid.UUID
}

func NewProfileID(value uuid.UUID) (ProfileID, error) {
	if value == uuid.Nil {
		return ProfileID{}, errors.New("profile ID cannot be nil")
	}
	return ProfileID{value: value}, nil
}

func NewProfileIDFromString(value string) (ProfileID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return ProfileID{}, errors.New("invalid profile ID format")
	}
	return NewProfileID(id)
}

func (p ProfileID) Value() uuid.UUID {
	return p.value
}

func (p ProfileID) String() string {
	return p.value.String()
}

func (p ProfileID) IsZero() bool {
	return p.value == uuid.Nil
}
