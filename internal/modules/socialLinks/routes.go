package sociallinks

import (
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	requireAuth func(http.Handler) http.Handler,
) {
	mux.Handle(
		"GET /socialLinks",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)
	mux.Handle(
		"POST /socialLinks",
		requireAuth(http.HandlerFunc(handler.Create)),
	)

	mux.Handle(
		"PUT /socialLinks/{id}",
		requireAuth(http.HandlerFunc(handler.Update)),
	)

	mux.Handle(
		"DELETE /socialLinks/{id}",
		requireAuth(http.HandlerFunc(handler.Delete)),
	)
}
