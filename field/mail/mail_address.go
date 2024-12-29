package mail

import (
	"net/mail"
)

// ValidMailAddress checks if the mail address is valid
func ValidMailAddress(address string) (string, error) {
	// Check if the mail address is empty
	if address == "" {
		return "", ErrInvalidMailAddress
	}

	// Check if the mail address is valid
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", ErrInvalidMailAddress
	}

	return addr.Address, nil
}
