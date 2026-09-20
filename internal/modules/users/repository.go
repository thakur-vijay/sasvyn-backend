package users

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByAppleID(ctx context.Context, appleID string) (*User, error) {
	var user User

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
	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()
	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, apple_id, full_name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.AppleID, user.FullName, user.Email, user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, apple_id, full_name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
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
	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()
	return &user, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	input UpdateUserDTO,
) (*User, error) {
	var user User

	err := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET full_name = $1,
		    date_of_birth = $2,
		    updated_at = NOW()
		WHERE id = $3
		RETURNING id, apple_id, full_name, email, date_of_birth, created_at, updated_at
	`,
		input.FullName,
		input.DateOfBirth,
		id,
	).Scan(
		&user.ID,
		&user.AppleID,
		&user.FullName,
		&user.Email,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()

	return &user, nil
}
