package sociallinks

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/idempotency"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *middleware.Middleware,
) {
	mux.Handle(
		"GET /socialLinks",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Fetch)),
	)
	mux.Handle(
		"POST /socialLinks",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create))),
	)

	mux.Handle(
		"PUT /socialLinks/{id}",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)

	mux.Handle(
		"GET /socialLinks/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchByID)),
	)

	mux.Handle(
		"DELETE /socialLinks/{id}",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete))),
	)
}
