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
		"POST /skills",
		requireAuth(http.HandlerFunc(handler.Create)),
	)
}
