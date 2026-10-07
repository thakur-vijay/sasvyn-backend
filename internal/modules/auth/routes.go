package auth

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *middleware.Middleware,
	clientIDMiddleware *middleware.ClientIDMiddleware,
) {
	mux.HandleFunc("POST /auth/socialLogin", handler.SocialLogin)
	mux.HandleFunc("POST /auth/refresh", handler.Refresh)

	mux.Handle(
		"POST /auth/logout",
		authMiddleware.RequireAuth(clientIDMiddleware.RequireClientID(http.HandlerFunc(handler.Logout))),
	)
}
