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
	structTypeName string,
	fieldName string,
	fieldType reflect.Type,
	tag string,
	required bool,
) {
	l.logger.Debug(
		"detected field on struct type: "+structTypeName,
		"field name: "+fieldName,
		fmt.Sprintf("field type: '%v'", fieldType),
		"field tag: "+tag,
		fmt.Sprintf("field is required: '%v'", required),
	)
}
