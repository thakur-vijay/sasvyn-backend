package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sasvyn/backend/internal/auth"
	"github.com/sasvyn/backend/internal/database"
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

	sessionRepository := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepository)

	authHandler := auth.NewHandler(sessionService)

	auth.RegisterRoutes(mux, authHandler)

	authMiddleware := auth.NewMiddleware(sessionService)

	userRepository := users.NewRepository(db)
	userHandler := users.NewHandler(userRepository, sessionService)

	users.RegisterRoutes(mux, userHandler, authMiddleware)

	log.Println("server listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
