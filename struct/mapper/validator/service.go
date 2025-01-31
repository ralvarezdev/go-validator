package validator

import (
	govalidatormapper "github.com/ralvarezdev/go-validator/struct/mapper"
	govalidatormapperparser "github.com/ralvarezdev/go-validator/struct/mapper/parser"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
	"reflect"
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
		CreateValidateFn(
			mapper *govalidatormapper.Mapper,
			validationsFns ...func(*govalidatormappervalidation.StructValidations) error,
		) func(
			dest interface{},
		) (interface{}, error)
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

// CreateValidateFn creates a validate function for the request body using the validator
// functions provided. It validates the required fields by default
func (d *DefaultService) CreateValidateFn(
	mapper *govalidatormapper.Mapper,
	validationsFns ...func(*govalidatormappervalidation.StructValidations) error,
) func(
	dest interface{},
) (interface{}, error) {
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
	}
}
