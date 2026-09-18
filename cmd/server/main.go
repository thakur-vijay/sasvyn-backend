package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sasvyn/backend/internal/auth"
	"github.com/sasvyn/backend/internal/database"
	"github.com/sasvyn/backend/internal/ratelimit"
	"github.com/sasvyn/backend/internal/sessions"
	"github.com/sasvyn/backend/internal/users"
)

func main() {
	db, err := database.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	defaultLimiter, err := ratelimit.New(
		ratelimit.Policies.Default.Limit,
		ratelimit.Policies.Default.Window,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer defaultLimiter.Close()

	authRefreshLimiter, err := ratelimit.New(
		ratelimit.Policies.AuthRefresh.Limit,
		ratelimit.Policies.AuthRefresh.Window,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer authRefreshLimiter.Close()

	socialLoginLimiter, err := ratelimit.New(
		ratelimit.Policies.SocialLogin.Limit,
		ratelimit.Policies.SocialLogin.Window,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer socialLoginLimiter.Close()

	sessionRepository := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepository)

	authHandler := auth.NewHandler(sessionService)

	auth.RegisterRoutes(mux, authHandler)

	authMiddleware := auth.NewMiddleware(sessionService)

	userRepository := users.NewRepository(db)
	userHandler := users.NewHandler(userRepository, sessionService)

	users.RegisterRoutes(mux, userHandler, authMiddleware)
	rateLimitedMux := ratelimit.PolicyMiddleware(
		defaultLimiter,
		authRefreshLimiter,
		socialLoginLimiter,
	)(mux)

	log.Println("server listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", rateLimitedMux); err != nil {
		log.Fatal(err)
	}
}
