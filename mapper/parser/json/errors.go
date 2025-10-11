package json

import (
	"errors"
)

const (
	ErrFieldNameAlreadyParsed = "field name already parsed: %s"
)

var (
	ErrNilFlattenedParsedValidations = errors.New("flattened parsed validations is nil")
)
