package queries

import "errors"

type FindUserByEmailQuery struct {
	email string
}

func NewFindUserByEmailQuery(email string) (FindUserByEmailQuery, error) {
	if email == "" {
		return FindUserByEmailQuery{}, errors.New("email cannot be empty")
	}
	return FindUserByEmailQuery{email: email}, nil
}

func (q FindUserByEmailQuery) Email() string { return q.email }
