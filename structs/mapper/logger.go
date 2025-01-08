package mapper

import (
	"fmt"
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
	"reflect"
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

// DetectedField prints a detected field
func (l *Logger) DetectedField(
	structName string,
	fieldName string,
	fieldType reflect.Type,
	tag string,
	required bool,
	parsed bool,
) {
	l.logger.Debug(
		fmt.Sprintf("detected field on struct '%v'", structName),
		fmt.Sprintf("field '%v'", fieldName),
		fmt.Sprintf("type: '%v'", fieldType),
		fmt.Sprintf("tag: '%v'", tag),
		fmt.Sprintf("required: '%v'", required),
		fmt.Sprintf("parsed: '%v'", parsed),
	)
}
