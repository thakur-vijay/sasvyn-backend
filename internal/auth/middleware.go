package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/sasvyn/backend/internal/sessions"
)

type contextKey string

const userIDKey contextKey = "user_id"

type Middleware struct {
	sessionService *sessions.Service
}

func NewMiddleware(sessionService *sessions.Service) *Middleware {
	return &Middleware{sessionService: sessionService}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "invalid, authorization header", http.StatusUnauthorized)
			return
		}

		session, err := m.sessionService.ValidateAccessToken(r.Context(), parts[1])
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
