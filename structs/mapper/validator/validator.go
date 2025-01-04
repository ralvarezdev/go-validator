package validator

import (
	"fmt"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	"reflect"
)

type (
	// Validator interface
	Validator interface {
		ValidateNilFields(
			data interface{},
			mapper *govalidatormapper.Mapper,
		) (validations govalidatorvalidations.Validations, err error)
	}

	// DefaultValidator struct
	DefaultValidator struct {
		mode *goflagsmode.Flag
	}
)

// NewDefaultValidator creates a new default mapper validator
func NewDefaultValidator(
	mode *goflagsmode.Flag,
) *DefaultValidator {
	return &DefaultValidator{
		mode: mode,
	}
}

// ValidateNilFields validates if the fields are not nil
func (d *DefaultValidator) ValidateNilFields(
	data interface{},
	mapper *govalidatormapper.Mapper,
) (validations govalidatorvalidations.Validations, err error) {
	// Check if either the data or the struct fields to validate are nil
	if data == nil {
		return nil, govalidatorvalidations.ErrNilData
	}
	if mapper == nil {
		return nil, govalidatorvalidations.ErrNilMapper
	}

	// Initialize struct fields validations
	validations = govalidatorvalidations.NewDefaultValidations()

	// Reflection of data
	valueReflection := reflect.ValueOf(data)

	// If data is a pointer, dereference it
	if valueReflection.Kind() == reflect.Ptr {
		valueReflection = valueReflection.Elem()
	}

	// Iterate over the fields
	fields := (*mapper).Fields
	nestedMappers := (*mapper).NestedMappers

	// Check if the struct has fields to validate
	if fields == nil && nestedMappers == nil {
		return nil, nil
	}

	// Iterate over the fields
	typeReflection := valueReflection.Type()
	for i := 0; i < valueReflection.NumField(); i++ {
		fieldValue := valueReflection.Field(i)
		fieldType := typeReflection.Field(i)

		// Print field on debug mode
		if d.mode != nil && d.mode.IsDebug() {
			fmt.Printf(
				"field '%v' of type '%v' value: %v\n",
				fieldType.Name,
				fieldType.Type,
				fieldValue,
			)
		}

		// Check if the field is a pointer
		if fieldValue.Kind() != reflect.Ptr {
			// Check if the field has to be validated
			if fields == nil {
				continue
			}
			validationName, ok := fields[fieldType.Name]
			if !ok {
				continue
			}

			// Check if the field is uninitialized
			if fieldValue.IsZero() {
				// Print error on debug mode
				if d.mode != nil && d.mode.IsDebug() {
					fmt.Printf("field is uninitialized: %v\n", fieldType.Name)
				}
				validations.AddFailedFieldValidationError(
					validationName,
					govalidatorvalidations.ErrFieldNotFound,
				)
			}
			continue
		}

		// Check if the field is a nested struct
		if fieldValue.Elem().Kind() != reflect.Struct {
			continue // It's an optional field
		}

		// Check if the nested struct has to be validated
		if fields == nil {
			continue
		}
		validationName, ok := fields[fieldType.Name]
		if !ok {
			continue
		}

		// Check if the field is initialized
		if fieldValue.IsNil() {
			// Print error on dev mode
			if d.mode != nil && d.mode.IsDev() {
				fmt.Printf("field is uninitialized: %v\n", fieldType.Name)
			}
			validations.AddFailedFieldValidationError(
				validationName,
				govalidatorvalidations.ErrFieldNotFound,
			)
			continue
		}

		// Get the nested struct fields to validate
		fieldNestedMapper, ok := nestedMappers[fieldType.Name]
		if !ok {
			continue
		}

		// Validate nested struct
		fieldNestedMapperValidations, err := d.ValidateNilFields(
			fieldValue.Addr().Interface(), // TEST IF THIS A POINTER OF THE STRUCT
			fieldNestedMapper,
		)
		if err != nil {
			return nil, err
		}

		// Add nested struct validations to the struct fields validations
		validations.SetNestedFieldsValidations(
			validationName,
			fieldNestedMapperValidations,
		)
	}

	return validations, nil
}
