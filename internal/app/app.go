package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/config"
	"github.com/sasvyn/backend/internal/database"
	"github.com/sasvyn/backend/internal/storage"
)

type App struct {
	DB       *sql.DB
	Router   http.Handler
	Limiters *RateLimiters
	R2Client *s3.Client
}

func New(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, nil
	}

	db, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(db, cfg.DatabaseURL); err != nil {
		db.Close()
		return nil, err
	}

	limiters, err := NewRateLimiters()
	if err != nil {
		db.Close()
		return nil, err
	}
	r2Client := storage.NewR2Client()
	router := BuildRouter(db, limiters, r2Client)

	return &App{
		DB:       db,
		Router:   router,
		Limiters: limiters,
		R2Client: r2Client,
	}, nil
}

func (a *App) Close() error {
	if a == nil {
		return nil
	}

	if a.Limiters != nil {
		a.Limiters.Close()
	}

	if a.DB != nil {
		return a.DB.Close()
	}

	return nil
}
