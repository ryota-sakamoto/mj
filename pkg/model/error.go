package model

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrLimitExceeded = errors.New("limit exceeded")

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
