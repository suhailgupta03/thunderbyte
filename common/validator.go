package common

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RequestValidator struct {
	validator *validator.Validate
}

func (cv *RequestValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewRequestValidator creates a new RequestValidator
func NewRequestValidator() echo.Validator {
	return &RequestValidator{validator: validator.New()}
}
