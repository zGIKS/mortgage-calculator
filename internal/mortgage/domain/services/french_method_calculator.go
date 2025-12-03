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

	periodsPerYear := mortgage.PeriodsPerYear()
	if periodsPerYear <= 0 {
		return errors.New("periods per year must be greater than zero")
	}

	totalPeriods := mortgage.TermMonths()
	if totalPeriods <= 0 && mortgage.TermYears() > 0 {
		totalPeriods = int(math.Round(periodsPerYear * float64(mortgage.TermYears())))
	}
	if totalPeriods <= 0 {
		return errors.New("term months must be greater than zero")
	}

	// 2. Convertir tasa de interés a tasa efectiva por periodo según la frecuencia
	periodicRate, err := fmc.convertToPeriodicRate(mortgage.InterestRate(), mortgage.RateType(), periodsPerYear)
	if err != nil {
		return err
	}
	mortgage.SetPeriodicRate(periodicRate)

	// 3. Ajustar principal si hay gracia total (capitalización de intereses)
	adjustedPrincipal := principalFinanced
	gracePeriods := 0

	if mortgage.GracePeriodType() != valueobjects.GracePeriodNone && mortgage.GracePeriodMonths() > 0 {
		gracePeriods = mortgage.GracePeriodMonths()
		if gracePeriods > totalPeriods {
			return errors.New("grace period months must be less than total periods")
		}

		if mortgage.GracePeriodType() == valueobjects.GracePeriodTotal {
			// P_gracia = P * (1 + i)^n_gracia
			adjustedPrincipal = principalFinanced * math.Pow(1+periodicRate, float64(gracePeriods))
		}
	}

	// 4. Calcular cuota fija para periodos posteriores a la gracia
	// A = P * [i(1+i)^n] / [(1+i)^n - 1]
	normalPeriods := totalPeriods - gracePeriods
	if normalPeriods <= 0 {
		return errors.New("term months must be greater than grace period months")
	}

	fixedInstallment := fmc.calculateFixedInstallment(adjustedPrincipal, periodicRate, normalPeriods)
	mortgage.SetFixedInstallment(fixedInstallment)

	// 5. Generar cronograma de pagos con cargos adicionales
	lifeRate := normalizeRate(mortgage.LifeInsuranceRate())
	propertyRate := normalizeRate(mortgage.PropertyInsuranceRate())
	propertyInsurancePerPeriod := 0.0
	if propertyRate > 0 {
		propertyInsurancePerPeriod = mortgage.PropertyPrice() * propertyRate / periodsPerYear
	}

	schedule, err := fmc.generatePaymentSchedule(
		mortgage,
		periodicRate,
		totalPeriods,
		propertyInsurancePerPeriod,
		lifeRate,
		mortgage.AdministrationFee(),
		mortgage.Portes(),
		mortgage.AdditionalCosts(),
		periodsPerYear,
	)
	if err != nil {
		return err
	}
	mortgage.SetPaymentSchedule(schedule)

	// 6. Calcular totales
	mortgage.SetTotalInterestPaid(schedule.TotalInterestPaid())
	mortgage.SetTotalPaid(schedule.TotalPaid())
	mortgage.SetTotalPaidWithFees(schedule.TotalPaidWithCharges())
	mortgage.SetTotalCharges(schedule.TotalCharges())
	mortgage.SetTotalInsurance(schedule.TotalInsurance())
	mortgage.SetTotalAdmin(schedule.TotalAdminFees())

	return nil
}

