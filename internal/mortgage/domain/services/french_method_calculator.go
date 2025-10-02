package services

import (
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"math"
)

// FrenchMethodCalculator implementa el método francés vencido ordinario
type FrenchMethodCalculator struct{}

func NewFrenchMethodCalculator() *FrenchMethodCalculator {
	return &FrenchMethodCalculator{}
}

// Calculate calcula el cronograma de pagos usando método francés
func (fmc *FrenchMethodCalculator) Calculate(mortgage *entities.Mortgage) error {
	// 1. Calcular principal financiado (después de aplicar el bono)
	principalFinanced := mortgage.LoanAmount() - mortgage.BonoTechoPropio()
	if principalFinanced <= 0 {
		return errors.New("principal financed must be greater than zero")
	}
	mortgage.SetPrincipalFinanced(principalFinanced)

	// 2. Convertir tasa de interés a tasa efectiva mensual (periodicRate)
	periodicRate, err := fmc.convertToPeriodicRate(mortgage.InterestRate(), mortgage.RateType())
	if err != nil {
		return err
	}
	mortgage.SetPeriodicRate(periodicRate)

	// 3. Ajustar principal si hay gracia total (capitalización de intereses)
	adjustedPrincipal := principalFinanced
	gracePeriods := 0

	if mortgage.GracePeriodType() != valueobjects.GracePeriodNone && mortgage.GracePeriodMonths() > 0 {
		gracePeriods = mortgage.GracePeriodMonths()

		if mortgage.GracePeriodType() == valueobjects.GracePeriodTotal {
			// P_gracia = P * (1 + i)^n_gracia
			adjustedPrincipal = principalFinanced * math.Pow(1+periodicRate, float64(gracePeriods))
		}
	}

	// 4. Calcular cuota fija para periodos posteriores a la gracia
	// A = P * [i(1+i)^n] / [(1+i)^n - 1]
	normalPeriods := mortgage.TermMonths() - gracePeriods
	if normalPeriods <= 0 {
		return errors.New("term months must be greater than grace period months")
	}

	fixedInstallment := fmc.calculateFixedInstallment(adjustedPrincipal, periodicRate, normalPeriods)
	mortgage.SetFixedInstallment(fixedInstallment)

	// 5. Generar cronograma de pagos
	schedule, err := fmc.generatePaymentSchedule(mortgage, adjustedPrincipal, periodicRate)
	if err != nil {
		return err
	}
	mortgage.SetPaymentSchedule(schedule)

	// 6. Calcular totales
	mortgage.SetTotalInterestPaid(schedule.TotalInterestPaid())
	mortgage.SetTotalPaid(schedule.TotalPaid())

	return nil
}

// convertToPeriodicRate convierte TNA o TEA a tasa efectiva mensual
func (fmc *FrenchMethodCalculator) convertToPeriodicRate(annualRate float64, rateType valueobjects.RateType) (float64, error) {
	if annualRate < 0 {
		return 0, errors.New("interest rate cannot be negative")
	}

	// Convertir porcentaje a decimal si es necesario (ej: 12% -> 0.12)
	rate := annualRate
	if rate > 1 {
		rate = rate / 100.0
	}

	switch rateType {
	case valueobjects.RateTypeNominal:
		// TNA: i_mensual = TNA / 12
		return rate / 12.0, nil
	case valueobjects.RateTypeEffective:
		// TEA: i_mensual = (1 + TEA)^(1/12) - 1
		return math.Pow(1+rate, 1.0/12.0) - 1, nil
	default:
		return 0, errors.New("invalid rate type")
	}
}

// calculateFixedInstallment calcula la cuota fija usando la fórmula del método francés
func (fmc *FrenchMethodCalculator) calculateFixedInstallment(principal, periodicRate float64, periods int) float64 {
	if periodicRate == 0 {
		// Si la tasa es 0%, la cuota es simplemente el principal dividido entre periodos
		return principal / float64(periods)
	}

	// A = P * [i(1+i)^n] / [(1+i)^n - 1]
	factor := math.Pow(1+periodicRate, float64(periods))
	installment := principal * (periodicRate * factor) / (factor - 1)
	return installment
}

