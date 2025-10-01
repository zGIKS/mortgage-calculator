package valueobjects

import "errors"

type UserID struct {
	value uint64
}

func NewUserID(value uint64) (UserID, error) {
	if value == 0 {
		return UserID{}, errors.New("user ID cannot be zero")
	}
	return UserID{value: value}, nil
}

func (u UserID) Value() uint64 {
	return u.value
}
