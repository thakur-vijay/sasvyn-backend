package upload

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *middleware.Middleware,
) {
	mux.Handle(
		"POST /upload-url",
		authMiddleware.RequireAuth(http.HandlerFunc(handler.CreateUploadURL)),
	)
}
