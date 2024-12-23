package service

import (
	"fmt"
	goflagmode "github.com/ralvarezdev/go-flags/mode"
	govalidatorbirthdate "github.com/ralvarezdev/go-validator/field/birthdate"
	govalidatormail "github.com/ralvarezdev/go-validator/field/mail"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type (
	// Service interface for the validator service
	Service interface {
		ModeFlag() *goflagmode.Flag
		ValidateEmail(
			emailField string,
			email string,
			validations *govalidatorvalidations.Validations,
		)
		ValidateBirthdate(
			birthdateField string,
			birthdate *timestamppb.Timestamp,
			validations *govalidatorvalidations.Validations,
		)
		ValidateNilFields(
			request interface{},
			mapper *govalidatormapper.Mapper,
		) (
			*govalidatorvalidations.Validations,
			error,
		)
		CheckValidations(validations *govalidatorvalidations.Validations) error
	}

	// DefaultValidator struct
	DefaultValidator struct {
		mode      *goflagmode.Flag
		generator *govalidatorvalidations.Generator
		validator *govalidatorvalidations.Validator
	}
)

// NewDefaultValidator creates a new default validator
func NewDefaultValidator(
	generator *govalidatorvalidations.Generator,
	validator *govalidatorvalidations.Validator,
	mode *goflagmode.Flag,
) (*DefaultValidator, error) {
	// Check if the generator or the validator is nil
	if generator == nil {
		return nil, govalidatorvalidations.NilGeneratorError
	}
	if validator == nil {
		return nil, govalidatorvalidations.NilValidatorError
	}

	return &DefaultValidator{
		mode:      mode,
		generator: generator,
		validator: validator,
	}, nil
}

// ModeFlag returns the mode flag
func (d *DefaultValidator) ModeFlag() *goflagmode.Flag {
	return d.mode
}

// ValidateEmail validates the email address field
func (d *DefaultValidator) ValidateEmail(
	emailField string,
	email string,
	validations *govalidatorvalidations.Validations,
) {
	if _, err := govalidatormail.ValidMailAddress(email); err != nil {
		(*validations).AddFailedFieldValidationError(
			emailField,
			govalidatormail.InvalidMailAddressError,
		)
	}
}

// ValidateBirthdate validates the birthdate field
func (d *DefaultValidator) ValidateBirthdate(
	birthdateField string,
	birthdate *timestamppb.Timestamp,
	validations *govalidatorvalidations.Validations,
) {
	if birthdate == nil || birthdate.AsTime().After(time.Now()) {
		(*validations).AddFailedFieldValidationError(
			birthdateField,
			govalidatorbirthdate.InvalidBirthdateError,
		)
	}
}

// ValidateNilFields validates the nil fields
func (d *DefaultValidator) ValidateNilFields(
	request interface{},
	mapper *govalidatormapper.Mapper,
) (*govalidatorvalidations.Validations, error) {
	return d.validator.ValidateNilFields(
		request,
		mapper,
		d.mode,
	)
}

// CheckValidations checks the validations and returns a pointer to the error message
func (d *DefaultValidator) CheckValidations(
	validations *govalidatorvalidations.Validations,
) error {
	// Get the error message from the validations if there are any
	if !(*validations).HasFailed() {
		return nil
	}

	// Get the validations message
	message, err := d.generator.Generate(validations)
	if err != nil {
		return FailedToGenerateMessageError
	}

	if message != nil {
		return fmt.Errorf(ValidationsError, *message)
	}
	return NilValidationsError
}
