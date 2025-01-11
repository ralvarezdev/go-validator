package validator

import (
	"errors"
)

var (
	ErrNilMapper               = errors.New("mapper cannot be nil")
	ErrNilValidator            = errors.New("mapper validator cannot be nil")
	ErrFieldTagNameNotFound    = "field tag name not found: %s"
	ErrFieldIsRequiredNotFound = "field is required not found: %s"
	ErrRequiredField           = "%s is required"
)
