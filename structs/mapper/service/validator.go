package service

import (
	"fmt"
	goflagmode "github.com/ralvarezdev/go-flags/mode"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
)

type (
	// Service interface for the validator service
	Service interface {
		ModeFlag() *goflagmode.Flag
		ValidateNilFields(
			request interface{},
			mapper *govalidatormapper.Mapper,
		) (
			govalidatorvalidations.Validations,
			error,
		)
		CheckValidations(validations govalidatorvalidations.Validations) error
	}

	// DefaultService struct
	DefaultService struct {
		mode      *goflagmode.Flag
		generator govalidatorvalidations.Generator
		validator govalidatorvalidations.Validator
	}
)

// NewDefaultService creates a new default validator service
func NewDefaultService(
	generator govalidatorvalidations.Generator,
	validator govalidatorvalidations.Validator,
	mode *goflagmode.Flag,
) (*DefaultService, error) {
	// Check if the generator or the validator is nil
	if generator == nil {
		return nil, govalidatorvalidations.ErrNilGenerator
	}
	if validator == nil {
		return nil, govalidatorvalidations.ErrNilValidator
	}

	return &DefaultService{
		mode:      mode,
		generator: generator,
		validator: validator,
	}, nil
}

// ModeFlag returns the mode flag
func (d *DefaultService) ModeFlag() *goflagmode.Flag {
	return d.mode
}

// ValidateNilFields validates the nil fields
func (d *DefaultService) ValidateNilFields(
	request interface{},
	mapper *govalidatormapper.Mapper,
) (govalidatorvalidations.Validations, error) {
	return d.validator.ValidateNilFields(
		request,
		mapper,
	)
}

// CheckValidations checks the validations and returns a pointer to the error message
func (d *DefaultService) CheckValidations(
	validations govalidatorvalidations.Validations,
) error {
	// Get the error message from the validations if there are any
	if !validations.HasFailed() {
		return nil
	}

	// Get the validations message
	message, err := d.generator.Generate(validations)
	if err != nil {
		return ErrFailedToGenerateMessage
	}

	if message != nil {
		return fmt.Errorf(ErrValidations, *message)
	}
	return ErrNilValidations
}
