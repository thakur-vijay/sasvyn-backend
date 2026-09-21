package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func DecodeJSON[T any](r *http.Request, destination *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read request body")
	}

	if len(strings.TrimSpace(string(body))) == 0 {
		return errors.New("request body must not be empty")
	}

	var raw map[string]json.RawMessage

	if err := json.Unmarshal(body, &raw); err != nil {
		return errors.New("request body must be a valid JSON object")
	}

	if len(raw) == 0 {
		return errors.New("request body must not be empty")
	}

	allowedFields, err := allowedJSONFields[T]()
	if err != nil {
		return err
	}

	for field := range raw {
		if _, exists := allowedFields[field]; !exists {
			return fmt.Errorf(
				"field %q is invalid; it must contain fields in %s",
				field,
				formatAllowedFields(allowedFields),
			)
		}
	}

	if err := json.Unmarshal(body, destination); err != nil {
		return fmt.Errorf(
			"request body is invalid; it must contain fields in %s",
			formatAllowedFields(allowedFields),
		)
	}

	return nil
}

func allowedJSONFields[T any]() (map[string]string, error) {
	var value T

	typ := reflect.TypeOf(value)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("request destination must be a struct")
	}

	fields := make(map[string]string)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.PkgPath != "" {
			continue
		}

		jsonTag := field.Tag.Get("json")

		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		name := strings.Split(jsonTag, ",")[0]

		if name == "" {
			continue
		}

		fields[name] = field.Type.String()
	}

	return fields, nil
}

func formatAllowedFields(fields map[string]string) string {
	var result []string

	for name, fieldType := range fields {
		result = append(result, fmt.Sprintf("%s: %s", name, fieldType))
	}

	return strings.Join(result, ", ")
}
