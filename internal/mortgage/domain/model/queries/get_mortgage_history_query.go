package queries

import (
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type GetMortgageHistoryQuery struct {
	UserID valueobjects.UserID
	Limit  int
	Offset int
}

func NewGetMortgageHistoryQuery(userID uint64) (*GetMortgageHistoryQuery, error) {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}
	return &GetMortgageHistoryQuery{
		UserID: uid,
		Limit:  50,
		Offset: 0,
	}, nil
}

func (q *GetMortgageHistoryQuery) WithPagination(limit, offset int) *GetMortgageHistoryQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}
