package parser

import (
	"fmt"
	"log/slog"

	gostringsconvert "github.com/ralvarezdev/go-strings/convert"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
)

type (
	// StructParsedValidations is the struct for the struct parsed validations
	StructParsedValidations struct {
		fieldName      *string
		structTypeName string
		nestedStructs  map[string]*StructParsedValidations
		fields         map[string]*FieldParsedValidations
	}

	// FieldParsedValidations is the struct for the field parsed validations
	FieldParsedValidations struct {
		errors []string
	}

	// FlattenedParsedValidations is the struct for the flattened parsed validations
	FlattenedParsedValidations struct {
		fields map[string]interface{}
	}

	// DefaultParser is a struct that holds the parser
	DefaultParser struct {
		logger *slog.Logger
	}
)

// NewStructParsedValidations creates a new StructParsedValidations struct
//
// Parameters:
//
//   - structTypeName: The name of the struct type
//
// Returns:
//
//   - *StructParsedValidations: The new StructParsedValidations struct
func NewStructParsedValidations(structTypeName string) *StructParsedValidations {
	return &StructParsedValidations{structTypeName: structTypeName}
}

// NewNestedStructParsedValidations creates a new nested StructParsedValidations struct
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - structTypeName: The name of the nested struct type
//
// Returns:
//
//   - *StructParsedValidations: The new nested StructParsedValidations struct
func NewNestedStructParsedValidations(
	fieldName string,
	structTypeName string,
) *StructParsedValidations {
	return &StructParsedValidations{
		fieldName:      &fieldName,
		structTypeName: structTypeName,
	}
}

// GetStructTypeName returns the struct type name from the struct parsed validations
//
// Returns:
//
//   - string: The struct type name
func (s *StructParsedValidations) GetStructTypeName() string {
	if s == nil {
		return ""
	}
	return s.structTypeName
}

// AddFieldParsedValidations adds a field parsed validations to the struct parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field
//   - fieldParsedValidations: The field parsed validations to add
func (s *StructParsedValidations) AddFieldParsedValidations(
	fieldName string,
	fieldParsedValidations *FieldParsedValidations,
) {
	if s == nil {
		return
	}

	// Check if the field name is empty or the field parsed validations are nil
	if fieldName == "" || fieldParsedValidations == nil {
		return
	}

	// Check if the fields are nil
	if s.fields == nil {
		fields := make(map[string]*FieldParsedValidations)
		s.fields = fields
	}

	// Add the field parsed validations to the struct parsed validations
	s.fields[fieldName] = fieldParsedValidations
}

// GetFieldParsedValidations returns the field parsed validations from the struct parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field
//
// Returns:
//
//   - *FieldParsedValidations: The field parsed validations
func (s *StructParsedValidations) GetFieldParsedValidations(fieldName string) *FieldParsedValidations {
	if s == nil {
		return nil
	}

	// Check if the fields are nil
	if s.fields == nil {
		return nil
	}
	return s.fields[fieldName]
}

// GetFieldsParsedValidations returns the fields parsed validations from the struct parsed validations
//
// Returns:
//
//   - map[string]*FieldParsedValidations: The fields parsed validations
func (s *StructParsedValidations) GetFieldsParsedValidations() map[string]*FieldParsedValidations {
	if s == nil {
		return nil
	}
	return s.fields
}

// AddNestedStructParsedValidations adds a nested struct parsed validations to the struct parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - nestedStructParsedValidations: The nested struct parsed validations to add
func (s *StructParsedValidations) AddNestedStructParsedValidations(
	fieldName string,
	nestedStructParsedValidations *StructParsedValidations,
) {
	if s == nil {
		return
	}

	// Check if the nested struct name is empty or the nested struct parsed validations are nil
	if fieldName == "" || nestedStructParsedValidations == nil {
		return
	}

	// Check if the nested structs are nil
	if s.nestedStructs == nil {
		s.nestedStructs = make(map[string]*StructParsedValidations)
	}

	// Add the nested struct parsed validations to the struct parsed validations
	s.nestedStructs[fieldName] = nestedStructParsedValidations
}

