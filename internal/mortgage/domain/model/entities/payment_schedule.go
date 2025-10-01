package entities

// PaymentScheduleItem representa una fila del cronograma de pagos
type PaymentScheduleItem struct {
	Period           int     `json:"period"`            // Número de periodo (mes)
	Installment      float64 `json:"installment"`       // Cuota a pagar (A)
	Interest         float64 `json:"interest"`          // Interés del periodo (I_k)
	Amortization     float64 `json:"amortization"`      // Amortización del capital (C_k)
	RemainingBalance float64 `json:"remaining_balance"` // Saldo restante después del pago
	IsGracePeriod    bool    `json:"is_grace_period"`   // Indica si es periodo de gracia
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

// TotalPaid calcula el total pagado (cuotas)
func (ps *PaymentSchedule) TotalPaid() float64 {
	total := 0.0
	for _, item := range ps.Items {
		total += item.Installment
	}
	return total
}
