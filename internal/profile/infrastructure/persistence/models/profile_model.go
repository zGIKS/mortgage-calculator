package models

import (
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"time"

	"github.com/google/uuid"
)

type ProfileModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;column:user_id"`
	DNIEncrypted   string    `gorm:"type:varchar(255);not null;uniqueIndex;column:dni_encrypted"`
	FirstName      string    `gorm:"type:varchar(100);not null;column:first_name"`
	FirstLastName  string    `gorm:"type:varchar(100);not null;column:first_last_name"`
	SecondLastName string    `gorm:"type:varchar(100);not null;column:second_last_name"`
	PhoneNumber    string    `gorm:"type:varchar(9);not null;column:phone_number"`
	MonthlyIncome  float64   `gorm:"type:decimal(12,2);not null;column:monthly_income"`
	Currency       string    `gorm:"type:varchar(3);not null;column:currency"`
	MaritalStatus  string    `gorm:"type:varchar(20);not null;column:marital_status"`
	IsFirstHome    bool      `gorm:"not null;column:is_first_home"`
	HasOwnLand     bool      `gorm:"not null;column:has_own_land"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (ProfileModel) TableName() string {
	return "profiles"
}

func (m *ProfileModel) ToEntity() (*entities.Profile, error) {
	profileID, err := valueobjects.NewProfileID(m.ID)
	if err != nil {
		return nil, err
	}

	userID, err := valueobjects.NewUserID(m.UserID)
	if err != nil {
		return nil, err
	}

	// Create DNI from encrypted value (will be decrypted by repository)
	dni := valueobjects.NewDNIFromEncrypted(m.DNIEncrypted)

	// Handle optional phone number
	var phoneNumber valueobjects.PhoneNumber
	if m.PhoneNumber == "" {
		phoneNumber = valueobjects.EmptyPhoneNumber()
	} else {
		phoneNumber, err = valueobjects.NewPhoneNumber(m.PhoneNumber)
		if err != nil {
			return nil, err
		}
	}

	// Handle optional monthly income
	var monthlyIncome valueobjects.MonthlyIncome
	if m.MonthlyIncome == 0 {
		monthlyIncome = valueobjects.EmptyMonthlyIncome()
	} else {
		monthlyIncome, err = valueobjects.NewMonthlyIncome(m.MonthlyIncome, valueobjects.Currency(m.Currency))
		if err != nil {
			return nil, err
		}
	}

	// Handle optional marital status
	var maritalStatus valueobjects.MaritalStatus
	if m.MaritalStatus == "" {
		maritalStatus = valueobjects.EmptyMaritalStatus()
	} else {
		maritalStatus, err = valueobjects.NewMaritalStatus(m.MaritalStatus)
		if err != nil {
			return nil, err
		}
	}

	return entities.ReconstructProfile(
		profileID,
		userID,
		dni,
		m.FirstName,
		m.FirstLastName,
		m.SecondLastName,
		phoneNumber,
		monthlyIncome,
		maritalStatus,
		m.IsFirstHome,
		m.HasOwnLand,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func FromEntity(profile *entities.Profile) *ProfileModel {
	return &ProfileModel{
		ID:             profile.ID().Value(),
		UserID:         profile.UserID().Value(),
		DNIEncrypted:   profile.DNI().EncryptedValue(),
		FirstName:      profile.FirstName(),
		FirstLastName:  profile.FirstLastName(),
		SecondLastName: profile.SecondLastName(),
		PhoneNumber:    profile.PhoneNumber().Value(),
		MonthlyIncome:  profile.MonthlyIncome().Amount(),
		Currency:       string(profile.MonthlyIncome().Currency()),
		MaritalStatus:  profile.MaritalStatus().String(),
		IsFirstHome:    profile.IsFirstHome(),
		HasOwnLand:     profile.HasOwnLand(),
		CreatedAt:      profile.CreatedAt(),
		UpdatedAt:      profile.UpdatedAt(),
	}
}
