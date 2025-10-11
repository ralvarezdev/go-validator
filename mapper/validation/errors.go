package validation

import (
	"errors"
)

var (
	ErrNilInstance          = errors.New("struct instance cannot be nil")
	ErrNilStructValidations = errors.New("struct validations is nil")
	ErrNilFieldValidations  = errors.New("field validations is nil")
)
