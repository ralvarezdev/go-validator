package parser

import (
	"errors"
)

const (
	ErrNilFieldNameAlreadyParsed = "field name already parsed: %s"
)

var (
	ErrNilParser                               = errors.New("mapper validations parser cannot be nil")
	ErrNilParsedValidations                    = errors.New("parsed validations is nil")
	ErrFlattenedParsedValidationsAlreadyExists = errors.New("flattened parsed validations already exists")
	ErrNilFlattenedParsedValidations           = errors.New("flattened parsed validations is nil")
)
