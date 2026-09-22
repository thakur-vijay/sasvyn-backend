package upload

import (
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	requireAuth func(http.Handler) http.Handler,
) {
	mux.Handle(
		"POST /upload-url",
		requireAuth(http.HandlerFunc(handler.CreateUploadURL)),
	)
}
