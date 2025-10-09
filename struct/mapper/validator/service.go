package validator

import (
	"fmt"
	"log/slog"
	"net/mail"
	"reflect"
	"time"

	goreflect "github.com/ralvarezdev/go-reflect"
	gostringscount "github.com/ralvarezdev/go-strings/count"
	govalidatorfieldbirthdate "github.com/ralvarezdev/go-validator/struct/field/birthdate"
	govalidatorfieldmail "github.com/ralvarezdev/go-validator/struct/field/mail"
	govalidatorfieldpassword "github.com/ralvarezdev/go-validator/struct/field/password"
	govalidatorfieldusername "github.com/ralvarezdev/go-validator/struct/field/username"
	govalidatormapper "github.com/ralvarezdev/go-validator/struct/mapper"
	govalidatormapperparser "github.com/ralvarezdev/go-validator/struct/mapper/parser"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
)

type (
	// DefaultService struct
	DefaultService struct {
		parser      govalidatormapperparser.Parser
		validator   Validator
		validateFns map[string]ValidateFn
		logger      *slog.Logger
	}

	// BirthdateOptions is the birthdate options struct
	BirthdateOptions struct {
		MinimumAge int
		MaximumAge int
	}

	// PasswordOptions is the password options struct
	PasswordOptions struct {
		MinimumLength       int
		MinimumSpecialCount int
		MinimumNumbersCount int
		MinimumCapsCount    int
	}
)

// NewDefaultService creates a new default validator service
//
// Parameters:
//
//   - parser: the parser to use
//   - validator: the validator to use
//   - logger: the logger to use
//
// Returns:
//
//   - *DefaultService: the default validator service
//   - error: if there was an error creating the service
func NewDefaultService(
	parser govalidatormapperparser.Parser,
	validator Validator,
	logger *slog.Logger,
) (*DefaultService, error) {
	// Check if the parser or the validator is nil
	if parser == nil {
		return nil, govalidatormapperparser.ErrNilParser
	}
	if validator == nil {
		return nil, ErrNilValidator
	}

	if logger != nil {
		logger = logger.With(slog.String("component", "validator_service"))
	}

	return &DefaultService{
		parser,
		validator,
		nil,
		logger,
	}, nil
}

// ValidateRequiredFields validates the required fields
//
// Parameters:
//
//   - rootStructValidations: the root struct validations
//   - mapper: the mapper to use
//
// Returns:
//
// - error: if there was an error validating the required fields
func (d *DefaultService) ValidateRequiredFields(
	rootStructValidations *govalidatormappervalidation.StructValidations,
	mapper *govalidatormapper.Mapper,
) error {
	if d == nil {
		return ErrNilService
	}
	return d.validator.ValidateRequiredFields(
		rootStructValidations,
		mapper,
	)
}

