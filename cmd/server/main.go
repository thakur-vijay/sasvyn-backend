package main

// @title Sasvyn Backend API
// @version 1.0
// @description API documentation for Sasvyn backend.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your token as: Bearer <token>
import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/sasvyn/backend/docs"
	"github.com/sasvyn/backend/internal/app"
	"github.com/sasvyn/backend/internal/config"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	if err := application.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}
}
