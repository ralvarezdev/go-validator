package validator

import (
	"errors"
)

var (
	ErrNilMapper     = errors.New("mapper cannot be nil")
	ErrNilValidator  = errors.New("mapper validator cannot be nil")
	ErrRequiredField = "%s is required"
)
