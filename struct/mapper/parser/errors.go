package parser

import (
	"errors"
)

var (
	ErrNilRawParser               = errors.New("mapper validations raw parser cannot be nil")
	ErrNilEndParser               = errors.New("mapper validations end parser cannot be nil")
	ErrNilStructParsedValidations = errors.New("struct parsed validations is nil")
)
