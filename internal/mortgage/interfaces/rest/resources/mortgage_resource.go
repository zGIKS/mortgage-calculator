package resources

import (
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"time"
)

// CalculateMortgageRequest representa la solicitud para calcular un crédito hipotecario
type CalculateMortgageRequest struct {
	PrecioVenta     float64 `json:"precio_venta" binding:"required,gt=0"`
	CuotaInicial    float64 `json:"cuota_inicial" binding:"required,gte=0"`
	MontoPrestamo   float64 `json:"monto_prestamo" binding:"required,gt=0"`
	BonoTechoPropio float64 `json:"bono_techo_propio" binding:"gte=0"`
	TasaAnual       float64 `json:"tasa_anual" binding:"required,gte=0"`
	TipoTasa        string  `json:"tipo_tasa" binding:"required,oneof=NOMINAL EFFECTIVE"`
	FrecuenciaPago  int     `json:"frecuencia_pago" binding:"required,gt=0"`
	DiasAnio        int     `json:"dias_anio" binding:"required,gt=0"`
	PlazoMeses      int     `json:"plazo_meses" binding:"required,gt=0"`
	MesesGracia     int     `json:"meses_gracia" binding:"gte=0"`
	TipoGracia      string  `json:"tipo_gracia" binding:"required,oneof=NONE TOTAL PARTIAL"`
	Moneda          string  `json:"moneda" binding:"required,oneof=PEN USD"`
	TasaDescuento   float64 `json:"tasa_descuento" binding:"gte=0"`
}

// UpdateMortgageRequest representa la solicitud para actualizar un crédito hipotecario
type UpdateMortgageRequest struct {
	PrecioVenta     *float64 `json:"precio_venta,omitempty" binding:"omitempty,gt=0"`
	CuotaInicial    *float64 `json:"cuota_inicial,omitempty" binding:"omitempty,gte=0"`
	MontoPrestamo   *float64 `json:"monto_prestamo,omitempty" binding:"omitempty,gt=0"`
	BonoTechoPropio *float64 `json:"bono_techo_propio,omitempty" binding:"omitempty,gte=0"`
	TasaAnual       *float64 `json:"tasa_anual,omitempty" binding:"omitempty,gte=0"`
	TipoTasa        *string  `json:"tipo_tasa,omitempty" binding:"omitempty,oneof=NOMINAL EFFECTIVE"`
	FrecuenciaPago  *int     `json:"frecuencia_pago,omitempty" binding:"omitempty,gt=0"`
	DiasAnio        *int     `json:"dias_anio,omitempty" binding:"omitempty,gt=0"`
	PlazoMeses      *int     `json:"plazo_meses,omitempty" binding:"omitempty,gt=0"`
	MesesGracia     *int     `json:"meses_gracia,omitempty" binding:"omitempty,gte=0"`
	TipoGracia      *string  `json:"tipo_gracia,omitempty" binding:"omitempty,oneof=NONE TOTAL PARTIAL"`
	Moneda          *string  `json:"moneda,omitempty" binding:"omitempty,oneof=PEN USD"`
	TasaDescuento   *float64 `json:"tasa_descuento,omitempty" binding:"omitempty,gte=0"`
}

// PaymentScheduleItemResource representa un item del cronograma
type PaymentScheduleItemResource struct {
	Periodo         int     `json:"periodo"`
	Cuota           float64 `json:"cuota"`
	Interes         float64 `json:"interes"`
	Amortizacion    float64 `json:"amortizacion"`
	SaldoFinal      float64 `json:"saldo_final"`
	EsPeriodoGracia bool    `json:"es_periodo_gracia"`
}