// generatePaymentSchedule genera el cronograma completo de pagos
func (fmc *FrenchMethodCalculator) generatePaymentSchedule(
	mortgage *entities.Mortgage,
	adjustedPrincipal float64,
	periodicRate float64,
) (*entities.PaymentSchedule, error) {
	schedule := entities.NewPaymentSchedule()
	balance := mortgage.PrincipalFinanced() // Saldo inicial (antes de gracia)

	totalMonths := mortgage.TermMonths()
	gracePeriods := 0
	if mortgage.GracePeriodType() != valueobjects.GracePeriodNone {
		gracePeriods = mortgage.GracePeriodMonths()
	}

	for period := 1; period <= totalMonths; period++ {
		var item entities.PaymentScheduleItem
		item.Period = period

		// Determinar si es periodo de gracia
		isGracePeriod := gracePeriods > 0 && period <= gracePeriods
		item.IsGracePeriod = isGracePeriod

		// Calcular interés del periodo: I_k = saldo * i
		interest := balance * periodicRate
		item.Interest = interest

		if isGracePeriod {
			// Periodo de gracia
			switch mortgage.GracePeriodType() {
			case valueobjects.GracePeriodTotal:
				// Gracia total: no se paga ni interés ni capital
				item.Installment = 0
				item.Amortization = 0
				// Los intereses se capitalizan (se suman al saldo)
				balance += interest
			case valueobjects.GracePeriodPartial:
				// Gracia parcial: solo se paga el interés
				item.Installment = interest
				item.Amortization = 0
				// El saldo no cambia
			case valueobjects.GracePeriodNone:
				// No debería llegar aquí, pero por seguridad tratamos como periodo normal
				item.Installment = mortgage.FixedInstallment()
				item.Amortization = item.Installment - interest
				balance -= item.Amortization
			}
		} else {
			// Periodo normal (después de gracia)
			item.Installment = mortgage.FixedInstallment()
			// Amortización: C_k = A - I_k
			item.Amortization = item.Installment - interest
			// Nuevo saldo: Saldo_k = Saldo_{k-1} - C_k
			balance -= item.Amortization
		}

		// Redondear para evitar errores de precisión
		if balance < 0.01 && balance > -0.01 {
			balance = 0
		}

		item.RemainingBalance = balance
		schedule.AddItem(item)
	}

	return schedule, nil
}

// CalculateNPV calcula el Valor Actual Neto (VAN)
// VAN = suma de [CF_k / (1 + j)^k] donde j es la tasa de descuento
func (fmc *FrenchMethodCalculator) CalculateNPV(mortgage *entities.Mortgage, discountRate float64) (float64, error) {
	if mortgage.PaymentSchedule() == nil {
		return 0, errors.New("payment schedule not calculated")
	}

	// Convertir tasa anual a mensual si es necesario
	monthlyDiscountRate := discountRate
	if discountRate > 1 {
		monthlyDiscountRate = discountRate / 100.0
	}
	monthlyDiscountRate = math.Pow(1+monthlyDiscountRate, 1.0/12.0) - 1

	// Flujo inicial: desembolso del préstamo (negativo)
	npv := -mortgage.PrincipalFinanced()

	// Sumar flujos de cada periodo (cuotas pagadas)
	items := mortgage.PaymentSchedule().GetItems()
	for _, item := range items {
		// CF_k / (1 + j)^k
		discountFactor := math.Pow(1+monthlyDiscountRate, float64(item.Period))
		npv += item.Installment / discountFactor
	}

	return npv, nil
}

// CalculateIRR calcula la Tasa Interna de Retorno (TIR) usando método de Newton-Raphson
// TIR es la tasa que hace VAN = 0
func (fmc *FrenchMethodCalculator) CalculateIRR(mortgage *entities.Mortgage) (float64, error) {
	if mortgage.PaymentSchedule() == nil {
		return 0, errors.New("payment schedule not calculated")
	}

	// Método de Newton-Raphson para encontrar TIR
	// Empezamos con una estimación inicial (tasa periódica del préstamo)
	irr := mortgage.PeriodicRate()
	tolerance := 0.0000001
	maxIterations := 1000

	items := mortgage.PaymentSchedule().GetItems()
	principal := mortgage.PrincipalFinanced()

	for i := 0; i < maxIterations; i++ {
		// Calcular VAN y su derivada
		npv := -principal
		npvDerivative := 0.0

		for _, item := range items {
			period := float64(item.Period)
			cashFlow := item.Installment

			factor := math.Pow(1+irr, period)
			npv += cashFlow / factor
			npvDerivative -= cashFlow * period / (factor * (1 + irr))
		}

		// Si VAN es suficientemente cercano a 0, hemos encontrado la TIR
		if math.Abs(npv) < tolerance {
			return irr, nil
		}

		// Actualizar estimación usando Newton-Raphson
		if npvDerivative == 0 {
			return 0, errors.New("cannot calculate IRR, derivative is zero")
		}
		irr = irr - npv/npvDerivative

		// Verificar que la TIR sea razonable
		if irr < -1 || irr > 10 {
			return 0, errors.New("IRR calculation diverged")
		}
	}

	return irr, nil
}

// CalculateTCEA calcula la Tasa de Costo Efectivo Anual
// TCEA = (1 + TIR_mensual)^12 - 1
func (fmc *FrenchMethodCalculator) CalculateTCEA(irr float64) float64 {
	return math.Pow(1+irr, 12) - 1
}
