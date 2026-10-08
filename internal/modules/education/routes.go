package education

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
		"GET /educations",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"GET /educations/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchByID)),
	)

	mux.Handle(
		"POST /educations",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create)))),
	)

	mux.Handle(
		"PUT /educations/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update)))),
	)

	mux.Handle(
		"DELETE /educations/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete)))),
	)
}
