package validations

import (
	"errors"
)

var (
	NilDataError       = errors.New("data cannot be nil")
	NilMapperError     = errors.New("mapper cannot be nil")
	FieldNotFoundError = errors.New("field not found")
)
