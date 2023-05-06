package model

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type ValidationError struct {
	reason string
}

func NewValidationError(reason string) ValidationError {
	return ValidationError{
		reason: reason,
	}
}

func (v ValidationError) Error() string {
	return v.reason
}
