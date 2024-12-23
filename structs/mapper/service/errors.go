package service

import (
	"errors"
)

var (
	ValidationsError             = "validations error: %v"
	FailedToGenerateMessageError = errors.New("failed to generate message")
	NilValidatorError            = errors.New("validator cannot be nil")
	NilMessageError              = errors.New("message cannot be nil")
	NilValidationsError          = errors.New("validations cannot be nil")
)
