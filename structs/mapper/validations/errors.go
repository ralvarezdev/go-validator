package validations

import (
	"errors"
)

var (
	NilDataError        = errors.New("data cannot be nil")
	NilMapperError      = errors.New("mapper cannot be nil")
	NilGeneratorError   = errors.New("mapper validations generator cannot be nil")
	NilValidatorError   = errors.New("mapper validator cannot be nil")
	NilValidationsError = errors.New("mapper validations cannot be nil")
	FieldNotFoundError  = errors.New("field not found")
)
