package valueobjects

import (
	"errors"
	"regexp"
)

type DNI struct {
	value          string // Plaintext DNI (for domain logic)
	encryptedValue string // Encrypted DNI (for persistence)
}

var dniRegex = regexp.MustCompile(`^\d{8}$`)

// NewDNI creates a new DNI from plaintext value
func NewDNI(value string) (DNI, error) {
	if value == "" {
		return DNI{}, errors.New("DNI cannot be empty")
	}
	if !dniRegex.MatchString(value) {
		return DNI{}, errors.New("DNI must be exactly 8 digits")
	}
	return DNI{value: value}, nil
}

// NewDNIFromEncrypted creates a DNI from encrypted value (for reconstruction)
func NewDNIFromEncrypted(encryptedValue string) DNI {
	return DNI{encryptedValue: encryptedValue}
}

// Value returns the plaintext DNI
func (d DNI) Value() string {
	return d.value
}

// EncryptedValue returns the encrypted DNI for persistence
func (d DNI) EncryptedValue() string {
	return d.encryptedValue
}

// SetEncryptedValue sets the encrypted value (used by repository)
func (d *DNI) SetEncryptedValue(encrypted string) {
	d.encryptedValue = encrypted
}

// SetPlainValue sets the plaintext value (used by repository after decryption)
func (d *DNI) SetPlainValue(plain string) {
	d.value = plain
}
