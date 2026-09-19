package auth

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *Middleware,
) {
	mux.HandleFunc("POST /socialLogin", handler.SocialLogin)
	mux.HandleFunc("POST /auth/refresh", handler.Refresh)

	mux.Handle(
		"POST /auth/logout",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.Logout)),
	)
}
