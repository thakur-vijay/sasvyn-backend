package skills

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
		"GET /skills",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"GET /skills/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchByID)),
	)

	mux.Handle(
		"POST /skills",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create)))),
	)

	mux.Handle(
		"PUT /skills/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update)))),
	)

	mux.Handle(
		"DELETE /skills/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete)))),
	)
}
