package validator

import (
	"fmt"
	"log/slog"
	"reflect"

	govalidatormapper "github.com/ralvarezdev/go-validator/mapper"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/mapper/validation"
)

type (
	// DefaultValidator struct
	DefaultValidator struct {
		logger *slog.Logger
	}
)

// NewDefaultValidator creates a new default mapper validator
//
// Parameters:
//
//   - logger: the logger to use
//
// Returns:
//
//   - *DefaultValidator: the default mapper validator
func NewDefaultValidator(
	logger *slog.Logger,
) *DefaultValidator {
	// Create a sub logger
	if logger != nil {
		logger = logger.With(
			slog.String("component", "struct_mapper_validator"),
		)
	}

	return &DefaultValidator{
		logger,
	}
}

// IsFieldInitialized checks if a field is initialized
//
// Parameters:
//
//   - fieldValue: the field value to check
//
// Returns:
//
//   - isInitialized: true if the field is initialized, false otherwise
func (d DefaultValidator) IsFieldInitialized(
	fieldValue reflect.Value,
) (isInitialized bool) {
	// Check if the field is not a pointer and is initialized
	if fieldValue.Kind() != reflect.Ptr {
		return !fieldValue.IsZero()
	}

	// Check if the field is initialized
	if fieldValue.IsNil() {
		return false
	}

	return true
}

// ValidateRequiredFields validates the required fields of a struct
//
// Parameters:
//
//   - rootStructValidations: the root struct validations to validate
//   - mapper: the struct mapper to use
//
// Returns:
//
//   - err: error if any
func (d DefaultValidator) ValidateRequiredFields(
	rootStructValidations *govalidatormappervalidation.StructValidations,
	mapper *govalidatormapper.Mapper,
) error {
	// Check if either the root struct validations, data or the struct fields to validate are nil
	if rootStructValidations == nil {
		return govalidatormappervalidation.ErrNilStructValidations
	}
	if mapper == nil {
		return ErrNilMapper
	}

	// Check if the struct has fields validations
	if !mapper.HasFieldsValidations() {
		return nil
	}

	// Iterate over the fields
	reflection := rootStructValidations.GetReflection()
	reflectedType := reflection.GetReflectedType()
	reflectedValue := reflection.GetReflectedValue()
	structTypeName := reflection.GetReflectedTypeName()
	for i := 0; i < reflectedValue.NumField(); i++ {
		// Get the field value and type
		fieldValue := reflectedValue.Field(i)
		structField := reflectedType.Field(i)
		fieldType := structField.Type
		fieldName := structField.Name
		
		// Check if the field is exported
	 	if structField.PkgPath != "" {
			if d.logger != nil {
				d.logger.Debug(
					"Skipping private (unexported) field on struct type",
					slog.String("struct_type", structTypeName),
					slog.String("field_name", fieldName),
					slog.Any("field_type", fieldType),
					slog.Any("field_value", fieldValue),
				)
			}
			continue
		}
		
		// Check if the field is required
		isRequired, ok := mapper.IsFieldRequired(fieldName)
		if !ok {
			return fmt.Errorf(ErrFieldIsRequiredNotFound, fieldName)
		}
		
		// Check if the field is initialized
		isInitialized := d.IsFieldInitialized(fieldValue)

		// Print field
		if d.logger != nil {
			if isInitialized {
				d.logger.Debug(
					"Detected initialized field on struct type",
					slog.String("struct_type", structTypeName),
					slog.String("field_name", fieldName),
					slog.Any("field_type", fieldType),
					slog.Any("field_value", fieldValue),
					slog.Bool("field_is_required", isRequired),
				)
			} else {
				d.logger.Debug(
					"Detected uninitialized field on struct type",
					slog.String("struct_type", structTypeName),
					slog.String("field_name", fieldName),
					slog.Any("field_type", fieldType),
					slog.Any("field_value", fieldValue),
					slog.Bool("field_is_required", isRequired),
				)
			}
		}

		// Check if the field has to be validated
		if !isRequired {
			continue
		}
		
		// Get the field tag name
		fieldTagName, ok := mapper.GetFieldTagName(fieldName)
		if !ok {
			// Print field
			if d.logger != nil {
				d.logger.Debug(
					"Field tag name not found on struct type",
					slog.String("struct_type", structTypeName),
					slog.String("field_name", fieldName),
					slog.String("field_tag_name", fieldTagName),
					slog.Any("field_value", fieldValue),
				)
			}
			return fmt.Errorf(ErrFieldTagNameNotFound, fieldName)
		}

		// Check if the is initialized
		if !isInitialized {
			rootStructValidations.AddFieldValidationError(
				fieldTagName,
				fmt.Errorf(ErrRequiredField, fieldTagName),
			)
			continue
		}

		// Check if the field is a pointer
		if fieldValue.Kind() != reflect.Ptr {
			if fieldType.Kind() != reflect.Struct {
				continue
			}
		// Check if the field is a scalar required or optional field
		} else if fieldValue.Elem().Kind() != reflect.Struct {
			continue
		}

		// Get the nested struct mapper
		fieldNestedMapper := mapper.GetFieldNestedMapper(fieldName)
		if fieldNestedMapper == nil {
			continue
		}

		// Initialize the nested struct mapper validations
		nestedStructValidations, err := govalidatormappervalidation.NewNestedStructValidations(
			fieldName,
			fieldValue.Interface(),
		)
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
			fieldTagName,
			nestedStructValidations,
		)
	}

	return nil
}
