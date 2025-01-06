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

// PrintFieldParsedValidations logs the parsed validations
func (l *Logger) PrintFieldParsedValidations(
	fieldName string,
	fieldValidations *FieldParsedValidations,
) {
	l.logger.Debug(
		fmt.Sprintf("parsed validations to field '%v'", fieldName),
		fmt.Sprintf("parsed validations: %v", fieldValidations),
	)
}

// PrintStructParsedValidations logs the parsed validations
func (l *Logger) PrintStructParsedValidations(
	structName string,
	nestedStructParsedValidations *StructParsedValidations,
) {
	l.logger.Debug(
		fmt.Sprintf("parsed validations to struct '%v'", structName),
		fmt.Sprintf("parsed validations: %v", nestedStructParsedValidations),
	)
}
