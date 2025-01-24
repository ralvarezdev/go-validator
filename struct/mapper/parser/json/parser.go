package json

import (
	"fmt"
	gostringsconvert "github.com/ralvarezdev/go-strings/convert"
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
)

type (
	// StructParsedValidations is the struct for the struct JSON parsed validations
	StructParsedValidations struct {
		fieldName      *string
		structTypeName string
		nestedStructs  *map[string]*StructParsedValidations
		fields         *map[string]*FieldParsedValidations
	}

	// FieldParsedValidations is the struct for the field JSON parsed validations
	FieldParsedValidations struct {
		errors *[]string
	}

	// FlattenedParsedValidations is the struct for the flattened JSON parsed validations
	FlattenedParsedValidations struct {
		fields *map[string]interface{}
	}

	// Parser is a struct that holds the JSON parser
	Parser struct {
		logger *Logger
	}
)

// NewStructParsedValidations creates a new StructParsedValidations struct
func NewStructParsedValidations(structTypeName string) *StructParsedValidations {
	return &StructParsedValidations{structTypeName: structTypeName}
}

// NewNestedStructParsedValidations creates a new nested StructParsedValidations struct
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
func (s *StructParsedValidations) GetStructTypeName() string {
	return s.structTypeName
}

// AddFieldParsedValidations adds a field parsed validations to the struct parsed validations
func (s *StructParsedValidations) AddFieldParsedValidations(
	fieldName string,
	fieldParsedValidations *FieldParsedValidations,
) {
	// Check if the field name is empty or the field parsed validations are nil
	if fieldName == "" || fieldParsedValidations == nil {
		return
	}

	// Check if the fields are nil
	if s.fields == nil {
		fields := make(map[string]*FieldParsedValidations)
		s.fields = &fields
	}

	// Add the field parsed validations to the struct parsed validations
	(*s.fields)[fieldName] = fieldParsedValidations
}

// GetFieldParsedValidations returns the field parsed validations from the struct parsed validations
func (s *StructParsedValidations) GetFieldParsedValidations(fieldName string) *FieldParsedValidations {
	// Check if the fields are nil
	if s.fields == nil {
		return nil
	}
	return (*s.fields)[fieldName]
}

// GetFieldsParsedValidations returns the fields parsed validations from the struct parsed validations
func (s *StructParsedValidations) GetFieldsParsedValidations() *map[string]*FieldParsedValidations {
	return s.fields
}

// AddNestedStructParsedValidations adds a nested struct parsed validations to the struct parsed validations
func (s *StructParsedValidations) AddNestedStructParsedValidations(
	fieldName string,
	nestedStructParsedValidations *StructParsedValidations,
) {
	// Check if the nested struct name is empty or the nested struct parsed validations are nil
	if fieldName == "" || nestedStructParsedValidations == nil {
		return
	}

	// Check if the nested structs are nil
	if s.nestedStructs == nil {
		nestedStructs := make(map[string]*StructParsedValidations)
		s.nestedStructs = &nestedStructs
	}

	// Add the nested struct parsed validations to the struct parsed validations
	(*s.nestedStructs)[fieldName] = nestedStructParsedValidations
}

// GetNestedStructParsedValidations returns the nested struct parsed validations from the struct parsed validations
func (s *StructParsedValidations) GetNestedStructParsedValidations(nestedStruct string) *StructParsedValidations {
	// Check if the nested structs are nil
	if s.nestedStructs == nil {
		return nil
	}
	return (*s.nestedStructs)[nestedStruct]
}

// NewFieldParsedValidations creates a new FieldParsedValidations struct
func NewFieldParsedValidations() *FieldParsedValidations {
	return &FieldParsedValidations{}
}

// AddErrors adds errors to the field parsed validations
func (f *FieldParsedValidations) AddErrors(errors *[]error) {
	// Check if the errors are nil
	if errors == nil {
		return
	}

	// Check if the field parsed validations errors are nil
	if f.errors == nil {
		f.errors = &[]string{}
	}

	// Append the errors to the field parsed validations
	mappedErrors := gostringsconvert.ErrorArrayToStringArray(errors)
	*f.errors = append(*f.errors, *mappedErrors...)
}

// AddError adds an error to the field parsed validations
func (f *FieldParsedValidations) AddError(error string) {
	// Check if the error is empty
	if error == "" {
		return
	}

	// Check if the errors are nil
	if f.errors == nil {
		var errors []string
		f.errors = &errors
	}

	// Append the error to the field parsed validations
	*f.errors = append(*f.errors, error)
}

// GetErrors returns the errors from the field parsed validations
func (f *FieldParsedValidations) GetErrors() *[]string {
	return f.errors
}

// NewFlattenedParsedValidations creates a new FlattenedParsedValidations struct
func NewFlattenedParsedValidations() *FlattenedParsedValidations {
	return &FlattenedParsedValidations{}
}

// AddFieldParsedValidations adds a field parsed validations to the flattened JSON parsed validations
func (f *FlattenedParsedValidations) AddFieldParsedValidations(
	fieldName string,
	fieldParsedValidations *FieldParsedValidations,
) error {
	// Check if the field name is empty or the field parsed validations are nil
	if fieldName == "" || fieldParsedValidations == nil {
		return nil
	}

	// Check if the fields are nil
	if f.fields == nil {
		fields := make(map[string]interface{})
		f.fields = &fields
	}

	// Check if the field name is already in the flattened JSON parsed validations
	if _, ok := (*f.fields)[fieldName]; ok {
		return fmt.Errorf(ErrNilFieldNameAlreadyParsed, fieldName)
	}

	// Add the field parsed validations to the flattened JSON parsed validations
	(*f.fields)[fieldName] = fieldParsedValidations.GetErrors()
	return nil
}

