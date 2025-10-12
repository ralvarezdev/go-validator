package validator

import (
	"time"

	govalidatormapper "github.com/ralvarezdev/go-validator/mapper"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/mapper/validation"
)

type (
	// Service interface for the validator service
	Service interface {
		ValidateRequiredFields(
			rootStructValidations *govalidatormappervalidation.StructValidations,
			mapper *govalidatormapper.Mapper,
		) error
		ParseValidations(rootStructValidations *govalidatormappervalidation.StructValidations) (
			interface{},
			error,
		)
		Email(
			emailField string,
			email string,
			validations *govalidatormappervalidation.StructValidations,
		)
		Username(
			usernameField string,
			username string,
			validations *govalidatormappervalidation.StructValidations,
		)
		Birthdate(
			birthdateField string,
			birthdate time.Time,
			validations *govalidatormappervalidation.StructValidations,
		)
		Password(
			passwordField string,
			password string,
			validations *govalidatormappervalidation.StructValidations,
		)
		CreateValidateFn(
			mapper *govalidatormapper.Mapper,
			cache bool,
			auxiliaryValidatorFns ...interface{},
		) (
			ValidateFn, error,
		)
		Validate(
			mapper *govalidatormapper.Mapper,
			auxiliaryValidatorFns ...interface{},
		) (interface{}, error)
	}

	// Validator interface
	Validator interface {
		ValidateRequiredFields(
			rootStructValidations *govalidatormappervalidation.StructValidations,
			mapper *govalidatormapper.Mapper,
		) (err error)
	}
)
