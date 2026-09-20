package users

import (
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	requireAuth func(http.Handler) http.Handler,
) {
	mux.Handle(
		"GET /users/{id}",
		requireAuth(http.HandlerFunc(handler.GetByID)),
	)

	mux.Handle(
		"PUT /users/{id}",
		requireAuth(http.HandlerFunc(handler.Update)),
	)
}
