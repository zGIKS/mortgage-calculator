package queries

import (
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"

	"github.com/google/uuid"
)

type GetBankByIDQuery struct {
	BankID valueobjects.BankID
}

func NewGetBankByIDQuery(bankID uuid.UUID) (*GetBankByIDQuery, error) {
	id, err := valueobjects.NewBankID(bankID)
	if err != nil {
		return nil, err
	}
	return &GetBankByIDQuery{BankID: id}, nil
}
