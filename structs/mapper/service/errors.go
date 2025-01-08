package service

import (
	"errors"
)

var (
	ErrNilService      = errors.New("mapper validator service cannot be nil")
	ErrInvalidBodyType = "invalid body type: expected '%T'"
)
