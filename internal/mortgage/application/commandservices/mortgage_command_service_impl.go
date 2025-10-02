package commandservices

import (
	"context"
	"errors"
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

func (s *MortgageCommandServiceImpl) HandleUpdateMortgage(
	ctx context.Context,
	cmd *commands.UpdateMortgageCommand,
) (*entities.Mortgage, error) {
	// Buscar hipoteca existente
	mortgage, err := s.repository.FindByID(ctx, cmd.MortgageID())
	if err != nil {
		return nil, err
	}

	// Variables para recálculo
	needsRecalculation := false

	// Actualizar campos si se proporcionan
	if cmd.PropertyPrice() != nil {
		mortgage = entities.ReconstructMortgage(
			mortgage.ID(),
			mortgage.UserID(),
			*cmd.PropertyPrice(),
			mortgage.DownPayment(),
			mortgage.LoanAmount(),
			mortgage.BonoTechoPropio(),
			mortgage.InterestRate(),
			mortgage.RateType(),
			mortgage.TermMonths(),
			mortgage.GracePeriodMonths(),
			mortgage.GracePeriodType(),
			mortgage.Currency(),
			mortgage.PrincipalFinanced(),
			mortgage.PeriodicRate(),
			mortgage.FixedInstallment(),
			mortgage.TotalInterestPaid(),
			mortgage.TotalPaid(),
			mortgage.NPV(),
			mortgage.IRR(),
			mortgage.TCEA(),
			mortgage.CreatedAt(),
		)
	}

	if cmd.DownPayment() != nil {
		needsRecalculation = true
	}

	if cmd.LoanAmount() != nil {
		needsRecalculation = true
	}

	if cmd.BonoTechoPropio() != nil {
		needsRecalculation = true
	}

	if cmd.InterestRate() != nil {
		needsRecalculation = true
	}

	if cmd.RateType() != nil {
		rateType, err := valueobjects.NewRateType(*cmd.RateType())
		if err != nil {
			return nil, err
		}
		mortgage = entities.ReconstructMortgage(
			mortgage.ID(),
			mortgage.UserID(),
			mortgage.PropertyPrice(),
			mortgage.DownPayment(),
			mortgage.LoanAmount(),
			mortgage.BonoTechoPropio(),
			mortgage.InterestRate(),
			rateType,
			mortgage.TermMonths(),
			mortgage.GracePeriodMonths(),
			mortgage.GracePeriodType(),
			mortgage.Currency(),
			mortgage.PrincipalFinanced(),
			mortgage.PeriodicRate(),
			mortgage.FixedInstallment(),
			mortgage.TotalInterestPaid(),
			mortgage.TotalPaid(),
			mortgage.NPV(),
			mortgage.IRR(),
			mortgage.TCEA(),
			mortgage.CreatedAt(),
		)
		needsRecalculation = true
	}

	if cmd.TermMonths() != nil {
		needsRecalculation = true
	}

	if cmd.GracePeriodMonths() != nil {
		needsRecalculation = true
	}

	if cmd.GracePeriodType() != nil {
		gracePeriodType, err := valueobjects.NewGracePeriodType(*cmd.GracePeriodType())
		if err != nil {
			return nil, err
		}
		mortgage = entities.ReconstructMortgage(
			mortgage.ID(),
			mortgage.UserID(),
			mortgage.PropertyPrice(),
			mortgage.DownPayment(),
			mortgage.LoanAmount(),
			mortgage.BonoTechoPropio(),
			mortgage.InterestRate(),
			mortgage.RateType(),
			mortgage.TermMonths(),
			mortgage.GracePeriodMonths(),
			gracePeriodType,
			mortgage.Currency(),
			mortgage.PrincipalFinanced(),
			mortgage.PeriodicRate(),
			mortgage.FixedInstallment(),
			mortgage.TotalInterestPaid(),
			mortgage.TotalPaid(),
			mortgage.NPV(),
			mortgage.IRR(),
			mortgage.TCEA(),
			mortgage.CreatedAt(),
		)
		needsRecalculation = true
	}

	if cmd.Currency() != nil {
		currency, err := valueobjects.NewCurrency(*cmd.Currency())
		if err != nil {
			return nil, err
		}
		mortgage = entities.ReconstructMortgage(
			mortgage.ID(),
			mortgage.UserID(),
			mortgage.PropertyPrice(),
			mortgage.DownPayment(),
			mortgage.LoanAmount(),
			mortgage.BonoTechoPropio(),
			mortgage.InterestRate(),
			mortgage.RateType(),
			mortgage.TermMonths(),
			mortgage.GracePeriodMonths(),
			mortgage.GracePeriodType(),
			currency,
			mortgage.PrincipalFinanced(),
			mortgage.PeriodicRate(),
			mortgage.FixedInstallment(),
			mortgage.TotalInterestPaid(),
			mortgage.TotalPaid(),
			mortgage.NPV(),
			mortgage.IRR(),
			mortgage.TCEA(),
			mortgage.CreatedAt(),
		)
	}

	// Recalcular si es necesario
	if needsRecalculation {
		// Crear nueva entidad con valores actualizados para recalcular
		updatedMortgage, err := entities.NewMortgage(
			mortgage.UserID(),
			valueOrDefault(cmd.PropertyPrice(), mortgage.PropertyPrice()),
			valueOrDefault(cmd.DownPayment(), mortgage.DownPayment()),
			valueOrDefault(cmd.LoanAmount(), mortgage.LoanAmount()),
			valueOrDefault(cmd.BonoTechoPropio(), mortgage.BonoTechoPropio()),
			valueOrDefault(cmd.InterestRate(), mortgage.InterestRate()),
			mortgage.RateType(),
			valueOrDefaultInt(cmd.TermMonths(), mortgage.TermMonths()),
			valueOrDefaultInt(cmd.GracePeriodMonths(), mortgage.GracePeriodMonths()),
			mortgage.GracePeriodType(),
			mortgage.Currency(),
		)
		if err != nil {
			return nil, err
		}

		// Mantener el ID original
		updatedMortgage.SetID(mortgage.ID())

		// Recalcular
		if err := s.calculator.Calculate(updatedMortgage); err != nil {
			return nil, err
		}

		// Calcular VAN si se proporciona tasa de descuento
		discountRate := valueOrDefault(cmd.NPVDiscountRate(), 0)
		if discountRate > 0 {
			npv, err := s.calculator.CalculateNPV(updatedMortgage, discountRate)
			if err != nil {
				return nil, err
			}
			updatedMortgage.SetNPV(npv)
		}

		// Calcular TIR
		irr, err := s.calculator.CalculateIRR(updatedMortgage)
		if err != nil {
			return nil, err
		}
		updatedMortgage.SetIRR(irr)

		// Calcular TCEA
		tcea := s.calculator.CalculateTCEA(irr)
		updatedMortgage.SetTCEA(tcea)

		mortgage = updatedMortgage
	}

	// Actualizar en repositorio
	if err := s.repository.Update(ctx, mortgage); err != nil {
		return nil, err
	}

	return mortgage, nil
}

func (s *MortgageCommandServiceImpl) HandleDeleteMortgage(
	ctx context.Context,
	cmd *commands.DeleteMortgageCommand,
) error {
	// Verificar que la hipoteca existe y pertenece al usuario
	mortgage, err := s.repository.FindByID(ctx, cmd.MortgageID())
	if err != nil {
		return err
	}

	// Verificar que pertenece al usuario
	if mortgage.UserID().String() != cmd.UserID().String() {
		return ErrUnauthorizedAccess
	}

	// Eliminar
	return s.repository.Delete(ctx, cmd.MortgageID())
}

// Helper functions
func valueOrDefault(ptr *float64, def float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return def
}

func valueOrDefaultInt(ptr *int, def int) int {
	if ptr != nil {
		return *ptr
	}
	return def
}

var ErrUnauthorizedAccess = errors.New("unauthorized access to mortgage")
