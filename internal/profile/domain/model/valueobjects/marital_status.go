package valueobjects

import "errors"

type MaritalStatus string

const (
	MaritalStatusSingle   MaritalStatus = "SOLTERO"
	MaritalStatusMarried  MaritalStatus = "CASADO"
	MaritalStatusDivorced MaritalStatus = "DIVORCIADO"
	MaritalStatusWidowed  MaritalStatus = "VIUDO"
)

func NewMaritalStatus(value string) (MaritalStatus, error) {
	status := MaritalStatus(value)
	switch status {
	case MaritalStatusSingle, MaritalStatusMarried, MaritalStatusDivorced, MaritalStatusWidowed:
		return status, nil
	default:
		return "", errors.New("invalid marital status")
	}
}

func (m MaritalStatus) String() string {
	return string(m)
}

func GetAllMaritalStatuses() []MaritalStatus {
	return []MaritalStatus{
		MaritalStatusSingle,
		MaritalStatusMarried,
		MaritalStatusDivorced,
		MaritalStatusWidowed,
	}
}
