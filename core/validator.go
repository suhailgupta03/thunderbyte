package core

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type RequestValidator struct {
	validator *validator.Validate
}

func (rv *RequestValidator) Validate(i interface{}) error {
	if reflect.Struct != reflect.ValueOf(i).Kind() {
		return errors.New("argument passed must be a struct")
	}

	if err := rv.validator.Struct(i); err != nil {
		return err
	}

	return nil
}
