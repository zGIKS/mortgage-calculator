package resources

import (
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"math"
	"time"
)

// CalculateMortgageRequest representa la solicitud para calcular un crédito hipotecario
type CalculateMortgageRequest struct {
	PrecioVenta     float64 `json:"precio_venta" binding:"required,gt=0"`
	CuotaInicial    float64 `json:"cuota_inicial" binding:"gte=0"`
	MontoPrestamo   float64 `json:"monto_prestamo" binding:"required,gt=0"`
	BonoTechoPropio float64 `json:"bono_techo_propio" binding:"gte=0"`
	TasaAnual       float64 `json:"tasa_anual" binding:"required,gte=0"`
	TipoTasa        string  `json:"tipo_tasa" binding:"required,oneof=NOMINAL EFFECTIVE"`
	Frecuencia      string  `json:"frecuencia,omitempty" binding:"omitempty,oneof=MENSUAL BIMESTRAL TRIMESTRAL"`
	FrecuenciaPago  int     `json:"frecuencia_pago" binding:"omitempty,gt=0"`
	DiasAnio        int     `json:"dias_anio" binding:"required,gt=0"`
	PlazoMeses      int     `json:"plazo_meses" binding:"omitempty,gt=0"`
	NumeroAnios     int     `json:"numero_anios" binding:"omitempty,gte=0"`
	MesesGracia     int     `json:"meses_gracia" binding:"gte=0"`
	TipoGracia      string  `json:"tipo_gracia" binding:"required,oneof=NONE TOTAL PARTIAL"`
	Moneda          string  `json:"moneda" binding:"required,oneof=PEN USD"`
	TasaDescuento   float64 `json:"tasa_descuento" binding:"gte=0"`
	COK             float64 `json:"cok" binding:"omitempty,gte=0"`
	Portes          float64 `json:"portes" binding:"omitempty,gte=0"`
	GastosAdm       float64 `json:"gastos_administrativos" binding:"omitempty,gte=0"`
	SeguroDesg      float64 `json:"seguro_desgravamen" binding:"omitempty,gte=0"`
	SeguroInmueble  float64 `json:"seguro_inmueble_anual" binding:"omitempty,gte=0"`
	ComisionEval    float64 `json:"comision_evaluacion" binding:"omitempty,gte=0"`
	ComisionDesem   float64 `json:"comision_desembolso" binding:"omitempty,gte=0"`
	CostosMensuales float64 `json:"costos_mensuales_adicionales" binding:"omitempty,gte=0"`
}

// UpdateMortgageRequest representa la solicitud para actualizar un crédito hipotecario
type UpdateMortgageRequest struct {
	PrecioVenta     *float64 `json:"precio_venta,omitempty" binding:"omitempty,gt=0"`
	CuotaInicial    *float64 `json:"cuota_inicial,omitempty" binding:"omitempty,gte=0"`
	MontoPrestamo   *float64 `json:"monto_prestamo,omitempty" binding:"omitempty,gt=0"`
	BonoTechoPropio *float64 `json:"bono_techo_propio,omitempty" binding:"omitempty,gte=0"`
	TasaAnual       *float64 `json:"tasa_anual,omitempty" binding:"omitempty,gte=0"`
	TipoTasa        *string  `json:"tipo_tasa,omitempty" binding:"omitempty,oneof=NOMINAL EFFECTIVE"`
	Frecuencia      *string  `json:"frecuencia,omitempty" binding:"omitempty,oneof=MENSUAL BIMESTRAL TRIMESTRAL"`
	FrecuenciaPago  *int     `json:"frecuencia_pago,omitempty" binding:"omitempty,gt=0"`
	DiasAnio        *int     `json:"dias_anio,omitempty" binding:"omitempty,gt=0"`
	PlazoMeses      *int     `json:"plazo_meses,omitempty" binding:"omitempty,gt=0"`
	NumeroAnios     *int     `json:"numero_anios,omitempty" binding:"omitempty,gte=0"`
	MesesGracia     *int     `json:"meses_gracia,omitempty" binding:"omitempty,gte=0"`
	TipoGracia      *string  `json:"tipo_gracia,omitempty" binding:"omitempty,oneof=NONE TOTAL PARTIAL"`
	Moneda          *string  `json:"moneda,omitempty" binding:"omitempty,oneof=PEN USD"`
	TasaDescuento   *float64 `json:"tasa_descuento,omitempty" binding:"omitempty,gte=0"`
	COK             *float64 `json:"cok,omitempty" binding:"omitempty,gte=0"`
	Portes          *float64 `json:"portes,omitempty" binding:"omitempty,gte=0"`
	GastosAdm       *float64 `json:"gastos_administrativos,omitempty" binding:"omitempty,gte=0"`
	SeguroDesg      *float64 `json:"seguro_desgravamen,omitempty" binding:"omitempty,gte=0"`
	SeguroInmueble  *float64 `json:"seguro_inmueble_anual,omitempty" binding:"omitempty,gte=0"`
	ComisionEval    *float64 `json:"comision_evaluacion,omitempty" binding:"omitempty,gte=0"`
	ComisionDesem   *float64 `json:"comision_desembolso,omitempty" binding:"omitempty,gte=0"`
	CostosMensuales *float64 `json:"costos_mensuales_adicionales,omitempty" binding:"omitempty,gte=0"`
}

