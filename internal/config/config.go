package config

import (
	"os"
	"strings"

	"github.com/sasvyn/backend/internal/apperrors"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() (*Config, error) {
	port := normalizePort(getEnv("APP_PORT", "8080"))

	cfg := &Config{
		Port:        port,
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if cfg.DatabaseURL == "" {
		return nil, apperrors.NewConfigError("DATABASE_URL is not set", nil)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func normalizePort(port string) string {
	if port == "" {
		return ":8080"
	}

	if strings.Contains(port, ":") {
		return port
	}

	return ":" + port
}
