package validator

import (
	"errors"
)

var (
	ErrNilService              = errors.New("mapper validator service cannot be nil")
	ErrNilMapper               = errors.New("mapper cannot be nil")
	ErrNilValidator            = errors.New("mapper validator cannot be nil")
	ErrFieldTagNameNotFound    = "field tag name not found: %s"
	ErrFieldIsRequiredNotFound = "field is required not found: %s"
	ErrRequiredField           = "%s is required"
)
