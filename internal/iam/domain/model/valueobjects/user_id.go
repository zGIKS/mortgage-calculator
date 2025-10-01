package valueobjects

import (
	"errors"
	"fmt"
)

type UserID struct {
	value uint64 `json:"value" gorm:"column:user_id"`
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

func (u UserID) String() string {
	return fmt.Sprintf("%d", u.value)
}

func (u UserID) IsZero() bool {
	return u.value == 0
}
