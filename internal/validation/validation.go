package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validator(dto any) error {
	err := validate.Struct(dto)
	if err == nil {
		return nil
	}

	return buildMessage(dto, err)
}
