package validator

import (
	govalidatorfieldbirthdate "github.com/ralvarezdev/go-validator/struct/field/birthdate"
	govalidatorfieldmail "github.com/ralvarezdev/go-validator/struct/field/mail"
	govalidatormapper "github.com/ralvarezdev/go-validator/struct/mapper"
	govalidatormapperparser "github.com/ralvarezdev/go-validator/struct/mapper/parser"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
	"reflect"
	"time"
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
		) error
		Birthdate(
			birthdateField string,
			birthdate *time.Time,
			validations *govalidatormappervalidation.StructValidations,
		) error
		CreateValidateFn(
			mapper *govalidatormapper.Mapper,
			validationsFns ...func(*govalidatormappervalidation.StructValidations) error,
		) (
			func(
				dest interface{},
			) (interface{}, error), error,
		)
	}

	// DefaultService struct
	DefaultService struct {
		parser    govalidatormapperparser.Parser
		validator Validator
	}
)

// NewDefaultService creates a new default validator service
func NewDefaultService(
	parser govalidatormapperparser.Parser,
	validator Validator,
) (*DefaultService, error) {
	// Check if the parser or the validator is nil
	if parser == nil {
		return nil, govalidatormapperparser.ErrNilParser
	}
	if validator == nil {
		return nil, ErrNilValidator
	}

	return &DefaultService{
		parser,
		validator,
	}, nil
}

// ValidateRequiredFields validates the required fields
func (d *DefaultService) ValidateRequiredFields(
	rootStructValidations *govalidatormappervalidation.StructValidations,
	mapper *govalidatormapper.Mapper,
) error {
	return d.validator.ValidateRequiredFields(
		rootStructValidations,
		mapper,
	)
}

// ParseValidations parses the validations
func (d *DefaultService) ParseValidations(
	rootStructValidations *govalidatormappervalidation.StructValidations,
) (interface{}, error) {
	// Check if there are any failed validations
	if !rootStructValidations.HasFailed() {
		return nil, nil
	}

	// Get the parsed validations from the validations
	parsedValidations, err := d.parser.ParseValidations(rootStructValidations)
	if err != nil {
		return nil, err
	}
	return parsedValidations, nil
}

// Email validates the email address field
func (d *DefaultService) Email(
	emailField string,
	email string,
	validations *govalidatormappervalidation.StructValidations,
) error {
	if _, err := govalidatorfieldmail.ValidMailAddress(email); err != nil {
		validations.AddFieldValidationError(
			emailField,
			govalidatorfieldmail.ErrInvalidMailAddress,
		)
	}
	return nil
}

// Birthdate validates the birthdate field
func (d *DefaultService) Birthdate(
	birthdateField string,
	birthdate *time.Time,
	validations *govalidatormappervalidation.StructValidations,
) error {
	if birthdate == nil || birthdate.After(time.Now()) {
		validations.AddFieldValidationError(
			birthdateField,
			govalidatorfieldbirthdate.ErrInvalidBirthdate,
		)
	}
	return nil
}

// CreateValidateFn creates a validate function for the request body using the validator
// functions provided. It validates the required fields by default
func (d *DefaultService) CreateValidateFn(
	mapper *govalidatormapper.Mapper,
	validationsFns ...func(*govalidatormappervalidation.StructValidations) error,
) (
	func(
		dest interface{},
	) (interface{}, error), error,
) {
	// Check if the mapper is nil
	if mapper == nil {
		return nil, govalidatormapper.ErrNilMapper
	}

	return func(
		dest interface{},
	) (
		interface{},
		error,
	) {
		// Check if the destination is a pointer
		if dest == nil {
			return nil, ErrNilDestination
		}
		if reflect.TypeOf(dest).Kind() != reflect.Ptr {
			return nil, ErrDestinationNotPointer
		}

		// Initialize struct fields validations from the request body
		rootStructValidations, err := govalidatormappervalidation.NewStructValidations(dest)
		if err != nil {
			return nil, err
		}

		// Validate the required fields
		if err = d.ValidateRequiredFields(
			rootStructValidations,
			mapper,
		); err != nil {
			return nil, err
		}

		// Run the validator functions
		for _, validatorFn := range validationsFns {
			if err = validatorFn(rootStructValidations); err != nil {
				return nil, err
			}
		}

		// Parse the validations
		return d.ParseValidations(rootStructValidations)
	}, nil
}
