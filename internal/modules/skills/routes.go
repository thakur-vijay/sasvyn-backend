package skills

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
		"GET /skills",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"GET /skills/{id}",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"POST /skills",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create))),
	)

	mux.Handle(
		"PUT /skills/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)

	mux.Handle(
		"DELETE /skills/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete))),
	)
}
