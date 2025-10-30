package mapper

import (
	"errors"
)

var (
	ErrNilGenerator            = errors.New("generator cannot be nil")
	ErrNilMapper               = errors.New("mapper cannot be nil")
	ErrProtobufTagNotFound     = "missing protobuf tag for field: %s"
	ErrProtobufTagNameNotFound = "missing protobuf tag name for field: %s"
	ErrEmptyJSONTag            = "empty json tag for field: %s"
	ErrNilStructInstance      = errors.New("struct instance cannot be nil")
	ErrStructInstanceNotStruct = errors.New("struct instance must be a struct")
	ErrInvalidStructInstance = errors.New("invalid struct instance")
)
