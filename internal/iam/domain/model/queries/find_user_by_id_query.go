package queries

import (
	"errors"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
)

type FindUserByIDQuery struct {
	userID valueobjects.UserID
}

func NewFindUserByIDQuery(userID valueobjects.UserID) (FindUserByIDQuery, error) {
	if userID.IsZero() {
		return FindUserByIDQuery{}, errors.New("user ID cannot be zero")
	}
	return FindUserByIDQuery{userID: userID}, nil
}

func (q FindUserByIDQuery) UserID() valueobjects.UserID { return q.userID }
