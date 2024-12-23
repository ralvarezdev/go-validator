package service

import (
	"errors"
)

var (
	ValidationsError             = "validations error: %v"
	FailedToGenerateMessageError = errors.New("failed to generate message")
	NilServiceError              = errors.New("validator service cannot be nil")
	NilValidationsError          = errors.New("validations cannot be nil")
)
