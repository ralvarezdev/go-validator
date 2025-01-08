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

// IsFieldInitialized checks if a field is initialized
func (d *DefaultValidator) IsFieldInitialized(
	fieldValue reflect.Value,
) (isInitialized bool) {
	// Check if the field is not a pointer and is initialized
	if fieldValue.Kind() != reflect.Ptr {
		if fieldValue.IsZero() {
			return false
		}
		return true
	}

	// Check if the field is initialized
	if fieldValue.IsNil() {
		return false
	}

	return true
}

// ValidateRequiredFields validates the required fields of a struct
func (d *DefaultValidator) ValidateRequiredFields(
	rootStructValidations *govalidatormappervalidations.StructValidations,
	mapper *govalidatormapper.Mapper,
) error {
	// Check if either the root struct validations, data or the struct fields to validate are nil
	if rootStructValidations == nil {
		return govalidatormappervalidations.ErrNilStructValidations
	}
	if mapper == nil {
		return ErrNilMapper
	}

	// Check if the struct has fields validations
	if !mapper.HasFieldsValidations() {
		return nil
	}

	// Iterate over the fields
	typeReflection := rootStructValidations.GetTypeReflection()
	valueReflection := rootStructValidations.GetValueReflection()
	structName := rootStructValidations.GetStructName()
	for i := 0; i < valueReflection.NumField(); i++ {
		// Get the field value and type
		fieldValue := valueReflection.Field(i)
		structField := typeReflection.Field(i)
		fieldType := structField.Type
		fieldName := structField.Name

		// Get the field tag name
		validationName, isRequired := mapper.GetFieldValidationName(fieldName)

		// Check if the field is parsed
		isParsed, ok := mapper.IsFieldParsed(fieldName)
		if !ok || !isParsed {
			continue
		}

		// Check if the field is initialized
		isInitialized := d.IsFieldInitialized(fieldValue)

		// Print field
		if d.logger != nil {
			if isInitialized {
				d.logger.InitializedField(
					structName,
					fieldName,
					fieldType,
					fieldValue,
					isRequired,
				)
			} else {
				d.logger.UninitializedField(
					structName,
					fieldName,
					fieldType,
					isRequired,
				)
			}
		}

		// Check if the field has to be validated
		if !isRequired {
			continue
		}

		// Check if the field is a pointer
		if !isInitialized {
			rootStructValidations.AddFieldValidationError(
				validationName,
				fmt.Errorf(ErrRequiredField, validationName),
			)
			continue
		}

		// Check if the field is a scalar required or optional field
		if fieldValue.Kind() != reflect.Ptr || fieldValue.Elem().Kind() != reflect.Struct {
			continue
		}

		// Get the nested struct mapper
		fieldNestedMapper := mapper.GetFieldNestedMapper(fieldName)
		if fieldNestedMapper == nil {
			continue
		}

		// Initialize the nested struct mapper validations
		nestedStructValidations, err := govalidatormappervalidations.NewStructValidations(fieldValue)
		if err != nil {
			return err
		}

		// Validate the nested struct
		err = d.ValidateRequiredFields(
			nestedStructValidations,
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
