package idempotency

import (
	"net/http"

	"github.com/sasvyn/backend/internal/response"
)

func RequireIdempotencyKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")

		if key == "" {
			response.Write(
				w,
				http.StatusBadRequest,
				"Idempotency-Key header is required",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
