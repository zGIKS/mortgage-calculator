package repositories

import (
	"context"
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	domainrepos "finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/infrastructure/persistence/models"
	"gorm.io/gorm"
)

type BankRepositoryImpl struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) domainrepos.BankRepository {
	return &BankRepositoryImpl{db: db}
}

func (r *BankRepositoryImpl) FindByID(ctx context.Context, id valueobjects.BankID) (*entities.Bank, error) {
	var model models.BankModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bank not found")
		}
		return nil, err
	}
	return r.toDomain(&model)
}

func (r *BankRepositoryImpl) List(ctx context.Context) ([]*entities.Bank, error) {
	var modelsList []models.BankModel
	if err := r.db.WithContext(ctx).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Bank, 0, len(modelsList))
	for _, model := range modelsList {
		bank, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		result = append(result, bank)
	}
	return result, nil
}

func (r *BankRepositoryImpl) toDomain(model *models.BankModel) (*entities.Bank, error) {
	bankID, err := valueobjects.NewBankID(model.ID)
	if err != nil {
		return nil, err
	}

	rateType, err := valueobjects.NewRateType(model.RateType)
	if err != nil {
		return nil, err
	}

	return entities.ReconstructBank(
		bankID,
		model.Name,
		rateType,
		model.PaymentFrequencyDays,
		model.DaysInYear,
		model.IncludesInflation,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
