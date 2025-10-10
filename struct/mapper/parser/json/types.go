package json

import (
	"fmt"

	govalidatormapperparser "github.com/ralvarezdev/go-validator/struct/mapper/parser"
)

type (
	// FlattenedParsedValidations is the struct for the flattened parsed validations
	FlattenedParsedValidations struct {
		fields map[string]interface{}
	}

	// DefaultEndParser is the default implementation of the EndParser interface
	DefaultEndParser struct{}
)

// NewFlattenedParsedValidations creates a new FlattenedParsedValidations struct
//
// Returns:
//
//   - *FlattenedParsedValidations: The new FlattenedParsedValidations struct
func NewFlattenedParsedValidations() *FlattenedParsedValidations {
	return &FlattenedParsedValidations{}
}

// AddFieldParsedValidations adds a field parsed validations to the flattened parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field
//   - fieldParsedValidations: The field parsed validations to add
//
// Returns:
//
//   - error: An error if the field name is already in the flattened parsed validations
func (f *FlattenedParsedValidations) AddFieldParsedValidations(
	fieldName string,
	fieldParsedValidations *govalidatormapperparser.FieldParsedValidations,
) error {
	if f == nil {
		return ErrNilFlattenedParsedValidations
	}

	// Check if the field name is empty or the field parsed validations are nil
	if fieldName == "" || fieldParsedValidations == nil {
		return nil
	}

	// Check if the fields are nil
	if f.fields == nil {
		f.fields = make(map[string]interface{})
	}

	// Check if the field name is already in the flattened parsed validations
	if _, ok := f.fields[fieldName]; ok {
		return fmt.Errorf(ErrNilFieldNameAlreadyParsed, fieldName)
	}

	// Add the field parsed validations to the flattened parsed validations
	f.fields[fieldName] = fieldParsedValidations.GetErrors()
	return nil
}

// AddNestedStructParsedValidations adds a nested struct parsed validations to the flattened parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - structParsedValidations: The struct parsed validations to add
//
// Returns:
//
//   - error: An error if the struct name is already in the flattened parsed validations
func (f *FlattenedParsedValidations) AddNestedStructParsedValidations(
	fieldName string,
	structParsedValidations *govalidatormapperparser.StructParsedValidations,
) error {
	if f == nil {
		return ErrNilFlattenedParsedValidations
	}

	// Check if the struct name is empty or the struct parsed validations are nil
	if structParsedValidations == nil {
		return nil
	}

	// Check if the fields are nil
	if f.fields == nil {
		f.fields = make(map[string]interface{})
	}

	// Check if the struct name is already in the flattened parsed validations
	if _, ok := f.fields[fieldName]; ok {
		return fmt.Errorf(ErrNilFieldNameAlreadyParsed, fieldName)
	}

	// Get the struct flattened parsed validations
	structFlattenedParsedValidations := NewFlattenedParsedValidations()
	err := structFlattenedParsedValidations.AddRootStructParsedValidations(structParsedValidations)
	if err != nil {
		return err
	}

	// Add the struct parsed validations to the flattened parsed validations
	f.fields[fieldName] = structFlattenedParsedValidations.GetFields()
	return nil
}

// AddRootStructParsedValidations adds the root struct parsed validations to the flattened parsed validations
//
// Parameters:
//
//   - structParsedValidations: The root struct parsed validations to add
//
// Returns:
//
//   - error: An error if the root struct parsed validations are nil or if the fields are already in the flattened parsed validations
func (f *FlattenedParsedValidations) AddRootStructParsedValidations(
	structParsedValidations *govalidatormapperparser.StructParsedValidations,
) error {
	if f == nil {
		return ErrNilFlattenedParsedValidations
	}

	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return govalidatormapperparser.ErrNilParsedValidations
	}

	// Check if the fields are nil
	if f.fields != nil {
		return ErrFlattenedParsedValidationsAlreadyExists
	}

	// Initialize the fields
	f.fields = make(map[string]interface{})

	// Add the struct parsed validations fields
	fieldsParsedValidations := structParsedValidations.GetFields()
	if fieldsParsedValidations != nil {
		for fieldName, fieldParsedValidations := range fieldsParsedValidations {
			// Check if the field name is already in the flattened parsed validations
			if _, ok := f.fields[fieldName]; ok {
				return fmt.Errorf(ErrNilFieldNameAlreadyParsed, fieldName)
			}

			// Add the field parsed validations
			f.fields[fieldName] = fieldParsedValidations.GetErrors()
		}
	}

	// Add the struct parsed validations nested structs
	nestedStructsParsedValidations := structParsedValidations.GetNestedStructs()
	if nestedStructsParsedValidations != nil {
		for fieldName, nestedStructParsedValidations := range nestedStructsParsedValidations {
			// Check if the nested struct name is already in the flattened parsed validations
			if _, ok := f.fields[fieldName]; ok {
				return fmt.Errorf(
					ErrNilFieldNameAlreadyParsed,
					fieldName,
				)
			}

			// Get the nested struct flattened parsed validations
			nestedStructFlattenedParsedValidations := NewFlattenedParsedValidations()
			err := nestedStructFlattenedParsedValidations.AddRootStructParsedValidations(
				nestedStructParsedValidations,
			)
			if err != nil {
				return err
			}

			// Add the nested struct parsed validations
			f.fields[fieldName] = nestedStructFlattenedParsedValidations.GetFields()
		}
	}

	return nil
}

// GetFields returns the fields from the flattened parsed validations
//
// Returns:
//
//   - map[string]interface{}: The fields
func (f *FlattenedParsedValidations) GetFields() map[string]interface{} {
	if f == nil {
		return nil
	}
	return f.fields
}

// NewDefaultEndParser creates a new DefaultEndParser
//
// Returns:
//
//   - DefaultEndParser: The new DefaultEndParser
func NewDefaultEndParser() DefaultEndParser {
	return DefaultEndParser{}
}

// ParseValidations parses the validations into a flattened map[string]interface{}
//
// Parameters:
//
//   - structValidations: The root struct validations
//
// Returns:
//
//   - interface{}: The parsed validations
//   - error: An error if the root struct validations are nil or if there was an error generating or flattening the parsed validations
func (d DefaultEndParser) ParseValidations(structParsedValidations *govalidatormapperparser.StructParsedValidations) (
	interface{},
	error,
) {
	// Flatten the parsed validations
	flattenedParsedValidations := NewFlattenedParsedValidations()
	err := flattenedParsedValidations.AddRootStructParsedValidations(
		structParsedValidations,
	)
	if err != nil {
		return nil, err
	}
	return flattenedParsedValidations.GetFields(), nil
}
