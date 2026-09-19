package auth

import (
	"context"
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByAppleID(ctx context.Context, appleID string) (*UserResponse, error) {
	var user UserResponse
	log.Printf("GetByAppleID: %v", appleID)
	err := r.db.QueryRowContext(ctx, `
		SELECT id, apple_id, full_name, email, created_at, updated_at
		FROM users
		WHERE apple_id = $1
	`, appleID).Scan(
		&user.ID,
		&user.AppleID,
		&user.FullName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user UserResponse) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, apple_id, full_name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`, user.ID, user.AppleID, user.FullName, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}
