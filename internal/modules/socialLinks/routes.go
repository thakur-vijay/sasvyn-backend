package sociallinks

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
		"GET /socialLinks",
		requireAuth(http.HandlerFunc(handler.Fetch)),
	)
	mux.Handle(
		"POST /socialLinks",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create))),
	)

	mux.Handle(
		"PUT /socialLinks/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)

	mux.Handle(
		"GET /socialLinks/{id}",
		requireAuth(http.HandlerFunc(handler.FetchByID)),
	)

	mux.Handle(
		"DELETE /socialLinks/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete))),
	)
}
