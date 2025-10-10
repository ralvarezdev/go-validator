package json

import (
	"errors"
)

const (
	ErrNilFieldNameAlreadyParsed = "field name already parsed: %s"
)

var (
	ErrFlattenedParsedValidationsAlreadyExists = errors.New("flattened parsed validations already exists")
	ErrNilFlattenedParsedValidations           = errors.New("flattened parsed validations is nil")
)
