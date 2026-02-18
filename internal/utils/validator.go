package utils

import (
	"strings"
	"unicode/utf8"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator handles validation logic
type Validator struct {
	errors []ValidationError
}

func NewValidator() *Validator {
	return &Validator{errors: []ValidationError{}}
}

func (v *Validator) Required(field, value string, message string) {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{Field: field, Message: message})
	}
}

func (v *Validator) MinLength(field, value string, min int, message string) {
	if utf8.RuneCountInString(value) < min {
		v.errors = append(v.errors, ValidationError{Field: field, Message: message})
	}
}

func (v *Validator) MaxLength(field, value string, max int, message string) {
	if utf8.RuneCountInString(value) > max {
		v.errors = append(v.errors, ValidationError{Field: field, Message: message})
	}
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}
