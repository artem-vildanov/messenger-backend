package validators

import (
	"messenger/internal/app/errors"
)

type StringValidator struct {
	value string
	err *errors.ValidationError
}

func String(value string) *StringValidator {
	return &StringValidator{value, nil}
}

func (s *StringValidator) Required() *StringValidator {
	if len(s.value) == 0 {
		s.err = errors.NewValidationError().WithReason("is required")
	}
	return s
}

func (s *StringValidator) MaxLen(maxLen int) *StringValidator {
	if len(s.value) > maxLen {
		s.err = errors.NewValidationError().WithReason("too long")
	}
	return s
}

func (s *StringValidator) MinLen(maxLen int) *StringValidator {
	if len(s.value) < maxLen {
		s.err = errors.NewValidationError().WithReason("too short")
	}
	return s.Required()
}

func (s *StringValidator) Validate() *errors.ValidationError {
	return s.err
}
