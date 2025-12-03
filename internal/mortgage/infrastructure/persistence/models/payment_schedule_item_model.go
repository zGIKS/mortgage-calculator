package models

import "github.com/google/uuid"

// PaymentScheduleItemModel representa un item del cronograma de pagos en la BD
type PaymentScheduleItemModel struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MortgageID        uint64    `gorm:"not null;index:idx_mortgage_period"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index"`
	Period            int       `gorm:"not null;index:idx_mortgage_period"`
	YearNumber        int       `gorm:"not null;default:1"`
	PeriodicRate      float64   `gorm:"not null;default:0"`
	Installment       float64   `gorm:"not null"`
	TotalInstallment  float64   `gorm:"not null"`
	Interest          float64   `gorm:"not null"`
	Amortization      float64   `gorm:"not null"`
	Administration    float64   `gorm:"not null;default:0"`
	Portes            float64   `gorm:"not null;default:0"`
	LifeInsurance     float64   `gorm:"not null;default:0"`
	PropertyInsurance float64   `gorm:"not null;default:0"`
	AdditionalCosts   float64   `gorm:"not null;default:0"`
	RemainingBalance  float64   `gorm:"not null"`
	IsGracePeriod     bool      `gorm:"default:false"`
	GraceType         string    `gorm:"type:varchar(20);default:''"`
}

func (PaymentScheduleItemModel) TableName() string {
	return "payment_schedule_items"
}
