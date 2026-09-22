package languages

import (
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	requireAuth func(http.Handler) http.Handler,
) {
	mux.Handle(
		"GET /languages",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)
	mux.Handle(
		"POST /languages",
		requireAuth(http.HandlerFunc(handler.Create)),
	)

	mux.Handle(
		"PUT /languages/{id}",
		requireAuth(http.HandlerFunc(handler.Update)),
	)

	mux.Handle(
		"DELETE /languages/{id}",
		requireAuth(http.HandlerFunc(handler.Delete)),
	)
}
