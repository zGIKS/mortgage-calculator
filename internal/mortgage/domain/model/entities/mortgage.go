package entities

import (
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"time"
)

// Mortgage representa un crédito hipotecario calculado con método francés
type Mortgage struct {
	id                   valueobjects.MortgageID
	userID               valueobjects.UserID
	propertyPrice        float64 // Precio de la vivienda
	downPayment          float64 // Cuota inicial
	loanAmount           float64 // Monto del préstamo solicitado
	bonoTechoPropio      float64 // Bono Techo Propio (subsidio)
	interestRate         float64 // Tasa de interés (TNA o TEA según rateType)
	rateType             valueobjects.RateType
	termMonths           int // Plazo en meses o número de periodos
	termYears            int // Plazo en años (se usa para derivar número de cuotas)
	gracePeriodMonths    int // Número de meses de gracia
	gracePeriodType      valueobjects.GracePeriodType
	currency             valueobjects.Currency
	paymentFrequencyDays int
	daysInYear           int
	adminFee             float64 // Gastos administrativos por periodo
	portes               float64 // Portes u otros costos fijos
	additionalCosts      float64 // Costos mensuales adicionales
	lifeInsuranceRate    float64 // Tasa mensual de seguro de desgravamen (decimal)
	propertyInsurance    float64 // Tasa anual de seguro de inmueble (decimal)
	evaluationFee        float64 // Comisión de evaluación (única)
	disbursementFee      float64 // Comisión de desembolso (única)

	// Resultados calculados
	principalFinanced float64          // Principal financiado = loanAmount - bonoTechoPropio
	periodicRate      float64          // Tasa efectiva por periodo (mensual)
	fixedInstallment  float64          // Cuota fija (después de gracia)
	paymentSchedule   *PaymentSchedule // Cronograma de pagos
	totalInterestPaid float64          // Total de intereses pagados
	totalPaid         float64          // Total pagado
	totalPaidWithFees float64          // Total pagado incluyendo seguros y gastos
	totalCharges      float64          // Total de cargos adicionales
	totalInsurance    float64          // Total de seguros
	totalAdmin        float64          // Total de gastos administrativos y portes
	npv               float64          // Valor Actual Neto (VAN)
	irr               float64          // Tasa Interna de Retorno (TIR) para cuota base
	flowIRR           float64          // TIR incluyendo cargos
	tcea              float64          // Tasa de Costo Efectivo Anual

	createdAt time.Time
}

func NewMortgage(
	userID valueobjects.UserID,
	propertyPrice float64,
	downPayment float64,
	loanAmount float64,
	bonoTechoPropio float64,
	interestRate float64,
	rateType valueobjects.RateType,
	termMonths int,
	termYears int,
	gracePeriodMonths int,
	gracePeriodType valueobjects.GracePeriodType,
	currency valueobjects.Currency,
	adminFee float64,
	portes float64,
	additionalCosts float64,
	lifeInsuranceRate float64,
	propertyInsurance float64,
	evaluationFee float64,
	disbursementFee float64,
) (*Mortgage, error) {
	return &Mortgage{
		userID:               userID,
		propertyPrice:        propertyPrice,
		downPayment:          downPayment,
		loanAmount:           loanAmount,
		bonoTechoPropio:      bonoTechoPropio,
		interestRate:         interestRate,
		rateType:             rateType,
		termMonths:           termMonths,
		termYears:            termYears,
		gracePeriodMonths:    gracePeriodMonths,
		gracePeriodType:      gracePeriodType,
		currency:             currency,
		paymentFrequencyDays: 30,
		daysInYear:           360,
		adminFee:             adminFee,
		portes:               portes,
		additionalCosts:      additionalCosts,
		lifeInsuranceRate:    lifeInsuranceRate,
		propertyInsurance:    propertyInsurance,
		evaluationFee:        evaluationFee,
		disbursementFee:      disbursementFee,
		createdAt:            time.Now(),
	}, nil
}

