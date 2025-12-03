package queries

import (
	"errors"
)

type FindProfileByDNIQuery struct {
	dni string
}

func NewFindProfileByDNIQuery(dni string) (FindProfileByDNIQuery, error) {
	if dni == "" {
		return FindProfileByDNIQuery{}, errors.New("DNI cannot be empty")
	}
	return FindProfileByDNIQuery{dni: dni}, nil
}

func (q FindProfileByDNIQuery) DNI() string {
	return q.dni
}
