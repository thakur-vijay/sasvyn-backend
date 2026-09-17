package main

import (
	"log"
	"net/http"

	"github.com/sasvyn/backend/internal/users"
)

func main() {
	mux := http.NewServeMux()

	users.RegisterRoutes(mux)

	log.Println("server listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}