package languages

import (
	"net/http"

	"github.com/sasvyn/backend/internal/modules/idempotency"
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
		"GET /languages/{id}",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)
	mux.Handle(
		"POST /languages",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create))),
	)

	mux.Handle(
		"PUT /languages/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)

	mux.Handle(
		"DELETE /languages/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete))),
	)
}
