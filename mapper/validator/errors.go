package validator

import (
	"errors"
)

var (
	ErrNilService              = errors.New("mapper validator service cannot be nil")
	ErrNilDestination          = errors.New("destination cannot be nil")
	ErrDestinationNotPointer   = errors.New("destination must be a pointer")
	ErrNilMapper               = errors.New("mapper cannot be nil")
	ErrNilValidator            = errors.New("mapper validator cannot be nil")
	ErrFieldTagNameNotFound    = "field tag name not found: %s"
	ErrFieldIsRequiredNotFound = "field is required not found: %s"
	ErrRequiredField           = "%s is required"
)
