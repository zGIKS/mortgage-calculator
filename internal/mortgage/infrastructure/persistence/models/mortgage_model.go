package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"finanzas-backend/internal/mortgage/domain/model/entities"
	"time"
)

// PaymentScheduleJSON es un tipo personalizado para serializar el cronograma como JSON en la BD
type PaymentScheduleJSON []entities.PaymentScheduleItem

func (p PaymentScheduleJSON) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PaymentScheduleJSON) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal PaymentScheduleJSON value")
	}
	return json.Unmarshal(bytes, p)
}

// MortgageModel es el modelo de persistencia para GORM
type MortgageModel struct {
	ID                uint64              `gorm:"primaryKey;autoIncrement"`
	UserID            uint64              `gorm:"not null;index"`
	PropertyPrice     float64             `gorm:"not null"`
	DownPayment       float64             `gorm:"not null"`
	LoanAmount        float64             `gorm:"not null"`
	BonoTechoPropio   float64             `gorm:"default:0"`
	InterestRate      float64             `gorm:"not null"`
	RateType          string              `gorm:"type:varchar(20);not null"`
	TermMonths        int                 `gorm:"not null"`
	GracePeriodMonths int                 `gorm:"default:0"`
	GracePeriodType   string              `gorm:"type:varchar(20);default:'NONE'"`
	Currency          string              `gorm:"type:varchar(3);not null"`

	// Resultados calculados
	PrincipalFinanced float64             `gorm:"not null"`
	PeriodicRate      float64             `gorm:"not null"`
	FixedInstallment  float64             `gorm:"not null"`
	PaymentSchedule   PaymentScheduleJSON `gorm:"type:jsonb"`
	TotalInterestPaid float64             `gorm:"not null"`
	TotalPaid         float64             `gorm:"not null"`
	NPV               float64             `gorm:"default:0"`
	IRR               float64             `gorm:"not null"`
	TCEA              float64             `gorm:"not null"`

	CreatedAt         time.Time           `gorm:"autoCreateTime"`
	UpdatedAt         time.Time           `gorm:"autoUpdateTime"`
}

func (MortgageModel) TableName() string {
	return "mortgages"
}