// AddNestedStructParsedValidations adds a nested struct parsed validations to the flattened JSON parsed validations
func (f *FlattenedParsedValidations) AddNestedStructParsedValidations(
	fieldName string,
	structParsedValidations *StructParsedValidations,
) error {
	// Check if the struct name is empty or the struct parsed validations are nil
	if structParsedValidations == nil {
		return nil
	}

	// Check if the fields are nil
	if f.fields == nil {
		fields := make(map[string]interface{})
		f.fields = &fields
	}

	// Check if the struct name is already in the flattened JSON parsed validations
	if _, ok := (*f.fields)[fieldName]; ok {
		return fmt.Errorf(ErrNilFieldNameAlreadyParsed, fieldName)
	}

	// Get the struct flattened JSON parsed validations
	structFlattenedParsedValidations := NewFlattenedParsedValidations()
	err := structFlattenedParsedValidations.AddRootStructParsedValidations(structParsedValidations)
	if err != nil {
		return err
	}

	// Add the struct parsed validations to the flattened JSON parsed validations
	(*f.fields)[fieldName] = structFlattenedParsedValidations.GetFields()
	return nil
}

// AddRootStructParsedValidations adds the root struct parsed validations to the flattened JSON parsed validations
func (f *FlattenedParsedValidations) AddRootStructParsedValidations(
	structParsedValidations *StructParsedValidations,
) error {
	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return ErrNilParsedValidations
	}

	// Check if the fields are nil
	if f.fields != nil {
		return ErrFlattenedParsedValidationsAlreadyExists
	}

	// Initialize the fields
	fields := make(map[string]interface{})
	f.fields = &fields

	// Add the struct parsed validations fields
	fieldsParsedValidations := structParsedValidations.GetFieldsParsedValidations()
	if fieldsParsedValidations != nil {
		for fieldName, fieldParsedValidations := range *fieldsParsedValidations {
			// Check if the field name is already in the flattened JSON parsed validations
			if _, ok := (*f.fields)[fieldName]; ok {
				return fmt.Errorf(ErrNilFieldNameAlreadyParsed, fieldName)
			}

			// Add the field parsed validations
			(*f.fields)[fieldName] = fieldParsedValidations.GetErrors()
		}
	}

	// Add the struct parsed validations nested structs
	nestedStructsParsedValidations := structParsedValidations.nestedStructs
	if nestedStructsParsedValidations != nil {
		for fieldName, nestedStructParsedValidations := range *nestedStructsParsedValidations {
			// Check if the nested struct name is already in the flattened JSON parsed validations
			if _, ok := (*f.fields)[fieldName]; ok {
				return fmt.Errorf(
					ErrNilFieldNameAlreadyParsed,
					fieldName,
				)
			}

			// Get the nested struct flattened JSON parsed validations
			nestedStructFlattenedParsedValidations := NewFlattenedParsedValidations()
			err := nestedStructFlattenedParsedValidations.AddRootStructParsedValidations(
				nestedStructParsedValidations,
			)
			if err != nil {
				return err
			}

			// Add the nested struct parsed validations
			(*f.fields)[fieldName] = nestedStructFlattenedParsedValidations.GetFields()
		}
	}

	return nil
}

// GetFields returns the fields from the flattened JSON parsed validations
func (f *FlattenedParsedValidations) GetFields() *map[string]interface{} {
	return f.fields
}

// NewParser creates a new Parser struct
func NewParser(logger *Logger) *Parser {
	return &Parser{logger: logger}
}

// GenerateParsedValidations returns a
func (p *Parser) GenerateParsedValidations(
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
	var fieldValidationsErrors *[]error
	if fieldsValidations != nil {
		for fieldName, fieldValidations := range *fieldsValidations {
			// Check if the field validations are nil
			if fieldValidations == nil {
				continue
			}

			// Check if the field has no errors
			fieldValidationsErrors = fieldValidations.GetErrors()
			if fieldValidationsErrors == nil || len(*fieldValidationsErrors) == 0 {
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
			if p.logger != nil {
				p.logger.FieldParsedValidations(
					structTypeName,
					fieldName,
					fieldParsedValidations,
				)
			}
		}
	}

	// Iterate over all nested structs validations
	if nestedStructsValidations != nil {
		for fieldName, nestedStructValidations := range *nestedStructsValidations {
			// Check if the nested struct validations are nil
			if nestedStructValidations == nil {
				continue
			}

			// Generate the nested JSON parsed validations
			nestedStructTypeName := nestedStructValidations.GetStructTypeName()
			nestedStructParsedValidations := NewNestedStructParsedValidations(
				fieldName,
				nestedStructTypeName,
			)
			err := p.GenerateParsedValidations(
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
			if p.logger != nil {
				p.logger.NestedStructParsedValidations(
					structTypeName,
					fieldName,
					nestedStructTypeName,
				)
			}
		}
	}

	return nil
}

// ParseValidations parses the validations into JSON
func (p *Parser) ParseValidations(rootStructValidations *govalidatormappervalidation.StructValidations) (
	interface{},
	error,
) {
	// Initialize the parsed validations
	rootParsedValidations := NewStructParsedValidations(rootStructValidations.GetStructTypeName())

	// Generate the JSON parsed validations
	err := p.GenerateParsedValidations(
		rootStructValidations,
		rootParsedValidations,
	)
	if err != nil {
		return nil, err
	}

	// Flatten the JSON parsed validations
	flattenedParsedValidations := NewFlattenedParsedValidations()
	err = flattenedParsedValidations.AddRootStructParsedValidations(
		rootParsedValidations,
	)
	if err != nil {
		return nil, err
	}

	return flattenedParsedValidations.GetFields(), nil
}
