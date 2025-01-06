package validations

import (
	"errors"
)

var (
	ErrNilStructValidations = errors.New("struct validations is nil")
	ErrNilFieldValidations  = errors.New("field validations is nil")
)
