package realtime

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *middleware.Middleware,
) {
	mux.Handle(
		"GET /ws",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.WebSocket)),
	)
}
