package auth

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *middleware.Middleware,
) {
	mux.HandleFunc("POST /auth/socialLogin", handler.SocialLogin)
	mux.HandleFunc("POST /auth/refresh", handler.Refresh)

	mux.Handle(
		"POST /auth/logout",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Logout)),
	)
}
