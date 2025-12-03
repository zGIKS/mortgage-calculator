package repositories

import (
	"context"
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/infrastructure/persistence/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MortgageRepositoryImpl struct {
	db *gorm.DB
}

func NewMortgageRepository(db *gorm.DB) repositories.MortgageRepository {
	return &MortgageRepositoryImpl{db: db}
}

func (r *MortgageRepositoryImpl) Save(ctx context.Context, mortgage *entities.Mortgage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Guardar mortgage
		mortgageModel := r.toModel(mortgage)
		if err := tx.Create(mortgageModel).Error; err != nil {
			return err
		}

		// Asignar ID generado
		id, err := valueobjects.NewMortgageID(mortgageModel.ID)
		if err != nil {
			return err
		}
		mortgage.SetID(id)

		// Guardar items del cronograma
		if mortgage.PaymentSchedule() != nil {
			scheduleItems := r.toScheduleItemModels(mortgageModel.ID, mortgageModel.UserID, mortgage.PaymentSchedule())
			if len(scheduleItems) > 0 {
				if err := tx.Create(&scheduleItems).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *MortgageRepositoryImpl) FindByID(ctx context.Context, id valueobjects.MortgageID) (*entities.Mortgage, error) {
	var model models.MortgageModel
	result := r.db.WithContext(ctx).
		Preload("PaymentScheduleItems").
		First(&model, id.Value())

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
		Preload("PaymentScheduleItems").
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

func (r *MortgageRepositoryImpl) Update(ctx context.Context, mortgage *entities.Mortgage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Actualizar mortgage
		mortgageModel := r.toModel(mortgage)
		result := tx.Model(&models.MortgageModel{}).
			Where("id = ?", mortgage.ID().Value()).
			Updates(mortgageModel)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("mortgage not found")
		}

		// Eliminar items antiguos del cronograma
		if err := tx.Where("mortgage_id = ?", mortgage.ID().Value()).
			Delete(&models.PaymentScheduleItemModel{}).Error; err != nil {
			return err
		}

		// Guardar nuevos items del cronograma
		if mortgage.PaymentSchedule() != nil {
			scheduleItems := r.toScheduleItemModels(mortgage.ID().Value(), mortgage.UserID().Value().String(), mortgage.PaymentSchedule())
			if len(scheduleItems) > 0 {
				if err := tx.Create(&scheduleItems).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *MortgageRepositoryImpl) Delete(ctx context.Context, id valueobjects.MortgageID) error {
	// El CASCADE en la FK eliminará automáticamente los items del cronograma
	result := r.db.WithContext(ctx).Delete(&models.MortgageModel{}, id.Value())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("mortgage not found")
	}

	return nil
}

func (r *MortgageRepositoryImpl) toModel(mortgage *entities.Mortgage) *models.MortgageModel {
	return &models.MortgageModel{
		ID:                   mortgage.ID().Value(),
		UserID:               mortgage.UserID().Value(),
		PropertyPrice:        mortgage.PropertyPrice(),
		DownPayment:          mortgage.DownPayment(),
		LoanAmount:           mortgage.LoanAmount(),
		BonoTechoPropio:      mortgage.BonoTechoPropio(),
		InterestRate:         mortgage.InterestRate(),
		RateType:             mortgage.RateType().String(),
		TermMonths:           mortgage.TermMonths(),
		TermYears:            mortgage.TermYears(),
		GracePeriodMonths:    mortgage.GracePeriodMonths(),
		GracePeriodType:      mortgage.GracePeriodType().String(),
		Currency:             mortgage.Currency().String(),
		PaymentFrequencyDays: mortgage.PaymentFrequencyDays(),
		DaysInYear:           mortgage.DaysInYear(),
		AdministrationFee:    mortgage.AdministrationFee(),
		Portes:               mortgage.Portes(),
		AdditionalCosts:      mortgage.AdditionalCosts(),
		LifeInsuranceRate:    mortgage.LifeInsuranceRate(),
		PropertyInsurance:    mortgage.PropertyInsuranceRate(),
		EvaluationFee:        mortgage.EvaluationFee(),
		DisbursementFee:      mortgage.DisbursementFee(),
		PrincipalFinanced:    mortgage.PrincipalFinanced(),
		PeriodicRate:         mortgage.PeriodicRate(),
		FixedInstallment:     mortgage.FixedInstallment(),
		TotalInterestPaid:    mortgage.TotalInterestPaid(),
		TotalPaid:            mortgage.TotalPaid(),
		TotalPaidWithFees:    mortgage.TotalPaidWithFees(),
		TotalCharges:         mortgage.TotalCharges(),
		TotalInsurance:       mortgage.TotalInsurance(),
		TotalAdmin:           mortgage.TotalAdmin(),
		NPV:                  mortgage.NPV(),
		IRR:                  mortgage.IRR(),
		FlowIRR:              mortgage.FlowIRR(),
		TCEA:                 mortgage.TCEA(),
		CreatedAt:            mortgage.CreatedAt(),
	}
}

func (r *MortgageRepositoryImpl) toScheduleItemModels(
	mortgageID uint64,
	userID interface{},
	schedule *entities.PaymentSchedule,
) []models.PaymentScheduleItemModel {
	items := schedule.GetItems()
	itemModels := make([]models.PaymentScheduleItemModel, 0, len(items))

	// Convertir userID según el tipo recibido
	var userUUID uuid.UUID
	switch v := userID.(type) {
	case uuid.UUID:
		userUUID = v
	case string:
		parsed, _ := uuid.Parse(v)
		userUUID = parsed
	}

	for _, item := range items {
		itemModels = append(itemModels, models.PaymentScheduleItemModel{
			MortgageID:        mortgageID,
			UserID:            userUUID,
			Period:            item.Period,
			YearNumber:        item.YearNumber,
			PeriodicRate:      item.PeriodicRateApplied,
			Installment:       item.Installment,
			TotalInstallment:  item.TotalInstallment,
			Interest:          item.Interest,
			Amortization:      item.Amortization,
			Administration:    item.AdministrationFee,
			Portes:            item.Portes,
			LifeInsurance:     item.LifeInsurance,
			PropertyInsurance: item.PropertyInsurance,
			AdditionalCosts:   item.AdditionalCosts,
			RemainingBalance:  item.RemainingBalance,
			IsGracePeriod:     item.IsGracePeriod,
			GraceType:         item.GraceType,
		})
	}

	return itemModels
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
		model.TermYears,
		model.GracePeriodMonths,
		gracePeriodType,
		currency,
		model.PaymentFrequencyDays,
		model.DaysInYear,
		model.AdministrationFee,
		model.Portes,
		model.AdditionalCosts,
		model.LifeInsuranceRate,
		model.PropertyInsurance,
		model.EvaluationFee,
		model.DisbursementFee,
		model.PrincipalFinanced,
		model.PeriodicRate,
		model.FixedInstallment,
		model.TotalInterestPaid,
		model.TotalPaid,
		model.TotalPaidWithFees,
		model.TotalCharges,
		model.TotalInsurance,
		model.TotalAdmin,
		model.NPV,
		model.IRR,
		model.FlowIRR,
		model.TCEA,
		model.CreatedAt,
	)

	// Reconstruir cronograma desde items
	if len(model.PaymentScheduleItems) > 0 {
		schedule := entities.NewPaymentSchedule()
		for _, itemModel := range model.PaymentScheduleItems {
			schedule.AddItem(entities.PaymentScheduleItem{
				Period:              itemModel.Period,
				YearNumber:          itemModel.YearNumber,
				PeriodicRateApplied: itemModel.PeriodicRate,
				Installment:         itemModel.Installment,
				TotalInstallment:    itemModel.TotalInstallment,
				Interest:            itemModel.Interest,
				Amortization:        itemModel.Amortization,
				AdministrationFee:   itemModel.Administration,
				Portes:              itemModel.Portes,
				LifeInsurance:       itemModel.LifeInsurance,
				PropertyInsurance:   itemModel.PropertyInsurance,
				AdditionalCosts:     itemModel.AdditionalCosts,
				RemainingBalance:    itemModel.RemainingBalance,
				IsGracePeriod:       itemModel.IsGracePeriod,
				GraceType:           itemModel.GraceType,
			})
		}
		mortgage.SetPaymentSchedule(schedule)
	}

	return mortgage, nil
}
