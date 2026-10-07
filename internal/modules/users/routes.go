package users

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
		"GET /users/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.GetByID)),
	)

	mux.Handle(
		"PUT /users/{id}",
		authMiddleware.RequireAuth(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update))),
	)
}
