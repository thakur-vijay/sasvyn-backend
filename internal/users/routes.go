package users

import (
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {
	mux.HandleFunc("GET /{id}", handler.GetByID)
}