// MortgageResponse representa la respuesta completa con todos los cálculos
type MortgageResponse struct {
	ID              uint64  `json:"id"`
	UserID          string  `json:"user_id"`
	PrecioVenta     float64 `json:"precio_venta"`
	CuotaInicial    float64 `json:"cuota_inicial"`
	MontoPrestamo   float64 `json:"monto_prestamo"`
	BonoTechoPropio float64 `json:"bono_techo_propio"`
	TasaAnual       float64 `json:"tasa_anual"`
	TipoTasa        string  `json:"tipo_tasa"`
	PlazoMeses      int     `json:"plazo_meses"`
	MesesGracia     int     `json:"meses_gracia"`
	TipoGracia      string  `json:"tipo_gracia"`
	Moneda          string  `json:"moneda"`
	FrecuenciaPago  int     `json:"frecuencia_pago"`
	DiasAnio        int     `json:"dias_anio"`

	// Resultados calculados
	SaldoFinanciar  float64                       `json:"saldo_financiar"`
	TasaPeriodo     float64                       `json:"tasa_periodo"`
	CuotaFija       float64                       `json:"cuota_fija"`
	CronogramaPagos []PaymentScheduleItemResource `json:"cronograma_pagos"`
	TotalIntereses  float64                       `json:"total_intereses"`
	TotalPagado     float64                       `json:"total_pagado"`
	VAN             float64                       `json:"van"`
	TIR             float64                       `json:"tir"`
	TCEA            float64                       `json:"tcea"`

	CreatedAt time.Time `json:"created_at"`
}

// MortgageSummaryResource representa un resumen de hipoteca (para listas)
type MortgageSummaryResource struct {
	ID            uint64    `json:"id"`
	UserID        string    `json:"user_id"`
	PrecioVenta   float64   `json:"precio_venta"`
	MontoPrestamo float64   `json:"monto_prestamo"`
	Moneda        string    `json:"moneda"`
	PlazoMeses    int       `json:"plazo_meses"`
	CuotaFija     float64   `json:"cuota_fija"`
	TCEA          float64   `json:"tcea"`
	CreatedAt     time.Time `json:"created_at"`
}

// TransformToMortgageResponse transforma una entidad Mortgage a MortgageResponse
func TransformToMortgageResponse(mortgage *entities.Mortgage) MortgageResponse {
	scheduleItems := make([]PaymentScheduleItemResource, 0)
	if mortgage.PaymentSchedule() != nil {
		for _, item := range mortgage.PaymentSchedule().GetItems() {
			scheduleItems = append(scheduleItems, PaymentScheduleItemResource{
				Periodo:         item.Period,
				Cuota:           item.Installment,
				Interes:         item.Interest,
				Amortizacion:    item.Amortization,
				SaldoFinal:      item.RemainingBalance,
				EsPeriodoGracia: item.IsGracePeriod,
			})
		}
	}

	return MortgageResponse{
		ID:              mortgage.ID().Value(),
		UserID:          mortgage.UserID().String(),
		PrecioVenta:     mortgage.PropertyPrice(),
		CuotaInicial:    mortgage.DownPayment(),
		MontoPrestamo:   mortgage.LoanAmount(),
		BonoTechoPropio: mortgage.BonoTechoPropio(),
		TasaAnual:       mortgage.InterestRate(),
		TipoTasa:        mortgage.RateType().String(),
		PlazoMeses:      mortgage.TermMonths(),
		MesesGracia:     mortgage.GracePeriodMonths(),
		TipoGracia:      mortgage.GracePeriodType().String(),
		Moneda:          mortgage.Currency().String(),
		FrecuenciaPago:  mortgage.PaymentFrequencyDays(),
		DiasAnio:        mortgage.DaysInYear(),
		SaldoFinanciar:  mortgage.PrincipalFinanced(),
		TasaPeriodo:     mortgage.PeriodicRate(),
		CuotaFija:       mortgage.FixedInstallment(),
		CronogramaPagos: scheduleItems,
		TotalIntereses:  mortgage.TotalInterestPaid(),
		TotalPagado:     mortgage.TotalPaid(),
		VAN:             mortgage.NPV(),
		TIR:             mortgage.IRR(),
		TCEA:            mortgage.TCEA(),
		CreatedAt:       mortgage.CreatedAt(),
	}
}

// TransformToMortgageSummary transforma una entidad Mortgage a MortgageSummaryResource
func TransformToMortgageSummary(mortgage *entities.Mortgage) MortgageSummaryResource {
	return MortgageSummaryResource{
		ID:            mortgage.ID().Value(),
		UserID:        mortgage.UserID().String(),
		PrecioVenta:   mortgage.PropertyPrice(),
		MontoPrestamo: mortgage.LoanAmount(),
		Moneda:        mortgage.Currency().String(),
		PlazoMeses:    mortgage.TermMonths(),
		CuotaFija:     mortgage.FixedInstallment(),
		TCEA:          mortgage.TCEA(),
		CreatedAt:     mortgage.CreatedAt(),
	}
}
