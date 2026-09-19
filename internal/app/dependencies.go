package app

import (
	"database/sql"
	"net/http"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/modules/users"
	"github.com/sasvyn/backend/internal/ratelimit"
)

func BuildRouter(db *sql.DB, limiters *RateLimiters) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	sessionRepository := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepository)

	authRepository := auth.NewRepository(db)
	authHandler := auth.NewHandler(authRepository, sessionService)
	auth.RegisterRoutes(mux, authHandler)

	userRepository := users.NewRepository(db)
	userHandler := users.NewHandler(userRepository, sessionService)
	users.RegisterRoutes(mux, userHandler)

	return ratelimit.PolicyMiddleware(
		limiters.Default,
		limiters.AuthRefresh,
		limiters.SocialLogin,
	)(mux)
}
