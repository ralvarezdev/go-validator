package json

import (
	"errors"
)

const (
	ErrNilFieldNameAlreadyParsed = "field name already parsed: %s"
)

var (
	ErrNilParsedValidations                    = errors.New("parsed validations is nil")
	ErrFlattenedParsedValidationsAlreadyExists = errors.New("flattened parsed validations already exists")
	ErrNilFlattenedParsedValidations           = errors.New("flattened parsed validations is nil")
)
