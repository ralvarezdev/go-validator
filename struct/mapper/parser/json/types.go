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

// NewFlattenedParsedValidations adds the root struct parsed validations to the flattened parsed validations
//
// Parameters:
//
//   - structParsedValidations: The root struct parsed validations to add
//
// Returns:
//
//   - error: An error if the root struct parsed validations are nil or if the fields are already in the flattened parsed validations
func NewFlattenedParsedValidations(
	structParsedValidations *govalidatormapperparser.StructParsedValidations,
) (*FlattenedParsedValidations, error) {
	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return nil, govalidatormapperparser.ErrNilStructParsedValidations
	}

	// Create the flattened parsed validations
	f := &FlattenedParsedValidations{
		fields: make(map[string]interface{}),
	}

	// Add the struct parsed validations fields
	fieldsParsedValidations := structParsedValidations.GetFields()
	if fieldsParsedValidations != nil {
		for fieldName, fieldParsedValidations := range fieldsParsedValidations {
			// Add the field parsed validations
			if err := f.AddField(
				fieldName,
				fieldParsedValidations,
			); err != nil {
				return nil, err
			}
		}
	}

	// Add the struct parsed validations nested structs
	nestedStructsParsedValidations := structParsedValidations.GetNestedStructs()
	if nestedStructsParsedValidations != nil {
		for fieldName, nestedStructParsedValidations := range nestedStructsParsedValidations {
			// Add the nested struct parsed validations
			if err := f.AddNestedStruct(
				fieldName,
				nestedStructParsedValidations,
			); err != nil {
				return nil, err
			}
		}
	}

	return f, nil
}

// AddField adds a field parsed validations to the flattened parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field
//   - fieldParsedValidations: The field parsed validations to add
//
// Returns:
//
//   - error: An error if the field name is already in the flattened parsed validations
func (f *FlattenedParsedValidations) AddField(
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
		return fmt.Errorf(ErrFieldNameAlreadyParsed, fieldName)
	}

	// Add the field parsed validations to the flattened parsed validations
	f.fields[fieldName] = fieldParsedValidations.GetErrors()
	return nil
}

// AddNestedStruct adds a nested struct parsed validations to the flattened parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - structParsedValidations: The struct parsed validations to add
//
// Returns:
//
//   - error: An error if the struct name is already in the flattened parsed validations
func (f *FlattenedParsedValidations) AddNestedStruct(
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
		return fmt.Errorf(ErrFieldNameAlreadyParsed, fieldName)
	}

	// Get the struct flattened parsed validations
	structFlattenedParsedValidations, err := NewFlattenedParsedValidations(structParsedValidations)
	if err != nil {
		return err
	}

	// Add the struct parsed validations to the flattened parsed validations
	f.fields[fieldName] = structFlattenedParsedValidations.GetFields()
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
	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return nil, govalidatormapperparser.ErrNilStructParsedValidations
	}

	// Flatten the parsed validations
	flattenedParsedValidations, err := NewFlattenedParsedValidations(
		structParsedValidations,
	)
	if err != nil {
		return nil, err
	}
	return flattenedParsedValidations.GetFields(), nil
}
