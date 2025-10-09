package validator

import (
	"time"

	govalidatormapper "github.com/ralvarezdev/go-validator/struct/mapper"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
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
			options *BirthdateOptions,
			validations *govalidatormappervalidation.StructValidations,
		)
		Password(
			passwordField string,
			password string,
			options *PasswordOptions,
			validations *govalidatormappervalidation.StructValidations,
		)
		CreateValidateFn(
			mapper *govalidatormapper.Mapper,
			auxiliaryValidatorFns ...AuxiliaryValidatorFn,
		) (
			ValidateFn, error,
		)
		CreateAndCacheValidateFnFromMapper(
			mapper *govalidatormapper.Mapper,
			auxiliaryValidatorFns ...AuxiliaryValidatorFn,
		) (ValidateFn, error)
		Validate(
			mapper *govalidatormapper.Mapper,
			auxiliaryValidatorFns ...AuxiliaryValidatorFn,
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
