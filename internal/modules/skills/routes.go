package skills

import (
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	requireAuth func(http.Handler) http.Handler,
) {
	mux.Handle(
		"GET /skills",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)
	mux.Handle(
		"POST /skills",
		requireAuth(http.HandlerFunc(handler.Create)),
	)

	mux.Handle(
		"PUT /skills/{id}",
		requireAuth(http.HandlerFunc(handler.Update)),
	)

	mux.Handle(
		"DELETE /skills/{id}",
		requireAuth(http.HandlerFunc(handler.Delete)),
	)
}
