package validations

import (
	"fmt"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"github.com/ralvarezdev/go-validator/structs/mapper"
	"reflect"
)

// ValidateMapperNilFields validates if the fields are not nil
func ValidateMapperNilFields(
	data interface{},
	mapper *mapper.Mapper,
	mode *goflagsmode.Flag,
) (mapperValidations *MapperValidations, err error) {
	// Check if either the data or the struct fields to validate are nil
	if data == nil {
		return nil, NilDataError
	}
	if mapper == nil {
		return nil, NilMapperError
	}

	// Initialize struct fields validations
	mapperValidations = NewMapperValidations()

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
		if mode != nil && mode.IsDebug() {
			fmt.Printf("field '%v' of type '%v' value: %v\n", fieldType.Name, fieldType.Type, fieldValue)
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
				if mode != nil && mode.IsDebug() {
					fmt.Printf("field is uninitialized: %v\n", fieldType.Name)
				}
				mapperValidations.AddFailedFieldValidationError(validationName, FieldNotFoundError)
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
			if mode != nil && mode.IsDev() {
				fmt.Printf("field is uninitialized: %v\n", fieldType.Name)
			}
			mapperValidations.AddFailedFieldValidationError(validationName, FieldNotFoundError)
			continue
		}

		// Get the nested struct fields to validate
		fieldNestedMapper, ok := nestedMappers[fieldType.Name]
		if !ok {
			continue
		}

		// Validate nested struct
		fieldNestedMapperValidations, err := ValidateMapperNilFields(
			fieldValue.Addr().Interface(), // TEST IF THIS A POINTER OF THE STRUCT
			fieldNestedMapper,
			mode,
		)
		if err != nil {
			return nil, err
		}

		// Add nested struct validations to the struct fields validations
		mapperValidations.SetNestedMapperValidations(validationName, fieldNestedMapperValidations)
	}

	return mapperValidations, nil
}