// GetNestedStructParsedValidations returns the nested struct parsed validations from the struct parsed validations
//
// Parameters:
//
//   - nestedStruct: The name of the nested struct
//
// Returns:
//
//   - *StructParsedValidations: The nested struct parsed validations
func (s *StructParsedValidations) GetNestedStructParsedValidations(nestedStruct string) *StructParsedValidations {
	if s == nil {
		return nil
	}

	// Check if the nested structs are nil
	if s.nestedStructs == nil {
		return nil
	}
	return s.nestedStructs[nestedStruct]
}

// NewFieldParsedValidations creates a new FieldParsedValidations struct
//
// Returns:
//
//   - *FieldParsedValidations: The new FieldParsedValidations struct
func NewFieldParsedValidations() *FieldParsedValidations {
	return &FieldParsedValidations{}
}

// AddErrors adds errors to the field parsed validations
//
// Parameters:
//
//   - errors: The errors to add
func (f *FieldParsedValidations) AddErrors(errors []error) {
	if f == nil {
		return
	}

	// Check if the errors are nil
	if errors == nil {
		return
	}

	// Check if the field parsed validations errors are nil
	if f.errors == nil {
		f.errors = []string{}
	}

	// Append the errors to the field parsed validations
	mappedErrors := gostringsconvert.ErrorArrayToStringArray(errors)
	if mappedErrors == nil || len(mappedErrors) == 0 {
		return
	}
	f.errors = append(f.errors, mappedErrors...)
}

// AddError adds an error to the field parsed validations
//
// Parameters:
//
//   - error: The error to add
func (f *FieldParsedValidations) AddError(error string) {
	if f == nil {
		return
	}

	// Check if the error is empty
	if error == "" {
		return
	}

	// Check if the errors are nil
	if f.errors == nil {
		f.errors = make([]string, 0)
	}

	// Append the error to the field parsed validations
	f.errors = append(f.errors, error)
}

// GetErrors returns the errors from the field parsed validations
//
// Returns:
//
//   - []string: The errors
func (f *FieldParsedValidations) GetErrors() []string {
	if f == nil {
		return nil
	}
	return f.errors
}

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
	fieldParsedValidations *FieldParsedValidations,
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
	structParsedValidations *StructParsedValidations,
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
	structParsedValidations *StructParsedValidations,
) error {
	if f == nil {
		return ErrNilFlattenedParsedValidations
	}

	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return ErrNilParsedValidations
	}

	// Check if the fields are nil
	if f.fields != nil {
		return ErrFlattenedParsedValidationsAlreadyExists
	}

	// Initialize the fields
	f.fields = make(map[string]interface{})

	// Add the struct parsed validations fields
	fieldsParsedValidations := structParsedValidations.GetFieldsParsedValidations()
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
	nestedStructsParsedValidations := structParsedValidations.nestedStructs
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

// NewParser creates a new DefaultParser struct
//
// Parameters:
//
//   - logger: The logger to use
//
// Returns:
//
//   - *DefaultParser: The new DefaultParser struct
func NewParser(logger *slog.Logger) *DefaultParser {
	if logger != nil {
		// Create a sub logger
		logger = logger.With(
			slog.String("component", "struct_mapper_parser"),
		)
	}

	return &DefaultParser{logger: logger}
}

