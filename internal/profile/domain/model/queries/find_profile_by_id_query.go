package queries

import (
	"finanzas-backend/internal/profile/domain/model/valueobjects"
)

type FindProfileByIDQuery struct {
	profileID valueobjects.ProfileID
}

func NewFindProfileByIDQuery(profileID string) (FindProfileByIDQuery, error) {
	id, err := valueobjects.NewProfileIDFromString(profileID)
	if err != nil {
		return FindProfileByIDQuery{}, err
	}
	return FindProfileByIDQuery{profileID: id}, nil
}

func (q FindProfileByIDQuery) ProfileID() valueobjects.ProfileID {
	return q.profileID
}
