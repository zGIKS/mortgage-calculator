package models

import "time"

// BankModel persists the configuration per financial institution.
type BankModel struct {
	ID                   string    `gorm:"type:varchar(50);primaryKey"`
	Name                 string    `gorm:"type:varchar(120);not null"`
	RateType             string    `gorm:"type:varchar(20);not null"`
	PaymentFrequencyDays int       `gorm:"not null"`
	DaysInYear           int       `gorm:"not null"`
	IncludesInflation    bool      `gorm:"default:false"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
}

func (BankModel) TableName() string {
	return "banks"
}