func ReconstructMortgage(
	id valueobjects.MortgageID,
	userID valueobjects.UserID,
	propertyPrice float64,
	downPayment float64,
	loanAmount float64,
	bonoTechoPropio float64,
	interestRate float64,
	rateType valueobjects.RateType,
	termMonths int,
	termYears int,
	gracePeriodMonths int,
	gracePeriodType valueobjects.GracePeriodType,
	currency valueobjects.Currency,
	paymentFrequencyDays int,
	daysInYear int,
	adminFee float64,
	portes float64,
	additionalCosts float64,
	lifeInsuranceRate float64,
	propertyInsurance float64,
	evaluationFee float64,
	disbursementFee float64,
	principalFinanced float64,
	periodicRate float64,
	fixedInstallment float64,
	totalInterestPaid float64,
	totalPaid float64,
	totalPaidWithFees float64,
	totalCharges float64,
	totalInsurance float64,
	totalAdmin float64,
	npv float64,
	irr float64,
	flowIRR float64,
	tcea float64,
	createdAt time.Time,
) *Mortgage {
	return &Mortgage{
		id:                   id,
		userID:               userID,
		propertyPrice:        propertyPrice,
		downPayment:          downPayment,
		loanAmount:           loanAmount,
		bonoTechoPropio:      bonoTechoPropio,
		interestRate:         interestRate,
		rateType:             rateType,
		termMonths:           termMonths,
		termYears:            termYears,
		gracePeriodMonths:    gracePeriodMonths,
		gracePeriodType:      gracePeriodType,
		currency:             currency,
		paymentFrequencyDays: paymentFrequencyDays,
		daysInYear:           daysInYear,
		adminFee:             adminFee,
		portes:               portes,
		additionalCosts:      additionalCosts,
		lifeInsuranceRate:    lifeInsuranceRate,
		propertyInsurance:    propertyInsurance,
		evaluationFee:        evaluationFee,
		disbursementFee:      disbursementFee,
		principalFinanced:    principalFinanced,
		periodicRate:         periodicRate,
		fixedInstallment:     fixedInstallment,
		totalInterestPaid:    totalInterestPaid,
		totalPaid:            totalPaid,
		totalPaidWithFees:    totalPaidWithFees,
		totalCharges:         totalCharges,
		totalInsurance:       totalInsurance,
		totalAdmin:           totalAdmin,
		npv:                  npv,
		irr:                  irr,
		flowIRR:              flowIRR,
		tcea:                 tcea,
		createdAt:            createdAt,
	}
}

// Getters
func (m *Mortgage) ID() valueobjects.MortgageID                   { return m.id }
func (m *Mortgage) UserID() valueobjects.UserID                   { return m.userID }
func (m *Mortgage) PropertyPrice() float64                        { return m.propertyPrice }
func (m *Mortgage) DownPayment() float64                          { return m.downPayment }
func (m *Mortgage) LoanAmount() float64                           { return m.loanAmount }
func (m *Mortgage) BonoTechoPropio() float64                      { return m.bonoTechoPropio }
func (m *Mortgage) InterestRate() float64                         { return m.interestRate }
func (m *Mortgage) RateType() valueobjects.RateType               { return m.rateType }
func (m *Mortgage) TermMonths() int                               { return m.termMonths }
func (m *Mortgage) TermYears() int                                { return m.termYears }
func (m *Mortgage) GracePeriodMonths() int                        { return m.gracePeriodMonths }
func (m *Mortgage) GracePeriodType() valueobjects.GracePeriodType { return m.gracePeriodType }
func (m *Mortgage) Currency() valueobjects.Currency               { return m.currency }
func (m *Mortgage) PaymentFrequencyDays() int                     { return m.paymentFrequencyDays }
func (m *Mortgage) DaysInYear() int                               { return m.daysInYear }
func (m *Mortgage) AdministrationFee() float64                    { return m.adminFee }
func (m *Mortgage) Portes() float64                               { return m.portes }
func (m *Mortgage) AdditionalCosts() float64                      { return m.additionalCosts }
func (m *Mortgage) LifeInsuranceRate() float64                    { return m.lifeInsuranceRate }
func (m *Mortgage) PropertyInsuranceRate() float64                { return m.propertyInsurance }
func (m *Mortgage) EvaluationFee() float64                        { return m.evaluationFee }
func (m *Mortgage) DisbursementFee() float64                      { return m.disbursementFee }
func (m *Mortgage) PeriodsPerYear() float64 {
	if m.paymentFrequencyDays > 0 && m.daysInYear > 0 {
		return float64(m.daysInYear) / float64(m.paymentFrequencyDays)
	}
	return 12.0
}
func (m *Mortgage) PrincipalFinanced() float64        { return m.principalFinanced }
func (m *Mortgage) PeriodicRate() float64             { return m.periodicRate }
func (m *Mortgage) FixedInstallment() float64         { return m.fixedInstallment }
func (m *Mortgage) PaymentSchedule() *PaymentSchedule { return m.paymentSchedule }
func (m *Mortgage) TotalInterestPaid() float64        { return m.totalInterestPaid }
func (m *Mortgage) TotalPaid() float64                { return m.totalPaid }
func (m *Mortgage) TotalPaidWithFees() float64        { return m.totalPaidWithFees }
func (m *Mortgage) TotalCharges() float64             { return m.totalCharges }
func (m *Mortgage) TotalInsurance() float64           { return m.totalInsurance }
func (m *Mortgage) TotalAdmin() float64               { return m.totalAdmin }
func (m *Mortgage) NPV() float64                      { return m.npv }
func (m *Mortgage) IRR() float64                      { return m.irr }
func (m *Mortgage) FlowIRR() float64                  { return m.flowIRR }
func (m *Mortgage) TCEA() float64                     { return m.tcea }
func (m *Mortgage) CreatedAt() time.Time              { return m.createdAt }