// convertToPeriodicRate convierte TNA o TEA a tasa efectiva por periodo según la frecuencia indicada.
func (fmc *FrenchMethodCalculator) convertToPeriodicRate(
	annualRate float64,
	rateType valueobjects.RateType,
	periodsPerYear float64,
) (float64, error) {
	if annualRate < 0 {
		return 0, errors.New("interest rate cannot be negative")
	}
	if periodsPerYear <= 0 {
		return 0, errors.New("periods per year must be greater than zero")
	}

	// Convertir porcentaje a decimal si es necesario (ej: 12% -> 0.12)
	rate := annualRate
	if rate > 1 {
		rate = rate / 100.0
	}

	switch rateType {
	case valueobjects.RateTypeNominal:
		// TNA: i_periodo = TNA / m
		return rate / periodsPerYear, nil
	case valueobjects.RateTypeEffective:
		// TEA: i_periodo = (1 + TEA)^(1/m) - 1
		return math.Pow(1+rate, 1.0/periodsPerYear) - 1, nil
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
	periodicRate float64,
	totalPeriods int,
	propertyInsurancePerPeriod float64,
	lifeInsuranceRate float64,
	administrationFee float64,
	portes float64,
	additionalCosts float64,
	periodsPerYear float64,
) (*entities.PaymentSchedule, error) {
	schedule := entities.NewPaymentSchedule()
	balance := mortgage.PrincipalFinanced() // Saldo inicial (antes de gracia)

	gracePeriods := 0
	if mortgage.GracePeriodType() != valueobjects.GracePeriodNone {
		gracePeriods = mortgage.GracePeriodMonths()
	}

	for period := 1; period <= totalPeriods; period++ {
		var item entities.PaymentScheduleItem
		item.Period = period
		item.YearNumber = int(math.Ceil(float64(period) / periodsPerYear))
		item.PeriodicRateApplied = periodicRate
		item.GraceType = mortgage.GracePeriodType().String()

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

		// Seguros y gastos adicionales
		item.LifeInsurance = balance * lifeInsuranceRate
		item.PropertyInsurance = propertyInsurancePerPeriod
		item.AdministrationFee = administrationFee
		item.Portes = portes
		item.AdditionalCosts = additionalCosts
		item.TotalInstallment = item.Installment +
			item.LifeInsurance +
			item.PropertyInsurance +
			item.AdministrationFee +
			item.Portes +
			item.AdditionalCosts

		// Redondear para evitar errores de precisión
		if balance < 0.01 && balance > -0.01 {
			balance = 0
		}

		item.RemainingBalance = balance
		schedule.AddItem(item)
	}

	return schedule, nil
}

// CalculateNPV calcula el Valor Actual Neto (VAN) considerando todos los cargos
// VAN = suma de [CF_k / (1 + j)^k] donde j es la tasa de descuento
func (fmc *FrenchMethodCalculator) CalculateNPV(mortgage *entities.Mortgage, discountRate float64) (float64, error) {
	if mortgage.PaymentSchedule() == nil {
		return 0, errors.New("payment schedule not calculated")
	}

	flows := fmc.buildCashFlows(mortgage, true)
	if len(flows) == 0 {
		return 0, errors.New("no cash flows to evaluate")
	}

	periodicDiscountRate, err := fmc.convertToPeriodicRate(
		discountRate,
		valueobjects.RateTypeEffective,
		mortgage.PeriodsPerYear(),
	)
	if err != nil {
		return 0, err
	}

	npv := 0.0
	for idx, cf := range flows {
		npv += cf / math.Pow(1+periodicDiscountRate, float64(idx))
	}

	return npv, nil
}

// CalculateIRR calcula la Tasa Interna de Retorno (TIR) de la cuota base
func (fmc *FrenchMethodCalculator) CalculateIRR(mortgage *entities.Mortgage) (float64, error) {
	if mortgage.PaymentSchedule() == nil {
		return 0, errors.New("payment schedule not calculated")
	}

	flows := fmc.buildCashFlows(mortgage, false)
	return fmc.irrFromFlows(flows, mortgage.PeriodicRate())
}

// CalculateFlowIRR calcula la TIR considerando seguros, gastos y comisiones
func (fmc *FrenchMethodCalculator) CalculateFlowIRR(mortgage *entities.Mortgage) (float64, error) {
	if mortgage.PaymentSchedule() == nil {
		return 0, errors.New("payment schedule not calculated")
	}

	flows := fmc.buildCashFlows(mortgage, true)
	return fmc.irrFromFlows(flows, mortgage.PeriodicRate())
}

func (fmc *FrenchMethodCalculator) buildCashFlows(mortgage *entities.Mortgage, includeCharges bool) []float64 {
	if mortgage.PaymentSchedule() == nil {
		return nil
	}

	initial := mortgage.PrincipalFinanced()
	if includeCharges {
		initial = initial - mortgage.EvaluationFee() - mortgage.DisbursementFee()
	}

	flows := make([]float64, 0, len(mortgage.PaymentSchedule().GetItems())+1)
	flows = append(flows, initial)

	for _, item := range mortgage.PaymentSchedule().GetItems() {
		payment := item.Installment
		if includeCharges {
			payment = item.TotalInstallment
		}
		flows = append(flows, -payment)
	}

	return flows
}

func (fmc *FrenchMethodCalculator) irrFromFlows(flows []float64, guess float64) (float64, error) {
	if len(flows) == 0 {
		return 0, errors.New("no cash flows to evaluate")
	}

	irr := guess
	if irr == 0 {
		irr = 0.01
	}
	tolerance := 0.0000001
	maxIterations := 1000

	for i := 0; i < maxIterations; i++ {
		npv := 0.0
		npvDerivative := 0.0

		for idx, cashFlow := range flows {
			period := float64(idx)
			factor := math.Pow(1+irr, period)
			npv += cashFlow / factor
			if period > 0 {
				npvDerivative -= cashFlow * period / (factor * (1 + irr))
			}
		}

		if math.Abs(npv) < tolerance {
			return irr, nil
		}

		if npvDerivative == 0 {
			return 0, errors.New("cannot calculate IRR, derivative is zero")
		}
		irr = irr - npv/npvDerivative

		if irr < -1 || irr > 10 {
			return 0, errors.New("IRR calculation diverged")
		}
	}

	return irr, errors.New("IRR did not converge")
}

func normalizeRate(rate float64) float64 {
	if rate > 1 {
		return rate / 100
	}
	return rate
}

// CalculateTCEA calcula la Tasa de Costo Efectivo Anual ajustada a la frecuencia configurada.
func (fmc *FrenchMethodCalculator) CalculateTCEA(irr float64, periodsPerYear float64) float64 {
	if periodsPerYear <= 0 {
		periodsPerYear = 12
	}
	return math.Pow(1+irr, periodsPerYear) - 1
}