// GenerateParsedValidations returns the parsed validations from the struct validations
//
// Parameters:
//
//   - rootStructValidations: The root struct validations
//   - rootStructParsedValidations: The root struct parsed validations to populate
//
// Returns:
//
//   - error: An error if the root struct validations or the root struct parsed validations are nil
func (d DefaultParser) GenerateParsedValidations(
	rootStructValidations *govalidatormappervalidation.StructValidations,
	rootStructParsedValidations *StructParsedValidations,
) error {
	// Check if the root struct validations or the root struct parsed validations are nil
	if rootStructValidations == nil {
		return govalidatormappervalidation.ErrNilStructValidations
	}
	if rootStructParsedValidations == nil {
		return ErrNilParsedValidations
	}

	// Check if there are failed validations
	if !rootStructValidations.HasFailed() {
		return nil
	}

	// Get the fields validations, the nested structs validations and the struct type name
	fieldsValidations := rootStructValidations.GetFieldsValidations()
	nestedStructsValidations := rootStructValidations.GetNestedStructsValidations()
	structTypeName := rootStructValidations.GetStructTypeName()

	// Iterate over all fields and their errors
	var fieldValidationsErrors []error
	if fieldsValidations != nil {
		for fieldName, fieldValidations := range fieldsValidations {
			// Check if the field validations are nil
			if fieldValidations == nil {
				continue
			}

			// Check if the field has no errors
			fieldValidationsErrors = fieldValidations.GetErrors()
			if fieldValidationsErrors == nil || len(fieldValidationsErrors) == 0 {
				continue
			}

			// Add the field parsed validations to the root struct parsed validations
			fieldParsedValidations := NewFieldParsedValidations()
			fieldParsedValidations.AddErrors(fieldValidationsErrors)
			rootStructParsedValidations.AddFieldParsedValidations(
				fieldName,
				fieldParsedValidations,
			)

			// Print the field parsed validations
			if d.logger != nil {
				// Get the errors
				errors := fieldValidations.GetErrors()
				if errors == nil {
					return nil
				}

				// Log the parsed validations
				d.logger.Debug(
					"parsed validations to struct type",
					slog.String("struct_type", structTypeName),
					slog.String("field_name", fieldName),
					slog.Any("errors", errors),
				)
			}
		}
	}

	// Iterate over all nested structs validations
	if nestedStructsValidations != nil {
		for fieldName, nestedStructValidations := range nestedStructsValidations {
			// Check if the nested struct validations are nil
			if nestedStructValidations == nil {
				continue
			}

			// Generate the nested parsed validations
			nestedStructTypeName := nestedStructValidations.GetStructTypeName()
			nestedStructParsedValidations := NewNestedStructParsedValidations(
				fieldName,
				nestedStructTypeName,
			)
			err := d.GenerateParsedValidations(
				nestedStructValidations,
				nestedStructParsedValidations,
			)
			if err != nil {
				return err
			}

			// Add the nested struct parsed validations to the root struct parsed validations
			rootStructParsedValidations.AddNestedStructParsedValidations(
				fieldName,
				nestedStructParsedValidations,
			)

			// Print the nested struct parsed validations
			if d.logger != nil {
				// Log the parsed validations
				d.logger.Debug(
					"parsed validations to struct type: "+structTypeName,
					slog.String("field_name", fieldName),
					slog.String(
						"nested_struct_type_name",
						nestedStructTypeName,
					),
				)
			}
		}
	}

	return nil
}

// ParseValidations parses the validations into a flattened map[string]interface{}
//
// Parameters:
//
//   - rootStructValidations: The root struct validations
//
// Returns:
//
//   - interface{}: The parsed validations
//   - error: An error if the root struct validations are nil or if there was an error generating or flattening the parsed validations
func (d DefaultParser) ParseValidations(rootStructValidations *govalidatormappervalidation.StructValidations) (
	interface{},
	error,
) {
	// Initialize the parsed validations
	rootParsedValidations := NewStructParsedValidations(rootStructValidations.GetStructTypeName())

	// Generate the parsed validations
	err := d.GenerateParsedValidations(
		rootStructValidations,
		rootParsedValidations,
	)
	if err != nil {
		return nil, err
	}

	// Flatten the parsed validations
	flattenedParsedValidations := NewFlattenedParsedValidations()
	err = flattenedParsedValidations.AddRootStructParsedValidations(
		rootParsedValidations,
	)
	if err != nil {
		return nil, err
	}

	return flattenedParsedValidations.GetFields(), nil
}
