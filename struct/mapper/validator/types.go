package validator

type (
	// ValidateFn is the type for the validate function
	ValidateFn func(toValidate interface{}) (interface{}, error)
)
