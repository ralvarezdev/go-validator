package parser

import (
	"errors"
)

var (
	ErrNilJSONParsedValidations = errors.New("json parsed validations is nil")
	ErrNilParser                = errors.New("mapper validations parser cannot be nil")
)
