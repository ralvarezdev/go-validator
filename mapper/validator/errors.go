package validator

import (
	"errors"
)

const (
	ErrFieldTagNameNotFound                   = "field tag name not found: %s"
	ErrFieldIsRequiredNotFound                = "field is required not found: %s"
	ErrStructValidationsAndMapperTypeMismatch = "struct validations and mapper type mismatch, both must be of the same type, mapper type: %s, struct validations type: %s"
)

var (
	ErrNilService                      = errors.New("mapper validator service cannot be nil")
	ErrNilDestination                  = errors.New("destination cannot be nil")
	ErrDestinationNotPointer           = errors.New("destination must be a pointer")
	ErrNilMapper                       = errors.New("mapper cannot be nil")
	ErrNilValidator                    = errors.New("mapper validator cannot be nil")
	ErrStructValidationsIsNotRootLevel = errors.New("struct validations is not root level")
	ErrRequiredField                   = "%s is required"
)
