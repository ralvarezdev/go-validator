package parser

import (
	"log/slog"

	gostringsconvert "github.com/ralvarezdev/go-strings/convert"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/mapper/validation"
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

	// DefaultRawParser is a struct that holds the default raw parser
	DefaultRawParser struct {
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

// AddField adds a field parsed validations to the struct parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field
//   - fieldParsedValidations: The field parsed validations to add
func (s *StructParsedValidations) AddField(
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

// GetField returns the field parsed validations from the struct parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field
//
// Returns:
//
//   - *FieldParsedValidations: The field parsed validations
func (s *StructParsedValidations) GetField(fieldName string) *FieldParsedValidations {
	if s == nil {
		return nil
	}

	// Check if the fields are nil
	if s.fields == nil {
		return nil
	}
	return s.fields[fieldName]
}

// GetFields returns the fields parsed validations from the struct parsed validations
//
// Returns:
//
//   - map[string]*FieldParsedValidations: The fields parsed validations
func (s *StructParsedValidations) GetFields() map[string]*FieldParsedValidations {
	if s == nil {
		return nil
	}
	return s.fields
}

// GetNestedStructs returns the nested structs parsed validations from the struct parsed validations
//
// Returns:
//
//   - map[string]*StructParsedValidations: The nested structs parsed validations
func (s *StructParsedValidations) GetNestedStructs() map[string]*StructParsedValidations {
	if s == nil {
		return nil
	}
	return s.nestedStructs
}

// AddNestedStruct adds a nested struct parsed validations to the struct parsed validations
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - nestedStructParsedValidations: The nested struct parsed validations to add
func (s *StructParsedValidations) AddNestedStruct(
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

// GetNestedStruct returns the nested struct parsed validations from the struct parsed validations
//
// Parameters:
//
//   - nestedStruct: The name of the nested struct
//
// Returns:
//
//   - *StructParsedValidations: The nested struct parsed validations
func (s *StructParsedValidations) GetNestedStruct(nestedStruct string) *StructParsedValidations {
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

// NewDefaultRawParser creates a new DefaultRawParser struct
//
// Parameters:
//
//   - logger: The logger to use
//
// Returns:
//
//   - *DefaultRawParser: The new DefaultRawParser struct
func NewDefaultRawParser(
	logger *slog.Logger,
) *DefaultRawParser {
	if logger != nil {
		// Create a sub logger
		logger = logger.With(
			slog.String("component", "struct_mapper_parser"),
		)
	}

	return &DefaultRawParser{
		logger: logger,
	}
}

// ParseValidations returns the parsed validations from the struct validations
//
// Parameters:
//
//   - structValidations: The root struct validations
//   - dest: The root struct parsed validations to populate
//
// Returns:
//
//   - error: An error if the root struct validations or the root struct parsed validations are nil
func (d DefaultRawParser) ParseValidations(
	structValidations *govalidatormappervalidation.StructValidations,
	dest *StructParsedValidations,
) error {
	// Check if the root struct validations or the root struct parsed validations are nil
	if structValidations == nil {
		return govalidatormappervalidation.ErrNilStructValidations
	}
	if dest == nil {
		return ErrNilStructParsedValidations
	}

	// Check if there are failed validations
	if !structValidations.HasFailed() {
		return nil
	}

	// Get the fields validations, the nested structs validations and the struct type name
	fieldsValidations := structValidations.GetFieldsValidations()
	nestedStructsValidations := structValidations.GetNestedStructsValidations()
	structTypeName := structValidations.GetStructTypeName()

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
			dest.AddField(
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
			err := d.ParseValidations(
				nestedStructValidations,
				nestedStructParsedValidations,
			)
			if err != nil {
				return err
			}

			// Add the nested struct parsed validations to the root struct parsed validations
			dest.AddNestedStruct(
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
