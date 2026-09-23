package users

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByAppleID(ctx context.Context, appleID string) (*User, error) {
	start := time.Now()

	defer func() {
		log.Printf("[DB] GetByAppleID: %v", time.Since(start))
	}()

	var user User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, apple_id, full_name, email, date_of_birth, img_key, sync_version, created_at, updated_at
		FROM users
		WHERE apple_id = $1
	`, appleID).Scan(
		&user.ID,
		&user.AppleID,
		&user.FullName,
		&user.Email,
		&user.DateOfBirth,
		&user.ImageKey,
		&user.SyncVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()

	r.attachImageURL(&user)

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user User) error {
	start := time.Now()

	defer func() {
		log.Printf("[DB] CreateUser: %v", time.Since(start))
	}()

	_, err := r.db.ExecContext(ctx, `
        INSERT INTO users (
            id,
            apple_id,
            full_name,
            email,
            created_at,
            updated_at,
            sync_version
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		user.ID,
		user.AppleID,
		user.FullName,
		user.Email,
		user.CreatedAt.UTC(),
		user.UpdatedAt.UTC(),
		int64(1),
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, apple_id, full_name, email, date_of_birth, img_key, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.AppleID,
		&user.FullName,
		&user.Email,
		&user.DateOfBirth,
		&user.ImageKey,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()

	r.attachImageURL(&user)

	return &user, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	input UpdateUserDTO,
) (*User, error) {
	start := time.Now()

	defer func() {
		log.Printf("[DB] UpdateUser %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if input.FullName != nil {
		set = append(set, fmt.Sprintf("full_name = $%d", arg))
		args = append(args, *input.FullName)
		arg++
	}

	if input.DateOfBirth != nil {
		set = append(set, fmt.Sprintf("date_of_birth = $%d", arg))
		args = append(args, *input.DateOfBirth)
		arg++
	}

	if input.ImageKey != nil {
		set = append(set, fmt.Sprintf("img_key = $%d", arg))
		args = append(args, *input.ImageKey)
		arg++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	set = append(set, "sync_version = sync_version + 1")

	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE users
         SET %s
         WHERE id = $%d
         RETURNING id, apple_id, full_name, email, date_of_birth,
                   img_key, created_at, updated_at, sync_version`,
		strings.Join(set, ", "),
		arg,
	)

	var user User

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.AppleID,
		&user.FullName,
		&user.Email,
		&user.DateOfBirth,
		&user.ImageKey,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.SyncVersion,
	)
	if err != nil {
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()

	r.attachImageURL(&user)

	return &user, nil
}

func (r *Repository) attachImageURL(user *User) {
	if user.ImageKey == nil {
		return
	}

	url := fmt.Sprintf(
		"%s/%s",
		os.Getenv("R2_PUBLIC_URL"),
		*user.ImageKey,
	)

	user.ImageUrl = &url
}
