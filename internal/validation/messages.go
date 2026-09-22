package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func buildMessage(dto any, err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("validation failed")
	}

	fieldError := validationErrors[0]

	t := reflect.TypeOf(dto)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, ok := t.FieldByName(fieldError.StructField())
	if !ok {
		return fmt.Errorf("validation failed")
	}

	jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
	fieldType := field.Type.String()

	if jsonName == "" {
		jsonName = fieldError.Field()
	}

	switch fieldError.Tag() {
	case "required":
		return fmt.Errorf(
			"%s of type %s is required",
			jsonName,
			fieldType,
		)

	case "min":
		return fmt.Errorf(
			"%s of type %s must be at least %s",
			jsonName,
			fieldType,
			fieldError.Param(),
		)

	case "max":
		return fmt.Errorf(
			"%s of type %s must be at most %s",
			jsonName,
			fieldType,
			fieldError.Param(),
		)

	default:
		return fmt.Errorf("%s is invalid", jsonName)
	}
}
