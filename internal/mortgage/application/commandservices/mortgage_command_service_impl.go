package commandservices

import (
	"context"
	"finanzas-backend/internal/mortgage/domain/model/commands"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/domain/services"
)

type MortgageCommandServiceImpl struct {
	repository repositories.MortgageRepository
	calculator *services.FrenchMethodCalculator
}

func NewMortgageCommandService(repository repositories.MortgageRepository) services.MortgageCommandService {
	return &MortgageCommandServiceImpl{
		repository: repository,
		calculator: services.NewFrenchMethodCalculator(),
	}
}

func (s *MortgageCommandServiceImpl) HandleCalculateMortgage(
	ctx context.Context,
	cmd *commands.CalculateMortgageCommand,
) (*entities.Mortgage, error) {
	// Crear value objects
	userID, err := valueobjects.NewUserID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	rateType, err := valueobjects.NewRateType(cmd.RateType)
	if err != nil {
		return nil, err
	}

	gracePeriodType, err := valueobjects.NewGracePeriodType(cmd.GracePeriodType)
	if err != nil {
		return nil, err
	}

	currency, err := valueobjects.NewCurrency(cmd.Currency)
	if err != nil {
		return nil, err
	}

	// Crear entidad Mortgage
	mortgage, err := entities.NewMortgage(
		userID,
		cmd.PropertyPrice,
		cmd.DownPayment,
		cmd.LoanAmount,
		cmd.BonoTechoPropio,
		cmd.InterestRate,
		rateType,
		cmd.TermMonths,
		cmd.GracePeriodMonths,
		gracePeriodType,
		currency,
	)
	if err != nil {
		return nil, err
	}

	// Calcular cronograma usando método francés
	if err := s.calculator.Calculate(mortgage); err != nil {
		return nil, err
	}

	// Calcular VAN si se proporciona tasa de descuento
	if cmd.NPVDiscountRate > 0 {
		npv, err := s.calculator.CalculateNPV(mortgage, cmd.NPVDiscountRate)
		if err != nil {
			return nil, err
		}
		mortgage.SetNPV(npv)
	}

	// Calcular TIR
	irr, err := s.calculator.CalculateIRR(mortgage)
	if err != nil {
		return nil, err
	}
	mortgage.SetIRR(irr)

	// Calcular TCEA
	tcea := s.calculator.CalculateTCEA(irr)
	mortgage.SetTCEA(tcea)

	// Guardar en repositorio
	if err := s.repository.Save(ctx, mortgage); err != nil {
		return nil, err
	}

	return mortgage, nil
}
