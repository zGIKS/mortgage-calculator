package repositories

import (
	"context"
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/infrastructure/persistence/models"
	"gorm.io/gorm"
)

type MortgageRepositoryImpl struct {
	db *gorm.DB
}

func NewMortgageRepository(db *gorm.DB) repositories.MortgageRepository {
	return &MortgageRepositoryImpl{db: db}
}

func (r *MortgageRepositoryImpl) Save(ctx context.Context, mortgage *entities.Mortgage) error {
	model := r.toModel(mortgage)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return result.Error
	}

	// Asignar ID generado
	id, err := valueobjects.NewMortgageID(model.ID)
	if err != nil {
		return err
	}
	mortgage.SetID(id)

	return nil
}

func (r *MortgageRepositoryImpl) FindByID(ctx context.Context, id valueobjects.MortgageID) (*entities.Mortgage, error) {
	var model models.MortgageModel
	result := r.db.WithContext(ctx).First(&model, id.Value())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("mortgage not found")
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

func (r *MortgageRepositoryImpl) FindByUserID(
	ctx context.Context,
	userID valueobjects.UserID,
	limit, offset int,
) ([]*entities.Mortgage, error) {
	var models []models.MortgageModel
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID.Value()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	mortgages := make([]*entities.Mortgage, 0, len(models))
	for _, model := range models {
		mortgage, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		mortgages = append(mortgages, mortgage)
	}

	return mortgages, nil
}

func (r *MortgageRepositoryImpl) toModel(mortgage *entities.Mortgage) *models.MortgageModel {
	var scheduleJSON models.PaymentScheduleJSON
	if mortgage.PaymentSchedule() != nil {
		scheduleJSON = mortgage.PaymentSchedule().GetItems()
	}

	return &models.MortgageModel{
		ID:                mortgage.ID().Value(),
		UserID:            mortgage.UserID().Value(),
		PropertyPrice:     mortgage.PropertyPrice(),
		DownPayment:       mortgage.DownPayment(),
		LoanAmount:        mortgage.LoanAmount(),
		BonoTechoPropio:   mortgage.BonoTechoPropio(),
		InterestRate:      mortgage.InterestRate(),
		RateType:          mortgage.RateType().String(),
		TermMonths:        mortgage.TermMonths(),
		GracePeriodMonths: mortgage.GracePeriodMonths(),
		GracePeriodType:   mortgage.GracePeriodType().String(),
		Currency:          mortgage.Currency().String(),
		PrincipalFinanced: mortgage.PrincipalFinanced(),
		PeriodicRate:      mortgage.PeriodicRate(),
		FixedInstallment:  mortgage.FixedInstallment(),
		PaymentSchedule:   scheduleJSON,
		TotalInterestPaid: mortgage.TotalInterestPaid(),
		TotalPaid:         mortgage.TotalPaid(),
		NPV:               mortgage.NPV(),
		IRR:               mortgage.IRR(),
		TCEA:              mortgage.TCEA(),
		CreatedAt:         mortgage.CreatedAt(),
	}
}

func (r *MortgageRepositoryImpl) toDomain(model *models.MortgageModel) (*entities.Mortgage, error) {
	id, err := valueobjects.NewMortgageID(model.ID)
	if err != nil {
		return nil, err
	}

	userID, err := valueobjects.NewUserIDFromUUID(model.UserID)
	if err != nil {
		return nil, err
	}

	rateType, err := valueobjects.NewRateType(model.RateType)
	if err != nil {
		return nil, err
	}

	gracePeriodType, err := valueobjects.NewGracePeriodType(model.GracePeriodType)
	if err != nil {
		return nil, err
	}

	currency, err := valueobjects.NewCurrency(model.Currency)
	if err != nil {
		return nil, err
	}

	mortgage := entities.ReconstructMortgage(
		id,
		userID,
		model.PropertyPrice,
		model.DownPayment,
		model.LoanAmount,
		model.BonoTechoPropio,
		model.InterestRate,
		rateType,
		model.TermMonths,
		model.GracePeriodMonths,
		gracePeriodType,
		currency,
		model.PrincipalFinanced,
		model.PeriodicRate,
		model.FixedInstallment,
		model.TotalInterestPaid,
		model.TotalPaid,
		model.NPV,
		model.IRR,
		model.TCEA,
		model.CreatedAt,
	)

	// Reconstruir cronograma
	if len(model.PaymentSchedule) > 0 {
		schedule := entities.NewPaymentSchedule()
		for _, item := range model.PaymentSchedule {
			schedule.AddItem(item)
		}
		mortgage.SetPaymentSchedule(schedule)
	}

	return mortgage, nil
}
