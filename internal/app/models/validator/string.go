package validator

import (
	"messenger/internal/app/errors"
)

type StringValidatorBuilder struct {
	value string
	err *errors.ValidationError
}

func String(value string) *StringValidatorBuilder {
	return &StringValidatorBuilder{value, nil}
}

func (s *StringValidatorBuilder) Required() *StringValidatorBuilder {
	if len(s.value) == 0 {
		s.err = errors.NewValidationError().WithReason("is required")
	}
	return s
}

func (s *StringValidatorBuilder) MaxLen(maxLen int) *StringValidatorBuilder {
	if len(s.value) > maxLen {
		s.err = errors.NewValidationError().WithReason("too long")
	}
	return s
}

func (s *StringValidatorBuilder) MinLen(maxLen int) *StringValidatorBuilder {
	if len(s.value) < maxLen {
		s.err = errors.NewValidationError().WithReason("too short")
	}
	return s.Required()
}

func (s *StringValidatorBuilder) Validate() *errors.ValidationError {
	return s.err
}