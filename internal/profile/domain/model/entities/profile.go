package entities

import (
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"time"
)

type Profile struct {
	id             valueobjects.ProfileID
	userID         valueobjects.UserID
	dni            valueobjects.DNI
	firstName      string
	firstLastName  string
	secondLastName string
	phoneNumber    valueobjects.PhoneNumber
	monthlyIncome  valueobjects.MonthlyIncome
	maritalStatus  valueobjects.MaritalStatus
	isFirstHome    bool
	hasOwnLand     bool
	createdAt      time.Time
	updatedAt      time.Time
}

func NewProfile(
	userID valueobjects.UserID,
	dni valueobjects.DNI,
	firstName string,
	firstLastName string,
	secondLastName string,
	phoneNumber valueobjects.PhoneNumber,
	monthlyIncome valueobjects.MonthlyIncome,
	maritalStatus valueobjects.MaritalStatus,
	isFirstHome bool,
	hasOwnLand bool,
) (*Profile, error) {
	return &Profile{
		userID:         userID,
		dni:            dni,
		firstName:      firstName,
		firstLastName:  firstLastName,
		secondLastName: secondLastName,
		phoneNumber:    phoneNumber,
		monthlyIncome:  monthlyIncome,
		maritalStatus:  maritalStatus,
		isFirstHome:    isFirstHome,
		hasOwnLand:     hasOwnLand,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
	}, nil
}

func ReconstructProfile(
	id valueobjects.ProfileID,
	userID valueobjects.UserID,
	dni valueobjects.DNI,
	firstName string,
	firstLastName string,
	secondLastName string,
	phoneNumber valueobjects.PhoneNumber,
	monthlyIncome valueobjects.MonthlyIncome,
	maritalStatus valueobjects.MaritalStatus,
	isFirstHome bool,
	hasOwnLand bool,
	createdAt, updatedAt time.Time,
) *Profile {
	return &Profile{
		id:             id,
		userID:         userID,
		dni:            dni,
		firstName:      firstName,
		firstLastName:  firstLastName,
		secondLastName: secondLastName,
		phoneNumber:    phoneNumber,
		monthlyIncome:  monthlyIncome,
		maritalStatus:  maritalStatus,
		isFirstHome:    isFirstHome,
		hasOwnLand:     hasOwnLand,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// Getters
func (p *Profile) ID() valueobjects.ProfileID           { return p.id }
func (p *Profile) UserID() valueobjects.UserID          { return p.userID }
func (p *Profile) DNI() valueobjects.DNI                { return p.dni }
func (p *Profile) FirstName() string                    { return p.firstName }
func (p *Profile) FirstLastName() string                { return p.firstLastName }
func (p *Profile) SecondLastName() string               { return p.secondLastName }
func (p *Profile) PhoneNumber() valueobjects.PhoneNumber { return p.phoneNumber }
func (p *Profile) MonthlyIncome() valueobjects.MonthlyIncome { return p.monthlyIncome }
func (p *Profile) MaritalStatus() valueobjects.MaritalStatus { return p.maritalStatus }
func (p *Profile) IsFirstHome() bool                    { return p.isFirstHome }
func (p *Profile) HasOwnLand() bool                     { return p.hasOwnLand }
func (p *Profile) CreatedAt() time.Time                 { return p.createdAt }
func (p *Profile) UpdatedAt() time.Time                 { return p.updatedAt }

func (p *Profile) FullName() string {
	return p.firstLastName + " " + p.secondLastName + " " + p.firstName
}

func (p *Profile) SetID(id valueobjects.ProfileID) {
	p.id = id
}

func (p *Profile) UpdatePhoneNumber(phoneNumber valueobjects.PhoneNumber) {
	p.phoneNumber = phoneNumber
	p.updatedAt = time.Now()
}

func (p *Profile) UpdateMonthlyIncome(monthlyIncome valueobjects.MonthlyIncome) {
	p.monthlyIncome = monthlyIncome
	p.updatedAt = time.Now()
}

func (p *Profile) UpdateMaritalStatus(maritalStatus valueobjects.MaritalStatus) {
	p.maritalStatus = maritalStatus
	p.updatedAt = time.Now()
}

func (p *Profile) UpdateIsFirstHome(isFirstHome bool) {
	p.isFirstHome = isFirstHome
	p.updatedAt = time.Now()
}

func (p *Profile) UpdateHasOwnLand(hasOwnLand bool) {
	p.hasOwnLand = hasOwnLand
	p.updatedAt = time.Now()
}
