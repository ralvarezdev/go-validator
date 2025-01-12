package json

import (
	"errors"
)

var (
	ErrNilParsedValidations                    = errors.New("parsed validations is nil")
	ErrNilFieldNameAlreadyParsed               = "field name already parsed: %s"
	ErrFlattenedParsedValidationsAlreadyExists = errors.New("flattened parsed validations already exists")
)
