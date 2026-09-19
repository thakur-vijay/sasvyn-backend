package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func New(statusCode int, message string, data any) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func Success(statusCode int, message string, data any) Response {
	return New(statusCode, message, data)
}

func Error(statusCode int, message string) Response {
	return New(statusCode, message, nil)
}

func Write(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(New(statusCode, message, data))
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	Write(w, statusCode, message, nil)
}
