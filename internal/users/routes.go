package users

import (
	"net/http"

	"github.com/sasvyn/backend/internal/auth"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *auth.Middleware,
) {
	mux.HandleFunc("POST /socialLogin", handler.SocialLogin)

	mux.Handle(
		"GET /me",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Me)),
	)
}
