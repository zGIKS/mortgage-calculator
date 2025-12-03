package queries

import (
	"finanzas-backend/internal/profile/domain/model/valueobjects"
)

type FindProfileByUserIDQuery struct {
	userID valueobjects.UserID
}

func NewFindProfileByUserIDQuery(userID string) (FindProfileByUserIDQuery, error) {
	id, err := valueobjects.NewUserIDFromString(userID)
	if err != nil {
		return FindProfileByUserIDQuery{}, err
	}
	return FindProfileByUserIDQuery{userID: id}, nil
}

func (q FindProfileByUserIDQuery) UserID() valueobjects.UserID {
	return q.userID
}
