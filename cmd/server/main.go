package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sasvyn/backend/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	if err := application.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
