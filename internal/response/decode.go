package response

import (
	"net/http"

	"github.com/sasvyn/backend/internal/validation"
)

func DecodeJSONAndValidate[T any](
	r *http.Request,
	destination *T,
) error {
	if err := DecodeJSON(r, destination); err != nil {
		return err
	}

	return validation.Validator(*destination)
}
