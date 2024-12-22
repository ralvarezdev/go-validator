package service

import (
	"errors"
	goflagmode "github.com/ralvarezdev/go-flags/mode"
	govalidatorbirthdate "github.com/ralvarezdev/go-validator/field/birthdate"
	govalidatormail "github.com/ralvarezdev/go-validator/field/mail"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/validations"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type (
	// Validator interface
	Validator interface {
		ModeFlag() *goflagmode.Flag
		ValidateEmail(
			emailField string,
			email string,
			mapperValidations *govalidatorvalidations.MapperValidations,
		)
		ValidateBirthdate(
			birthdateField string,
			birthdate *timestamppb.Timestamp,
			mapperValidations *govalidatorvalidations.MapperValidations,
		)
		ValidateNilFields(request interface{}, mapper *govalidatormapper.Mapper) (
			*govalidatorvalidations.MapperValidations,
			error,
		)
		CheckValidations(mapperValidations *govalidatorvalidations.MapperValidations) error
	}

	// DefaultValidator struct
	DefaultValidator struct {
		mode *goflagmode.Flag
	}
)

// NewDefaultValidator creates a new default validator
func NewDefaultValidator(mode *goflagmode.Flag) *DefaultValidator {
	return &DefaultValidator{
		mode: mode,
	}
}

// ModeFlag returns the mode flag
func (d *DefaultValidator) ModeFlag() *goflagmode.Flag {
	return d.mode
}

// ValidateEmail validates the email address field
func (d *DefaultValidator) ValidateEmail(
	emailField string,
	email string,
	mapperValidations *govalidatorvalidations.MapperValidations,
) {
	if _, err := govalidatormail.ValidMailAddress(email); err != nil {
		mapperValidations.AddFailedFieldValidationError(
			emailField,
			govalidatormail.InvalidMailAddressError,
		)
	}
}

// ValidateBirthdate validates the birthdate field
func (d *DefaultValidator) ValidateBirthdate(
	birthdateField string,
	birthdate *timestamppb.Timestamp,
	mapperValidations *govalidatorvalidations.MapperValidations,
) {
	if birthdate == nil || birthdate.AsTime().After(time.Now()) {
		mapperValidations.AddFailedFieldValidationError(
			birthdateField,
			govalidatorbirthdate.InvalidBirthdateError,
		)
	}
}

// ValidateNilFields validates the nil fields
func (d *DefaultValidator) ValidateNilFields(
	request interface{},
	mapper *govalidatormapper.Mapper,
) (*govalidatorvalidations.MapperValidations, error) {
	return govalidatorvalidations.ValidateMapperNilFields(
		request,
		mapper,
		d.mode,
	)
}

// CheckValidations checks the validations and returns a pointer to the error message
func (d *DefaultValidator) CheckValidations(
	mapperValidations *govalidatorvalidations.MapperValidations,
) error {
	// Get the error message from the validations if there are any
	if mapperValidations.HasFailed() {
		// Get the validations
		validations := mapperValidations.StringPtr()

		if validations != nil {
			return errors.New(*validations)
		}
		return NilValidationsError
	}
	return nil
}
