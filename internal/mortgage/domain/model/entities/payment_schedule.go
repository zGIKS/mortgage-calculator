package entities

// PaymentScheduleItem representa una fila del cronograma de pagos
type PaymentScheduleItem struct {
	Period              int     `json:"period"`                // Número de periodo (mes, bimestre, trimestre, etc.)
	YearNumber          int     `json:"year_number"`           // Año al que pertenece el periodo (1-n)
	PeriodicRateApplied float64 `json:"periodic_rate_applied"` // Tasa efectiva del periodo (por ejemplo, TET)
	Installment         float64 `json:"installment"`           // Cuota base (sin seguros ni gastos)
	TotalInstallment    float64 `json:"total_installment"`     // Cuota total (incluye seguros y gastos)
	Interest            float64 `json:"interest"`              // Interés del periodo (I_k)
	Amortization        float64 `json:"amortization"`          // Amortización del capital (C_k)
	AdministrationFee   float64 `json:"administration_fee"`    // Gastos administrativos del periodo
	Portes              float64 `json:"portes"`                // Portes u otros costos fijos del periodo
	LifeInsurance       float64 `json:"life_insurance"`        // Seguro de desgravamen del periodo
	PropertyInsurance   float64 `json:"property_insurance"`    // Seguro de inmueble del periodo
	AdditionalCosts     float64 `json:"additional_costs"`      // Otros costos mensuales adicionales
	RemainingBalance    float64 `json:"remaining_balance"`     // Saldo restante después del pago
	IsGracePeriod       bool    `json:"is_grace_period"`       // Indica si es periodo de gracia
	GraceType           string  `json:"grace_type,omitempty"`  // Tipo de gracia aplicada en el periodo
}

// PaymentSchedule representa el cronograma completo de pagos
type PaymentSchedule struct {
	Items []PaymentScheduleItem `json:"items"`
}

func NewPaymentSchedule() *PaymentSchedule {
	return &PaymentSchedule{
		Items: make([]PaymentScheduleItem, 0),
	}
}

func (ps *PaymentSchedule) AddItem(item PaymentScheduleItem) {
	ps.Items = append(ps.Items, item)
}

func (ps *PaymentSchedule) GetItems() []PaymentScheduleItem {
	return ps.Items
}

// TotalInterestPaid calcula el total de intereses pagados
func (ps *PaymentSchedule) TotalInterestPaid() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.Interest
	}
	return total
}

// TotalPaid calcula el total pagado (cuotas base)
func (ps *PaymentSchedule) TotalPaid() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.Installment
	}
	return total
}

// TotalPaidWithCharges calcula el total pagado considerando seguros y gastos
func (ps *PaymentSchedule) TotalPaidWithCharges() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.TotalInstallment
	}
	return total
}

// TotalCharges calcula la suma de cargos adicionales (seguros, portes, gastos, etc.)
func (ps *PaymentSchedule) TotalCharges() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.TotalInstallment - item.Installment
	}
	return total
}

// TotalInsurance calcula el total de seguros pagados
func (ps *PaymentSchedule) TotalInsurance() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.LifeInsurance + item.PropertyInsurance
	}
	return total
}

// TotalAdminFees calcula la suma de gastos administrativos y portes
func (ps *PaymentSchedule) TotalAdminFees() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.AdministrationFee + item.Portes + item.AdditionalCosts
	}
	return total
}
