package service

import (
	"errors"
)

var (
	ErrValidations             = "validations error: %v"
	ErrFailedToGenerateMessage = errors.New("failed to generate message")
	ErrNilService              = errors.New("validator service cannot be nil")
	ErrNilValidations          = errors.New("validations cannot be nil")
)
