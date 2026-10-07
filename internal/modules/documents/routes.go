package documents

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/idempotency"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *middleware.Middleware,
	clientIDMiddleware *middleware.ClientIDMiddleware,
) {
	mux.Handle(
		"GET /documents",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"GET /documents/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchByID)),
	)

	mux.Handle(
		"POST /documents",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create)))),
	)

	mux.Handle(
		"DELETE /documents/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete)))),
	)
}
