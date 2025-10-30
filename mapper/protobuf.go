package mapper

import (
	"log/slog"
	"reflect"

	goreflect "github.com/ralvarezdev/go-reflect"
)

type (
	// ProtobufGenerator is a generator for Protobuf mappers
	ProtobufGenerator struct {
		logger *slog.Logger
	}
)

// NewProtobufGenerator creates a new Protobuf generator
//
// Parameters:
//
//   - logger: optional logger to use for logging detected fields
//
// Returns:
//
//   - *ProtobufGenerator: instance of the Protobuf generator
func NewProtobufGenerator(logger *slog.Logger) *ProtobufGenerator {
	if logger != nil {
		// Create a sub logger
		logger = logger.With(
			slog.String("component", "struct_mapper_protobuf_generator"),
		)
	}

	return &ProtobufGenerator{
		logger,
	}
}

// NewMapper creates the fields to validate from a Protobuf compiled struct
//
// Parameters:
//
//   - structInstance: instance of the Protobuf compiled struct
//
// Returns:
//
//   - *Mapper: instance of the mapper
//   - error: error if any
func (p ProtobufGenerator) NewMapper(structInstance any) (
	*Mapper,
	error,
) {	
	// Reflection of data
	reflectedType := goreflect.GetDereferencedType(structInstance)

	// Get the struct type name
	structTypeName := goreflect.GetTypeName(reflectedType)

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

		// Omit protobuf internal fields
		if IsProtobufGeneratedField(fieldName) {
			continue
		}

		// Get the Protobuf tag of the field
		protobufTag, err := GetProtobufTag(structField, fieldName)
		if err != nil {
			return nil, err
		}

		// Get the field name from the Protobuf tag
		protobufName, err := GetProtobufTagName(protobufTag, fieldName)
		if err != nil {
			return nil, err
		}

		// Add the field to the fields map
		rootMapper.AddFieldTagName(fieldName, protobufName)

		// Check if the field is a pointer
		if fieldType.Kind() == reflect.Ptr {
			// Dereference the pointer
			fieldType = fieldType.Elem()

			// Check if the element type is not a struct and the tag to determine if it contains 'oneof', which means it
			// is an optional struct field
			if fieldType.Kind() != reflect.Struct || IsProtobufFieldOptional(protobufTag) {
				// Set field as not required
				rootMapper.SetFieldIsRequired(fieldName, false)

				// Print field
				DetectedField(
					structTypeName,
					fieldName,
					fieldType,
					protobufTag,
					false,
					p.logger,
				)
				continue
			}

			// Create a new Mapper for the nested struct field
			fieldNestedMapper, err := p.NewMapper(
				reflect.New(fieldType).Interface(),
			)
			if err != nil {
				return nil, err
			}

			// Add the nested fields to the map
			rootMapper.AddFieldNestedMapper(fieldName, fieldNestedMapper)
		}

		// Check if the field is an interface (special case for oneof fields)
		if fieldType.Kind() == reflect.Interface {
			// Set field as not required
			rootMapper.SetFieldIsRequired(fieldName, false)

			// Print field
			DetectedField(
				structTypeName,
				fieldName,
				fieldType,
				protobufTag,
				false,
				p.logger,
			)
			continue
		}

		// Set field as required
		rootMapper.SetFieldIsRequired(fieldName, true)

		// Print field
		DetectedField(
			structTypeName,
			fieldName,
			fieldType,
			protobufTag,
			true,
			p.logger,
		)
	}

	return rootMapper, nil
}

// NewMapperWithNoError creates the fields to validate from a Protobuf compiled struct
//
// Parameters:
//
//   - structInstance: instance of the Protobuf compiled struct
//
// Returns:
//
//   - *Mapper: instance of the mapper
func (p ProtobufGenerator) NewMapperWithNoError(structInstance any) *Mapper {
	mapper, err := p.NewMapper(structInstance)
	if err != nil {
		panic(err)
	}
	return mapper
}