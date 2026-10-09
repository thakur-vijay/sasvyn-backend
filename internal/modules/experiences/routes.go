package experiences

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
	// MARK: - Experiences

	mux.Handle(
		"GET /experiences",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Fetch)),
	)

	mux.Handle(
		"GET /experiences/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchByID)),
	)

	mux.Handle(
		"POST /experiences",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Create)))),
	)

	mux.Handle(
		"PUT /experiences/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Update)))),
	)

	mux.Handle(
		"DELETE /experiences/{id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.Delete)))),
	)

	// MARK: - Experience Responsibilities

	mux.Handle(
		"GET /experiences/{experience_id}/responsibilities",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchResponsibilities)),
	)

	mux.Handle(
		"GET /experiences/{experience_id}/responsibilities/{responsibility_id}",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.FetchResponsibilityByID)),
	)

	mux.Handle(
		"POST /experiences/{experience_id}/responsibilities",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.CreateResponsibility)))),
	)

	mux.Handle(
		"PUT /experiences/{experience_id}/responsibilities/{responsibility_id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.UpdateResponsibility)))),
	)

	mux.Handle(
		"DELETE /experiences/{experience_id}/responsibilities/{responsibility_id}",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(idempotency.RequireIdempotencyKey(http.HandlerFunc(handler.DeleteResponsibility)))),
	)
}