// PaymentScheduleItemResource representa un item del cronograma
type PaymentScheduleItemResource struct {
	Periodo               int     `json:"periodo"`
	NumeroAnio            int     `json:"numero_anio"`
	TasaPeriodo           float64 `json:"tasa_periodo"`
	Cuota                 float64 `json:"cuota"`
	CuotaTotal            float64 `json:"cuota_total"`
	Interes               float64 `json:"interes"`
	Amortizacion          float64 `json:"amortizacion"`
	Portes                float64 `json:"portes"`
	GastosAdministrativos float64 `json:"gastos_administrativos"`
	SeguroDesgravamen     float64 `json:"seguro_desgravamen"`
	SeguroInmueble        float64 `json:"seguro_inmueble"`
	CostosAdicionales     float64 `json:"costos_adicionales"`
	SaldoFinal            float64 `json:"saldo_final"`
	EsPeriodoGracia       bool    `json:"es_periodo_gracia"`
	TipoGracia            string  `json:"tipo_gracia,omitempty"`
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
	NumeroAnios     int     `json:"numero_anios"`
	MesesGracia     int     `json:"meses_gracia"`
	TipoGracia      string  `json:"tipo_gracia"`
	Moneda          string  `json:"moneda"`
	FrecuenciaPago  int     `json:"frecuencia_pago"`
	DiasAnio        int     `json:"dias_anio"`
	Portes          float64 `json:"portes"`
	GastosAdm       float64 `json:"gastos_administrativos"`
	SeguroDesg      float64 `json:"seguro_desgravamen"`
	SeguroInmueble  float64 `json:"seguro_inmueble_anual"`
	ComisionEval    float64 `json:"comision_evaluacion"`
	ComisionDesem   float64 `json:"comision_desembolso"`
	CostosMensuales float64 `json:"costos_mensuales_adicionales"`
	CuotasPorAnio   int     `json:"cuotas_por_anio"`
	NumeroCuotas    int     `json:"numero_cuotas"`

	// Resultados calculados
	SaldoFinanciar    float64                       `json:"saldo_financiar"`
	TasaPeriodo       float64                       `json:"tasa_periodo"`
	CuotaFija         float64                       `json:"cuota_fija"`
	CuotaTotal        float64                       `json:"cuota_total"`
	CronogramaPagos   []PaymentScheduleItemResource `json:"cronograma_pagos"`
	TotalIntereses    float64                       `json:"total_intereses"`
	TotalPagado       float64                       `json:"total_pagado"`
	TotalPagadoCargos float64                       `json:"total_pagado_con_cargos"`
	TotalCargos       float64                       `json:"total_cargos"`
	TotalSeguros      float64                       `json:"total_seguros"`
	TotalGastos       float64                       `json:"total_gastos"`
	VAN               float64                       `json:"van"`
	TIR               float64                       `json:"tir"`
	TIRFlujo          float64                       `json:"tir_flujo"`
	TEA               float64                       `json:"tea"`
	TCEA              float64                       `json:"tcea"`

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
				Periodo:               item.Period,
				NumeroAnio:            item.YearNumber,
				TasaPeriodo:           item.PeriodicRateApplied,
				Cuota:                 item.Installment,
				CuotaTotal:            item.TotalInstallment,
				Interes:               item.Interest,
				Amortizacion:          item.Amortization,
				Portes:                item.Portes,
				GastosAdministrativos: item.AdministrationFee,
				SeguroDesgravamen:     item.LifeInsurance,
				SeguroInmueble:        item.PropertyInsurance,
				CostosAdicionales:     item.AdditionalCosts,
				SaldoFinal:            item.RemainingBalance,
				EsPeriodoGracia:       item.IsGracePeriod,
				TipoGracia:            item.GraceType,
			})
		}
	}

	cuotasPorAnio := int(math.Round(mortgage.PeriodsPerYear()))
	numeroCuotas := mortgage.TermMonths()
	if numeroCuotas == 0 && len(scheduleItems) > 0 {
		numeroCuotas = len(scheduleItems)
	}

	cuotaTotal := mortgage.FixedInstallment()
	if len(scheduleItems) > 0 {
		cuotaTotal = scheduleItems[0].CuotaTotal
	}

	tea := math.Pow(1+mortgage.PeriodicRate(), mortgage.PeriodsPerYear()) - 1

	return MortgageResponse{
		ID:                mortgage.ID().Value(),
		UserID:            mortgage.UserID().String(),
		PrecioVenta:       mortgage.PropertyPrice(),
		CuotaInicial:      mortgage.DownPayment(),
		MontoPrestamo:     mortgage.LoanAmount(),
		BonoTechoPropio:   mortgage.BonoTechoPropio(),
		TasaAnual:         mortgage.InterestRate(),
		TipoTasa:          mortgage.RateType().String(),
		PlazoMeses:        mortgage.TermMonths(),
		NumeroAnios:       mortgage.TermYears(),
		MesesGracia:       mortgage.GracePeriodMonths(),
		TipoGracia:        mortgage.GracePeriodType().String(),
		Moneda:            mortgage.Currency().String(),
		FrecuenciaPago:    mortgage.PaymentFrequencyDays(),
		DiasAnio:          mortgage.DaysInYear(),
		Portes:            mortgage.Portes(),
		GastosAdm:         mortgage.AdministrationFee(),
		SeguroDesg:        mortgage.LifeInsuranceRate(),
		SeguroInmueble:    mortgage.PropertyInsuranceRate(),
		ComisionEval:      mortgage.EvaluationFee(),
		ComisionDesem:     mortgage.DisbursementFee(),
		CostosMensuales:   mortgage.AdditionalCosts(),
		CuotasPorAnio:     cuotasPorAnio,
		NumeroCuotas:      numeroCuotas,
		SaldoFinanciar:    mortgage.PrincipalFinanced(),
		TasaPeriodo:       mortgage.PeriodicRate(),
		CuotaFija:         mortgage.FixedInstallment(),
		CuotaTotal:        cuotaTotal,
		CronogramaPagos:   scheduleItems,
		TotalIntereses:    mortgage.TotalInterestPaid(),
		TotalPagado:       mortgage.TotalPaid(),
		TotalPagadoCargos: mortgage.TotalPaidWithFees(),
		TotalCargos:       mortgage.TotalCharges(),
		TotalSeguros:      mortgage.TotalInsurance(),
		TotalGastos:       mortgage.TotalAdmin(),
		VAN:               mortgage.NPV(),
		TIR:               mortgage.IRR(),
		TIRFlujo:          mortgage.FlowIRR(),
		TEA:               tea,
		TCEA:              mortgage.TCEA(),
		CreatedAt:         mortgage.CreatedAt(),
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
