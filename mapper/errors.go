package mapper

import (
	"errors"
)

var (
	ErrNilGenerator            = errors.New("generator cannot be nil")
	ErrNilMapper               = errors.New("mapper cannot be nil")
	ErrProtobufTagNotFound     = "missing protobuf tag: %s"
	ErrProtobufTagNameNotFound = "missing protobuf tag name: %s"
	ErrEmptyJSONTag            = "empty json tag: %s"
	ErrNilStructInstance      = errors.New("struct instance cannot be nil")
	ErrStructInstanceNotStruct = errors.New("struct instance must be a struct")
	ErrInvalidStructInstance = errors.New("invalid struct instance")
)
