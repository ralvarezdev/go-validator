package json

import (
	"fmt"
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
)

// Logger is the JWT validator logger
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger creates a new JWT validator logger
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// FieldParsedValidations logs the parsed validations
func (l *Logger) FieldParsedValidations(
	structTypeName,
	fieldName string,
	fieldValidations *FieldParsedValidations,
) {
	// Get the errors
	errors := fieldValidations.GetErrors()
	if errors == nil {
		return
	}

	// Log the parsed validations
	l.logger.Debug(
		"parsed validations to struct type: "+structTypeName,
		"field name: "+fieldName,
		fmt.Sprintf("field validations errors: %v", *errors),
	)
}

// NestedStructParsedValidations logs the parsed validations
func (l *Logger) NestedStructParsedValidations(
	structTypeName,
	fieldName,
	nestedStructTypeName string,
) {
	// Log the parsed validations
	l.logger.Debug(
		"parsed validations to struct type: "+structTypeName,
		"field name: "+fieldName,
		"nested struct type: "+nestedStructTypeName,
	)
}
