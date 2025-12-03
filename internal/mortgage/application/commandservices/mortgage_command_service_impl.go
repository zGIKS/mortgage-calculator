package commandservices

import (
	"context"
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/commands"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"finanzas-backend/internal/mortgage/domain/repositories"
	"finanzas-backend/internal/mortgage/domain/services"
	"math"
)

type MortgageCommandServiceImpl struct {
	repository repositories.MortgageRepository
	calculator *services.FrenchMethodCalculator
}

func NewMortgageCommandService(
	repository repositories.MortgageRepository,
) services.MortgageCommandService {
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

	// Validate and create RateType
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
		cmd.TermYears,
		cmd.GracePeriodMonths,
		gracePeriodType,
		currency,
		cmd.AdministrationFee,
		cmd.Portes,
		cmd.AdditionalCosts,
		cmd.LifeInsuranceRate,
		cmd.PropertyInsurance,
		cmd.EvaluationFee,
		cmd.DisbursementFee,
	)
	if err != nil {
		return nil, err
	}

	// Set payment configuration from command
	mortgage.SetPaymentFrequencyDays(cmd.PaymentFrequencyDays)
	mortgage.SetDaysInYear(cmd.DaysInYear)

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

	// Calcular TIR de cuota base
	irr, err := s.calculator.CalculateIRR(mortgage)
	if err != nil {
		return nil, err
	}
	mortgage.SetIRR(irr)

	flowIRR, err := s.calculator.CalculateFlowIRR(mortgage)
	if err != nil {
		return nil, err
	}
	mortgage.SetFlowIRR(flowIRR)

	// Calcular TCEA con flujos completos
	tcea := s.calculator.CalculateTCEA(flowIRR, mortgage.PeriodsPerYear())
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

	// Valores base actuales
	propertyPrice := mortgage.PropertyPrice()
	downPayment := mortgage.DownPayment()
	loanAmount := mortgage.LoanAmount()
	bono := mortgage.BonoTechoPropio()
	interestRate := mortgage.InterestRate()
	rateType := mortgage.RateType()
	paymentFrequencyDays := mortgage.PaymentFrequencyDays()
	daysInYear := mortgage.DaysInYear()
	termMonths := mortgage.TermMonths()
	termYears := mortgage.TermYears()
	graceMonths := mortgage.GracePeriodMonths()
	graceType := mortgage.GracePeriodType()
	currency := mortgage.Currency()
	adminFee := mortgage.AdministrationFee()
	portes := mortgage.Portes()
	additionalCosts := mortgage.AdditionalCosts()
	lifeInsurance := mortgage.LifeInsuranceRate()
	propertyInsurance := mortgage.PropertyInsuranceRate()
	evaluationFee := mortgage.EvaluationFee()
	disbursementFee := mortgage.DisbursementFee()

	discountRate := valueOrDefault(cmd.NPVDiscountRate(), 0)
	needsRecalculation := false

	// Aplicar cambios del comando
	if cmd.PropertyPrice() != nil {
		propertyPrice = *cmd.PropertyPrice()
		needsRecalculation = true
	}
	if cmd.DownPayment() != nil {
		downPayment = *cmd.DownPayment()
		needsRecalculation = true
	}
	if cmd.LoanAmount() != nil {
		loanAmount = *cmd.LoanAmount()
		needsRecalculation = true
	}
	if cmd.BonoTechoPropio() != nil {
		bono = *cmd.BonoTechoPropio()
		needsRecalculation = true
	}
	if cmd.InterestRate() != nil {
		interestRate = *cmd.InterestRate()
		needsRecalculation = true
	}
	if cmd.RateType() != nil {
		newRate, err := valueobjects.NewRateType(*cmd.RateType())
		if err != nil {
			return nil, err
		}
		rateType = newRate
		needsRecalculation = true
	}
	if cmd.PaymentFrequencyDays() != nil {
		paymentFrequencyDays = *cmd.PaymentFrequencyDays()
		needsRecalculation = true
	}
	if cmd.DaysInYear() != nil {
		daysInYear = *cmd.DaysInYear()
		needsRecalculation = true
	}
	if cmd.TermMonths() != nil {
		termMonths = *cmd.TermMonths()
		needsRecalculation = true
	}
	if cmd.TermYears() != nil {
		termYears = *cmd.TermYears()
		needsRecalculation = true
	}
	if cmd.GracePeriodMonths() != nil {
		graceMonths = *cmd.GracePeriodMonths()
		needsRecalculation = true
	}
	if cmd.GracePeriodType() != nil {
		newGrace, err := valueobjects.NewGracePeriodType(*cmd.GracePeriodType())
		if err != nil {
			return nil, err
		}
		graceType = newGrace
		needsRecalculation = true
	}
	if cmd.Currency() != nil {
		// Validación: Si cambia moneda, DEBE actualizar todos los montos
		if cmd.PropertyPrice() == nil || cmd.DownPayment() == nil || cmd.LoanAmount() == nil {
			return nil, errors.New("when changing currency, you must update all monetary amounts (property_price, down_payment, loan_amount)")
		}
		newCurrency, err := valueobjects.NewCurrency(*cmd.Currency())
		if err != nil {
			return nil, err
		}
		currency = newCurrency
		needsRecalculation = true
	}
	if cmd.AdministrationFee() != nil {
		adminFee = *cmd.AdministrationFee()
		needsRecalculation = true
	}
	if cmd.Portes() != nil {
		portes = *cmd.Portes()
		needsRecalculation = true
	}
	if cmd.AdditionalCosts() != nil {
		additionalCosts = *cmd.AdditionalCosts()
		needsRecalculation = true
	}
	if cmd.LifeInsuranceRate() != nil {
		lifeInsurance = *cmd.LifeInsuranceRate()
		needsRecalculation = true
	}
	if cmd.PropertyInsurance() != nil {
		propertyInsurance = *cmd.PropertyInsurance()
		needsRecalculation = true
	}
	if cmd.EvaluationFee() != nil {
		evaluationFee = *cmd.EvaluationFee()
		needsRecalculation = true
	}
	if cmd.DisbursementFee() != nil {
		disbursementFee = *cmd.DisbursementFee()
		needsRecalculation = true
	}

	// Recalcular si corresponde
	if needsRecalculation {
		if termMonths <= 0 && termYears > 0 && paymentFrequencyDays > 0 && daysInYear > 0 {
			periodsPerYear := float64(daysInYear) / float64(paymentFrequencyDays)
			termMonths = int(math.Round(periodsPerYear * float64(termYears)))
		}

		calculated, err := entities.NewMortgage(
			mortgage.UserID(),
			propertyPrice,
			downPayment,
			loanAmount,
			bono,
			interestRate,
			rateType,
			termMonths,
			termYears,
			graceMonths,
			graceType,
			currency,
			adminFee,
			portes,
			additionalCosts,
			lifeInsurance,
			propertyInsurance,
			evaluationFee,
			disbursementFee,
		)
		if err != nil {
			return nil, err
		}

		calculated.SetPaymentFrequencyDays(paymentFrequencyDays)
		calculated.SetDaysInYear(daysInYear)

		if err := s.calculator.Calculate(calculated); err != nil {
			return nil, err
		}

		// Calcular VAN si se proporciona tasa de descuento
		if discountRate > 0 {
			npv, err := s.calculator.CalculateNPV(calculated, discountRate)
			if err != nil {
				return nil, err
			}
			calculated.SetNPV(npv)
		}

		// Calcular TIR para cuota base y flujos completos
		irr, err := s.calculator.CalculateIRR(calculated)
		if err != nil {
			return nil, err
		}
		calculated.SetIRR(irr)

		flowIRR, err := s.calculator.CalculateFlowIRR(calculated)
		if err != nil {
			return nil, err
		}
		calculated.SetFlowIRR(flowIRR)

		// Calcular TCEA con flujos completos
		tcea := s.calculator.CalculateTCEA(flowIRR, calculated.PeriodsPerYear())
		calculated.SetTCEA(tcea)

		// Reconstruir conservando metadata original
		mortgage = entities.ReconstructMortgage(
			mortgage.ID(),
			mortgage.UserID(),
			propertyPrice,
			downPayment,
			loanAmount,
			bono,
			interestRate,
			rateType,
			termMonths,
			termYears,
			graceMonths,
			graceType,
			currency,
			paymentFrequencyDays,
			daysInYear,
			adminFee,
			portes,
			additionalCosts,
			lifeInsurance,
			propertyInsurance,
			evaluationFee,
			disbursementFee,
			calculated.PrincipalFinanced(),
			calculated.PeriodicRate(),
			calculated.FixedInstallment(),
			calculated.TotalInterestPaid(),
			calculated.TotalPaid(),
			calculated.TotalPaidWithFees(),
			calculated.TotalCharges(),
			calculated.TotalInsurance(),
			calculated.TotalAdmin(),
			calculated.NPV(),
			calculated.IRR(),
			calculated.FlowIRR(),
			calculated.TCEA(),
			mortgage.CreatedAt(),
		)
		mortgage.SetPaymentSchedule(calculated.PaymentSchedule())
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
