package models

import (
	"time"

	"github.com/google/uuid"
)

// MortgageModel es el modelo de persistencia para GORM
type MortgageModel struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement"`
	UserID               uuid.UUID `gorm:"type:uuid;not null;index"`
	PropertyPrice        float64   `gorm:"not null"`
	DownPayment          float64   `gorm:"not null"`
	LoanAmount           float64   `gorm:"not null"`
	BonoTechoPropio      float64   `gorm:"default:0"`
	InterestRate         float64   `gorm:"not null"`
	RateType             string    `gorm:"type:varchar(20);not null"`
	TermMonths           int       `gorm:"not null"`
	TermYears            int       `gorm:"default:0"`
	GracePeriodMonths    int       `gorm:"default:0"`
	GracePeriodType      string    `gorm:"type:varchar(20);default:'NONE'"`
	Currency             string    `gorm:"type:varchar(3);not null"`
	PaymentFrequencyDays int       `gorm:"not null;default:30"`
	DaysInYear           int       `gorm:"not null;default:360"`
	AdministrationFee    float64   `gorm:"default:0"`
	Portes               float64   `gorm:"default:0"`
	AdditionalCosts      float64   `gorm:"default:0"`
	LifeInsuranceRate    float64   `gorm:"default:0"`
	PropertyInsurance    float64   `gorm:"default:0"`
	EvaluationFee        float64   `gorm:"default:0"`
	DisbursementFee      float64   `gorm:"default:0"`

	// Resultados calculados
	PrincipalFinanced float64 `gorm:"not null"`
	PeriodicRate      float64 `gorm:"not null"`
	FixedInstallment  float64 `gorm:"not null"`
	TotalInterestPaid float64 `gorm:"not null"`
	TotalPaid         float64 `gorm:"not null"`
	TotalPaidWithFees float64 `gorm:"not null"`
	TotalCharges      float64 `gorm:"not null"`
	TotalInsurance    float64 `gorm:"not null"`
	TotalAdmin        float64 `gorm:"not null"`
	NPV               float64 `gorm:"default:0"`
	IRR               float64 `gorm:"not null"`
	FlowIRR           float64 `gorm:"not null"`
	TCEA              float64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relación con los items del cronograma
	PaymentScheduleItems []PaymentScheduleItemModel `gorm:"foreignKey:MortgageID;constraint:OnDelete:CASCADE"`
}

func (MortgageModel) TableName() string {
	return "mortgages"
}
