package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sasvyn/backend/internal/apperrors"
)

func Migrate(db *sql.DB, databaseURL string) error {
	if databaseURL == "" {
		return apperrors.NewConfigError("DATABASE_URL is not set", nil)
	}

	m, err := migrate.New(
		"file://internal/database/migrations",
		databaseURL,
	)
	if err != nil {
		return apperrors.NewServiceUnavailable("database migration setup failed", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return apperrors.NewServiceUnavailable("database migration failed", err)
	}

	return nil
}
