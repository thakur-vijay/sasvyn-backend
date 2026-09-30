package users

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
		"GET /users/{id}",
		requireAuth(http.HandlerFunc(handler.GetByID)),
	)

	mux.Handle(
		"PUT /users/{id}",
		requireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)
}
