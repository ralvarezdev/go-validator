package mapper

import (
	"errors"
)

var (
	ErrNilGenerator            = errors.New("generator cannot be nil")
	ErrProtobufTagNotFound     = "missing protobuf tag: %s"
	ErrProtobufTagNameNotFound = "missing protobuf tag name: %s"
	ErrEmptyJSONTag            = "empty json tag: %s"
)
