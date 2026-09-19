package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/sasvyn/backend/internal/database"
)

type App struct {
	DB       *sql.DB
	Router   http.Handler
	Limiters *RateLimiters
}

func New() (*App, error) {
	db, err := database.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	limiters, err := NewRateLimiters()
	if err != nil {
		db.Close()
		return nil, err
	}

	router := BuildRouter(db, limiters)

	return &App{
		DB:       db,
		Router:   router,
		Limiters: limiters,
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
