package mapper

import (
	"log/slog"
	"reflect"
	"strings"

	goreflect "github.com/ralvarezdev/go-reflect"
	gostringsjson "github.com/ralvarezdev/go-strings/json"
)

type (
	// JSONGenerator is a generator for JSON mappers
	JSONGenerator struct {
		logger *slog.Logger
	}
)

// NewJSONGenerator creates a new JSON generator
//
// Parameters:
//
//   - logger: optional logger to use for logging detected fields
//
// Returns:
//
//   - *JSONGenerator: instance of the JSON generator
func NewJSONGenerator(logger *slog.Logger) *JSONGenerator {
	if logger != nil {
		// Create a sub logger
		logger = logger.With(
			slog.String("component", "struct_mapper_json_generator"),
		)
	}

	return &JSONGenerator{
		logger,
	}
}

// NewMapper creates the fields to validate from a JSON struct
//
// Parameters:
//
//   - structInstance: instance of the JSON struct
//
// Returns:
//
//   - *Mapper: instance of the mapper
//   - error: error if any
func (j JSONGenerator) NewMapper(structInstance any) (
	*Mapper,
	error,
) {
	// Check if the struct instance is nil
	if structInstance == nil {
		return nil, ErrNilStructInstance
	}
	
	// Reflection of data
	reflectedType := goreflect.GetDereferencedType(structInstance)

	// Get the struct type name
	structTypeName := reflectedType.Name()

	// Initialize the root map of fields and the map of nested mappers
	rootMapper, err := NewMapper(structInstance)
	if err != nil {
		return nil, err
	}

	// Reflection of the type of data
	for i := 0; i < reflectedType.NumField(); i++ {
		// Get the field type through reflection
		structField := reflectedType.Field(i)
		fieldType := structField.Type
		fieldName := structField.Name

		// Check if the field is unexported
		if !goreflect.IsStructFieldExported(structField) {
			// Set field as not required
			rootMapper.SetFieldIsRequired(fieldName, false)
			continue
		}
			
		// Get the JSON tag of the field
		jsonTag, err := gostringsjson.GetJSONTag(structField, fieldName)
		if err != nil {
			return nil, err
		}

		// Get the JSON name from the tag
		jsonName, err := gostringsjson.GetJSONTagName(jsonTag, fieldName)
		if err != nil {
			return nil, err
		}

		// Add field tag name to the map and set the field as parsed
		rootMapper.AddFieldTagName(fieldName, jsonName)

		// Check if the JSON tag is unassigned or if it contains 'omitempty', which means it is an optional field
		if jsonTag == "-" || strings.Contains(jsonTag, gostringsjson.JSONOmitempty) {
			// Set field name as not required
			rootMapper.SetFieldIsRequired(fieldName, false)

			// Print field
			DetectedField(
				structTypeName,
				fieldName,
				fieldType,
				jsonTag,
				false,
				j.logger,
			)
			continue
		}

		// Set field name as required
		rootMapper.SetFieldIsRequired(fieldName, true)

		// Dereference the pointer
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Check if the element type is a struct
		if fieldType.Kind() == reflect.Struct {
			// Create a new Mapper for the nested struct field
			fieldNestedMapper, err := j.NewMapper(
				reflect.New(fieldType).Interface(),
			)
			if err != nil {
				return nil, err
			}

			// Add the nested fields to the map
			rootMapper.AddFieldNestedMapper(fieldName, fieldNestedMapper)
		}

		// Print field
		DetectedField(
			structTypeName,
			fieldName,
			fieldType,
			jsonTag,
			true,
			j.logger,
		)
	}

	return rootMapper, nil
}

// NewMapperWithNoError creates the fields to validate from a JSON struct
//
// Parameters:
//
//   - structInstance: instance of the JSON struct
//
// Returns:
//
//   - *Mapper: instance of the mapper
func (j JSONGenerator) NewMapperWithNoError(structInstance any) *Mapper {
	mapper, err := j.NewMapper(structInstance)
	if err != nil {
		panic(err)
	}
	return mapper
}
