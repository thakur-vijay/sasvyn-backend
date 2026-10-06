package realtime

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	requireAuth func(http.Handler) http.Handler,
) {
	mux.Handle(
		"GET /ws",
		requireAuth(http.HandlerFunc(handler.WebSocket)),
	)
}
