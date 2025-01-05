package validator

import (
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatormappervalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	"reflect"
)

type (
	// Validator interface
	Validator interface {
		ValidateNilFields(
			validations govalidatormappervalidations.Validations,
			data interface{},
			mapper *govalidatormapper.Mapper,
		) (err error)
	}

	// DefaultValidator struct
	DefaultValidator struct {
		newValidationsFn func() govalidatormappervalidations.Validations
		logger           *Logger
	}
)

// NewDefaultValidator creates a new default mapper validator
func NewDefaultValidator(
	newValidationsFn func() govalidatormappervalidations.Validations,
	logger *Logger,
) *DefaultValidator {
	return &DefaultValidator{
		newValidationsFn: newValidationsFn,
		logger:           logger,
	}
}

// ValidateNilFields validates if the fields are not nil
func (d *DefaultValidator) ValidateNilFields(
	validations govalidatormappervalidations.Validations,
	data interface{},
	mapper *govalidatormapper.Mapper,
) (err error) {
	// Check if either the validations, data or the struct fields to validate are nil
	if validations == nil {
		return govalidatormapper.ErrNilValidations
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

				validations.AddFieldValidationError(
					validationName,
					govalidatormappervalidations.ErrFieldNotFound,
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

			validations.AddFieldValidationError(
				validationName,
				govalidatormappervalidations.ErrFieldNotFound,
			)
			continue
		}

		// Get the nested struct mapper
		fieldNestedMapper, ok := nestedMappers[fieldType.Name]
		if !ok {
			continue
		}

		// Initialize nested struct mapper validations
		fieldNestedMapperValidations := d.newValidationsFn()

		// Validate nested struct
		err = d.ValidateNilFields(
			fieldNestedMapperValidations,
			fieldValue.Addr().Interface(), // TEST IF THIS A POINTER OF THE STRUCT
			fieldNestedMapper,
		)
		if err != nil {
			return err
		}

		// Add nested struct validations to the struct fields validations
		validations.SetNestedFieldsValidations(
			validationName,
			fieldNestedMapperValidations,
		)
	}

	return nil
}
