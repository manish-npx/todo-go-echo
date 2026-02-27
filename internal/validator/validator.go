package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps validator
type CustomValidator struct {
	validator *validator.Validate
}

// New creates validator instance
func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate implements echo.Validator.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
