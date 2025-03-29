package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidator creates and configures a validator instance
func NewValidator() *validator.Validate {
	v := validator.New()

	// Register a function to get the field name from the json tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

// ValidateStruct validates a struct
func ValidateStruct(v *validator.Validate, s interface{}) error {
	return v.Struct(s)
}

// FormatValidationErrors formats validation errors into a string
func FormatValidationErrors(err error) string {
	if err == nil {
		return ""
	}

	// If it's not a validator.ValidationErrors, return the error message
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return err.Error()
	}

	var errMsgs []string
	for _, e := range validationErrors {
		errMsgs = append(errMsgs, formatValidationError(e))
	}

	return strings.Join(errMsgs, "; ")
}

// formatValidationError formats a single validation error
func formatValidationError(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "email":
		return fieldError.Field() + " must be a valid email address"
	case "min":
		return fieldError.Field() + " must be at least " + fieldError.Param() + " characters long"
	case "max":
		return fieldError.Field() + " must be at most " + fieldError.Param() + " characters long"
	case "oneof":
		return fieldError.Field() + " must be one of " + fieldError.Param()
	default:
		return fieldError.Field() + " failed on " + fieldError.Tag() + " validation"
	}
}
