package database

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(ctx context.Context) (*sql.DB, error) {
	return sql.Open("pgx", "postgres://sasvyn:sasvyn@localhost:5433/sasvyn?sslmode=disable")
}
