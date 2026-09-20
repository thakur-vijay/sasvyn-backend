package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}