// ParseValidations parses the validations
//
// Parameters:
//
//   - rootStructValidations: the root struct validations
//
// Returns:
//
//   - interface{}: the parsed validations
//   - error: if there was an error parsing the validations
func (d *DefaultService) ParseValidations(
	rootStructValidations *govalidatormappervalidation.StructValidations,
) (interface{}, error) {
	if d == nil {
		return nil, ErrNilService
	}

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

// Username validates the username field
//
// Parameters:
//
//   - usernameField: the username field name
//   - username: the username to validate
//   - validations: the struct validations
func (d *DefaultService) Username(
	usernameField string,
	username string,
	validations *govalidatormappervalidation.StructValidations,
) {
	if d == nil {
		return
	}

	// Check if the username contains non-alphanumeric characters
	if gostringscount.Alphanumeric(username) != len(username) {
		validations.AddFieldValidationError(
			usernameField,
			govalidatorfieldusername.ErrMustBeAlphanumeric,
		)
	}
}

// Email validates the email address field
//
// Parameters:
//
//   - emailField: the email field name
//   - email: the email to validate
//   - validations: the struct validations
func (d *DefaultService) Email(
	emailField string,
	email string,
	validations *govalidatormappervalidation.StructValidations,
) {
	if d == nil {
		return
	}

	// Check if the mail address is empty
	if email == "" {
		validations.AddFieldValidationError(
			emailField,
			govalidatorfieldmail.ErrInvalidMailAddress,
		)
	}

	// Check if the mail address is valid
	if _, err := mail.ParseAddress(email); err != nil {
		validations.AddFieldValidationError(
			emailField,
			govalidatorfieldmail.ErrInvalidMailAddress,
		)
	}
}

// Birthdate validates the birthdate field
//
// Parameters:
//
// - birthdateField: the birthdate field name
// - birthdate: the birthdate to validate
// - options: the birthdate options
// - validations: the struct validations
func (d *DefaultService) Birthdate(
	birthdateField string,
	birthdate time.Time,
	options *BirthdateOptions,
	validations *govalidatormappervalidation.StructValidations,
) {
	if d == nil {
		return
	}

	// Check if the birthdate is after the current time
	if birthdate.After(time.Now()) {
		validations.AddFieldValidationError(
			birthdateField,
			govalidatorfieldbirthdate.ErrInvalidBirthdate,
		)
	}

	// Check if the birthdate options are nil
	if options == nil {
		return
	}

	// Check if the birthdate is before the minimum age
	if options.MinimumAge > 0 {
		if time.Now().AddDate(-options.MinimumAge, 0, 0).Before(birthdate) {
			validations.AddFieldValidationError(
				birthdateField,
				fmt.Errorf(
					govalidatorfieldbirthdate.ErrMinimumAge,
					options.MinimumAge,
				),
			)
		}
	}

	// Check if the birthdate is after the maximum age
	if options.MaximumAge > 0 {
		if time.Now().AddDate(-options.MaximumAge, 0, 0).After(birthdate) {
			validations.AddFieldValidationError(
				birthdateField,
				fmt.Errorf(
					govalidatorfieldbirthdate.ErrMaximumAge,
					options.MaximumAge,
				),
			)
		}
	}
}

// Password validates the password field
//
// Parameters:
//
// - passwordField: the password field name
// - password: the password to validate
// - options: the password options
// - validations: the struct validations
func (d *DefaultService) Password(
	passwordField string,
	password string,
	options *PasswordOptions,
	validations *govalidatormappervalidation.StructValidations,
) {
	if d == nil {
		return
	}

	// Check if the password length is less than the minimum length
	if options.MinimumLength > 0 && len(password) < options.MinimumLength {
		validations.AddFieldValidationError(
			passwordField,
			fmt.Errorf(
				govalidatorfieldpassword.ErrMinimumLength,
				options.MinimumLength,
			),
		)
	}

	// Check if the password contains the minimum special characters
	if options.MinimumSpecialCount > 0 {
		if count := gostringscount.Special(password); count < options.MinimumSpecialCount {
			validations.AddFieldValidationError(
				passwordField,
				fmt.Errorf(
					govalidatorfieldpassword.ErrMinimumSpecialCount,
					options.MinimumSpecialCount,
				),
			)
		}
	}

	// Check if the password contains the minimum numbers
	if options.MinimumNumbersCount > 0 {
		if count := gostringscount.Numbers(password); count < options.MinimumNumbersCount {
			validations.AddFieldValidationError(
				passwordField,
				fmt.Errorf(
					govalidatorfieldpassword.ErrMinimumNumbersCount,
					options.MinimumNumbersCount,
				),
			)
		}
	}

	// Check if the password contains the minimum caps
	if options.MinimumCapsCount > 0 {
		if count := gostringscount.Caps(password); count < options.MinimumCapsCount {
			validations.AddFieldValidationError(
				passwordField,
				fmt.Errorf(
					govalidatorfieldpassword.ErrMinimumCapsCount,
					options.MinimumCapsCount,
				),
			)
		}
	}
}

// CreateValidateFn creates a validate function for a given mapper
//
// Parameters:
//
//   - mapper: the mapper to use
//   - cache: whether to cache the validate function or not
//   - auxiliaryValidatorFns: the auxiliary validator functions to use
//
// Returns:
//
//   - ValidateFn: the validate function
//   - error: if there was an error creating the validate function
func (d *DefaultService) CreateValidateFn(
	mapper *govalidatormapper.Mapper,
	cache bool,
	auxiliaryValidatorFns ...AuxiliaryValidatorFn,
) (
	ValidateFn, error,
) {
	if d == nil {
		return nil, ErrNilService
	}

	// Check if the mapper is nil
	if mapper == nil {
		return nil, govalidatormapper.ErrNilMapper
	}

	// Check if the cache parameter is true, if so call try to get the validate function from the cache
	if cache {
		if validateFn, ok := d.validateFns[goreflect.UniqueTypeReference(mapper.GetStructInstance())]; ok {
			return validateFn, nil
		}
	}

	// Create the validate function
	validateFn := func(
		toValidate interface{},
	) (
		interface{},
		error,
	) {
		// Check if the destination is a pointer
		if toValidate == nil {
			return nil, ErrNilDestination
		}
		if reflect.TypeOf(toValidate).Kind() != reflect.Ptr {
			return nil, ErrDestinationNotPointer
		}

		// Initialize struct fields validations from the request body
		rootStructValidations, err := govalidatormappervalidation.NewStructValidations(toValidate)
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

		// Call the validate function
		for _, auxiliaryValidatorFn := range auxiliaryValidatorFns {
			_, err = goreflect.SafeCallFunction(
				auxiliaryValidatorFn,
				toValidate,
				rootStructValidations,
			)
			if err != nil {
				if d.logger != nil {
					d.logger.Error(
						"error calling auxiliary validator function",
						slog.String("error", err.Error()),
					)
				}
				return nil, err
			}
		}

		// Parse the validations
		return d.ParseValidations(rootStructValidations)
	}

	// If cache is true, cache the validate function
	if cache {
		if d.validateFns == nil {
			d.validateFns = make(map[string]ValidateFn)
		}
		d.validateFns[goreflect.UniqueTypeReference(mapper.GetStructInstance())] = validateFn
	}

	return validateFn, nil
}

// Validate is the function that creates (if not cached), caches and executes the validation
//
// Parameters:
//
//   - mapper: the mapper to use
//   - auxiliaryValidatorFns: auxiliary validator functions to use in the validation
//
// Returns:
//
//   - interface{}: the parsed validations
//   - error: if there was an error validating the request
func (d *DefaultService) Validate(
	mapper *govalidatormapper.Mapper,
	auxiliaryValidatorFns ...AuxiliaryValidatorFn,
) (interface{}, error) {
	if d == nil {
		return nil, ErrNilService
	}

	// Create and cache the validate function
	validateFn, err := d.CreateValidateFn(
		mapper,
		true,
		auxiliaryValidatorFns...,
	)
	if err != nil {
		return nil, err
	}

	// Execute the validate function
	return validateFn(mapper)
}
