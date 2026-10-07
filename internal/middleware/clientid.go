package middleware

import (
	"net/http"

	"github.com/sasvyn/backend/internal/response"
)

const clientIDHeader = "X-Client-ID"

type ClientIDMiddleware struct{}

func NewClientIDMiddleware() *ClientIDMiddleware {
	return &ClientIDMiddleware{}
}

func (m *ClientIDMiddleware) RequireClientID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get(clientIDHeader)

		if clientID == "" {
			response.Write(w, http.StatusBadRequest, "client id is required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
