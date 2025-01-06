package validator

import (
	"fmt"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatormappervalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	"reflect"
)

type (
	// Validator interface
	Validator interface {
		ValidateRequiredFields(
			rootStructValidations *govalidatormappervalidations.StructValidations,
			data interface{},
			mapper *govalidatormapper.Mapper,
		) (err error)
	}

	// DefaultValidator struct
	DefaultValidator struct {
		logger *Logger
	}
)

// NewDefaultValidator creates a new default mapper validator
func NewDefaultValidator(
	logger *Logger,
) *DefaultValidator {
	return &DefaultValidator{
		logger: logger,
	}
}

// ValidateRequiredFields validates the required fields of a struct
func (d *DefaultValidator) ValidateRequiredFields(
	rootStructValidations *govalidatormappervalidations.StructValidations,
	data interface{},
	mapper *govalidatormapper.Mapper,
) (err error) {
	// Check if either the root struct validations, data or the struct fields to validate are nil
	if rootStructValidations == nil {
		return govalidatormappervalidations.ErrNilStructValidations
	}
	if data == nil {
		return ErrNilData
	}
	if mapper == nil {
		return ErrNilMapper
	}

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
		return nil
	}

	// Iterate over the fields
	typeReflection := valueReflection.Type()
	for i := 0; i < valueReflection.NumField(); i++ {
		fieldValue := valueReflection.Field(i)
		fieldType := typeReflection.Field(i)

		// Print field
		if d.logger != nil {
			d.logger.PrintField(fieldType.Name, fieldType.Type, fieldValue)
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
				if d.logger != nil {
					d.logger.UninitializedField(fieldType.Name)
				}

				rootStructValidations.AddFieldValidationError(
					validationName,
					fmt.Errorf(ErrRequiredField, validationName),
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
			if d.logger != nil {
				d.logger.UninitializedField(fieldType.Name)
			}

			rootStructValidations.AddFieldValidationError(
				validationName,
				fmt.Errorf(ErrRequiredField, validationName),
			)
			continue
		}

		// Get the nested struct mapper
		fieldNestedMapper, ok := nestedMappers[fieldType.Name]
		if !ok {
			continue
		}

		// Initialize the nested struct mapper validations
		nestedStructValidations := govalidatormappervalidations.NewStructValidations()

		// Validate the nested struct
		err = d.ValidateRequiredFields(
			nestedStructValidations,
			fieldValue,
			fieldNestedMapper,
		)
		if err != nil {
			return err
		}

		// Add the nested struct validations to the root struct validations
		rootStructValidations.AddNestedStructValidations(
			validationName,
			nestedStructValidations,
		)
	}

	return nil
}
