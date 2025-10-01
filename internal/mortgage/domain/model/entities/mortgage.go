package entities

import (
	"finanzas-backend/internal/mortgage/domain/model/valueobjects"
	"time"
)

// Mortgage representa un crédito hipotecario calculado con método francés
type Mortgage struct {
	id                  valueobjects.MortgageID
	userID              valueobjects.UserID
	propertyPrice       float64 // Precio de la vivienda
	downPayment         float64 // Cuota inicial
	loanAmount          float64 // Monto del préstamo solicitado
	bonoTechoPropio     float64 // Bono Techo Propio (subsidio)
	interestRate        float64 // Tasa de interés (TNA o TEA según rateType)
	rateType            valueobjects.RateType
	termMonths          int // Plazo en meses
	gracePeriodMonths   int // Número de meses de gracia
	gracePeriodType     valueobjects.GracePeriodType
	currency            valueobjects.Currency

	// Resultados calculados
	principalFinanced   float64          // Principal financiado = loanAmount - bonoTechoPropio
	periodicRate        float64          // Tasa efectiva por periodo (mensual)
	fixedInstallment    float64          // Cuota fija (después de gracia)
	paymentSchedule     *PaymentSchedule // Cronograma de pagos
	totalInterestPaid   float64          // Total de intereses pagados
	totalPaid           float64          // Total pagado
	npv                 float64          // Valor Actual Neto (VAN)
	irr                 float64          // Tasa Interna de Retorno (TIR)
	tcea                float64          // Tasa de Costo Efectivo Anual

	createdAt           time.Time
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
	gracePeriodMonths int,
	gracePeriodType valueobjects.GracePeriodType,
	currency valueobjects.Currency,
) (*Mortgage, error) {
	return &Mortgage{
		userID:            userID,
		propertyPrice:     propertyPrice,
		downPayment:       downPayment,
		loanAmount:        loanAmount,
		bonoTechoPropio:   bonoTechoPropio,
		interestRate:      interestRate,
		rateType:          rateType,
		termMonths:        termMonths,
		gracePeriodMonths: gracePeriodMonths,
		gracePeriodType:   gracePeriodType,
		currency:          currency,
		createdAt:         time.Now(),
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
	gracePeriodMonths int,
	gracePeriodType valueobjects.GracePeriodType,
	currency valueobjects.Currency,
	principalFinanced float64,
	periodicRate float64,
	fixedInstallment float64,
	totalInterestPaid float64,
	totalPaid float64,
	npv float64,
	irr float64,
	tcea float64,
	createdAt time.Time,
) *Mortgage {
	return &Mortgage{
		id:                id,
		userID:            userID,
		propertyPrice:     propertyPrice,
		downPayment:       downPayment,
		loanAmount:        loanAmount,
		bonoTechoPropio:   bonoTechoPropio,
		interestRate:      interestRate,
		rateType:          rateType,
		termMonths:        termMonths,
		gracePeriodMonths: gracePeriodMonths,
		gracePeriodType:   gracePeriodType,
		currency:          currency,
		principalFinanced: principalFinanced,
		periodicRate:      periodicRate,
		fixedInstallment:  fixedInstallment,
		totalInterestPaid: totalInterestPaid,
		totalPaid:         totalPaid,
		npv:               npv,
		irr:               irr,
		tcea:              tcea,
		createdAt:         createdAt,
	}
}

// Getters
func (m *Mortgage) ID() valueobjects.MortgageID                  { return m.id }
func (m *Mortgage) UserID() valueobjects.UserID                  { return m.userID }
func (m *Mortgage) PropertyPrice() float64                       { return m.propertyPrice }
func (m *Mortgage) DownPayment() float64                         { return m.downPayment }
func (m *Mortgage) LoanAmount() float64                          { return m.loanAmount }
func (m *Mortgage) BonoTechoPropio() float64                     { return m.bonoTechoPropio }
func (m *Mortgage) InterestRate() float64                        { return m.interestRate }
func (m *Mortgage) RateType() valueobjects.RateType              { return m.rateType }
func (m *Mortgage) TermMonths() int                              { return m.termMonths }
func (m *Mortgage) GracePeriodMonths() int                       { return m.gracePeriodMonths }
func (m *Mortgage) GracePeriodType() valueobjects.GracePeriodType { return m.gracePeriodType }
func (m *Mortgage) Currency() valueobjects.Currency              { return m.currency }
func (m *Mortgage) PrincipalFinanced() float64                   { return m.principalFinanced }
func (m *Mortgage) PeriodicRate() float64                        { return m.periodicRate }
func (m *Mortgage) FixedInstallment() float64                    { return m.fixedInstallment }
func (m *Mortgage) PaymentSchedule() *PaymentSchedule            { return m.paymentSchedule }
func (m *Mortgage) TotalInterestPaid() float64                   { return m.totalInterestPaid }
func (m *Mortgage) TotalPaid() float64                           { return m.totalPaid }
func (m *Mortgage) NPV() float64                                 { return m.npv }
func (m *Mortgage) IRR() float64                                 { return m.irr }
func (m *Mortgage) TCEA() float64                                { return m.tcea }
func (m *Mortgage) CreatedAt() time.Time                         { return m.createdAt }

// Setters para resultados calculados
func (m *Mortgage) SetID(id valueobjects.MortgageID)           { m.id = id }
func (m *Mortgage) SetPrincipalFinanced(value float64)         { m.principalFinanced = value }
func (m *Mortgage) SetPeriodicRate(value float64)              { m.periodicRate = value }
func (m *Mortgage) SetFixedInstallment(value float64)          { m.fixedInstallment = value }
func (m *Mortgage) SetPaymentSchedule(schedule *PaymentSchedule) { m.paymentSchedule = schedule }
func (m *Mortgage) SetTotalInterestPaid(value float64)         { m.totalInterestPaid = value }
func (m *Mortgage) SetTotalPaid(value float64)                 { m.totalPaid = value }
func (m *Mortgage) SetNPV(value float64)                       { m.npv = value }
func (m *Mortgage) SetIRR(value float64)                       { m.irr = value }
func (m *Mortgage) SetTCEA(value float64)                      { m.tcea = value }
