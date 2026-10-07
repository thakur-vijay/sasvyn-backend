package languages

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
		"GET /languages",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"GET /languages/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchByID)),
	)
	mux.Handle(
		"POST /languages",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create))),
	)

	mux.Handle(
		"PUT /languages/{id}",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)

	mux.Handle(
		"DELETE /languages/{id}",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete))),
	)
}
