package validations

import (
	"errors"
)

var (
	ErrNilStructData        = errors.New("struct data cannot be nil")
	ErrNilStructValidations = errors.New("struct validations is nil")
	ErrNilFieldValidations  = errors.New("field validations is nil")
)
