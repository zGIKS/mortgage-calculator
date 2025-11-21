package models

import (
	"time"

	"github.com/google/uuid"
)

// MortgageModel es el modelo de persistencia para GORM
type MortgageModel struct {
	ID                   uint64     `gorm:"primaryKey;autoIncrement"`
	UserID               uuid.UUID  `gorm:"type:uuid;not null;index"`
	PropertyPrice        float64    `gorm:"not null"`
	DownPayment          float64    `gorm:"not null"`
	LoanAmount           float64    `gorm:"not null"`
	BonoTechoPropio      float64    `gorm:"default:0"`
	InterestRate         float64    `gorm:"not null"`
	RateType             string     `gorm:"type:varchar(20);not null"`
	BankID               *uuid.UUID `gorm:"type:uuid"`
	BankName             string     `gorm:"type:varchar(120)"`
	TermMonths           int        `gorm:"not null"`
	GracePeriodMonths    int        `gorm:"default:0"`
	GracePeriodType      string     `gorm:"type:varchar(20);default:'NONE'"`
	Currency             string     `gorm:"type:varchar(3);not null"`
	PaymentFrequencyDays int        `gorm:"not null;default:30"`
	DaysInYear           int        `gorm:"not null;default:360"`

	// Resultados calculados
	PrincipalFinanced float64 `gorm:"not null"`
	PeriodicRate      float64 `gorm:"not null"`
	FixedInstallment  float64 `gorm:"not null"`
	TotalInterestPaid float64 `gorm:"not null"`
	TotalPaid         float64 `gorm:"not null"`
	NPV               float64 `gorm:"default:0"`
	IRR               float64 `gorm:"not null"`
	TCEA              float64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relación con los items del cronograma
	PaymentScheduleItems []PaymentScheduleItemModel `gorm:"foreignKey:MortgageID;constraint:OnDelete:CASCADE"`
}

func (MortgageModel) TableName() string {
	return "mortgages"
}
