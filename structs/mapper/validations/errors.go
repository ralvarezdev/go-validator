package validations

import (
	"errors"
)

var (
	ErrNilData        = errors.New("data cannot be nil")
	ErrNilMapper      = errors.New("mapper cannot be nil")
	ErrNilGenerator   = errors.New("mapper validations generator cannot be nil")
	ErrNilValidator   = errors.New("mapper validator cannot be nil")
	ErrNilValidations = errors.New("mapper validations cannot be nil")
	ErrFieldNotFound  = errors.New("field not found")
)
