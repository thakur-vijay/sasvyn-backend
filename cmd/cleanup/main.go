package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sasvyn/backend/internal/config"
	"github.com/sasvyn/backend/internal/database"
	"github.com/sasvyn/backend/internal/modules/idempotency"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repository := idempotency.NewRepository(db)

	if err := repository.DeleteExpiredIdempotencyKeys(context.Background()); err != nil {
		log.Fatal(err)
	}
}
