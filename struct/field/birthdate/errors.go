package birthdate

import (
	"errors"
)

var (
	ErrInvalidBirthdate = errors.New("invalid birthdate")
	ErrMinimumAge       = "age must be greater than or equal to %d"
	ErrMaximumAge       = "age must be less than or equal to %d"
)
