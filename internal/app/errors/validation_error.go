package errors

import (
	"fmt"
)

type ValidationError struct {
	*Error
	field  *string
	reason *string
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Error: BadRequestError().WithVerbose("validation error"),
	}
}

func (v *ValidationError) WithReason(reason string) *ValidationError {
	v.reason = &reason
	return v
}

func (v *ValidationError) WithFieldName(fieldName string) *ValidationError {
	v.field = &fieldName
	return v
}

func (v *ValidationError) BuildError() *Error {
	if v.reason != nil {
		v.WithVerbose(fmt.Sprintf("for reason: [%s]", *v.reason))
	}
	if v.field != nil {
		v.WithVerbose(fmt.Sprintf("in field: [%s]", *v.field))
	}
	v.message = errorMessage(v.verbose)
	return v.Error
}
