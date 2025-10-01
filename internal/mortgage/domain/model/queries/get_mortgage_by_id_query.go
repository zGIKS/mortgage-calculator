package queries

import (
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
)

type GetMortgageByIDQuery struct {
	MortgageID valueobjects.MortgageID
}

func NewGetMortgageByIDQuery(mortgageID uint64) (*GetMortgageByIDQuery, error) {
	id, err := valueobjects.NewMortgageID(mortgageID)
	if err != nil {
		return nil, err
	}
	return &GetMortgageByIDQuery{MortgageID: id}, nil
}