// Setters para resultados calculados
func (m *Mortgage) SetID(id valueobjects.MortgageID)             { m.id = id }
func (m *Mortgage) SetPrincipalFinanced(value float64)           { m.principalFinanced = value }
func (m *Mortgage) SetPeriodicRate(value float64)                { m.periodicRate = value }
func (m *Mortgage) SetFixedInstallment(value float64)            { m.fixedInstallment = value }
func (m *Mortgage) SetPaymentSchedule(schedule *PaymentSchedule) { m.paymentSchedule = schedule }
func (m *Mortgage) SetTotalInterestPaid(value float64)           { m.totalInterestPaid = value }
func (m *Mortgage) SetTotalPaid(value float64)                   { m.totalPaid = value }
func (m *Mortgage) SetTotalPaidWithFees(value float64)           { m.totalPaidWithFees = value }
func (m *Mortgage) SetTotalCharges(value float64)                { m.totalCharges = value }
func (m *Mortgage) SetTotalInsurance(value float64)              { m.totalInsurance = value }
func (m *Mortgage) SetTotalAdmin(value float64)                  { m.totalAdmin = value }
func (m *Mortgage) SetNPV(value float64)                         { m.npv = value }
func (m *Mortgage) SetIRR(value float64)                         { m.irr = value }
func (m *Mortgage) SetFlowIRR(value float64)                     { m.flowIRR = value }
func (m *Mortgage) SetTCEA(value float64)                        { m.tcea = value }
func (m *Mortgage) SetRateType(value valueobjects.RateType)      { m.rateType = value }
func (m *Mortgage) SetPaymentFrequencyDays(value int) {
	if value > 0 {
		m.paymentFrequencyDays = value
	}
}
func (m *Mortgage) SetDaysInYear(value int) {
	if value > 0 {
		m.daysInYear = value
	}
}
func (m *Mortgage) SetTermYears(value int) {
	if value >= 0 {
		m.termYears = value
	}
}
func (m *Mortgage) SetAdministrationFee(value float64) {
	if value >= 0 {
		m.adminFee = value
	}
}
func (m *Mortgage) SetPortes(value float64) {
	if value >= 0 {
		m.portes = value
	}
}
func (m *Mortgage) SetAdditionalCosts(value float64) {
	if value >= 0 {
		m.additionalCosts = value
	}
}
func (m *Mortgage) SetLifeInsuranceRate(value float64) {
	if value >= 0 {
		m.lifeInsuranceRate = value
	}
}
func (m *Mortgage) SetPropertyInsuranceRate(value float64) {
	if value >= 0 {
		m.propertyInsurance = value
	}
}
func (m *Mortgage) SetEvaluationFee(value float64) {
	if value >= 0 {
		m.evaluationFee = value
	}
}
func (m *Mortgage) SetDisbursementFee(value float64) {
	if value >= 0 {
		m.disbursementFee = value
	}
}
