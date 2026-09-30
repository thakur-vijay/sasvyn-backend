package idempotency

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

func (r *Repository) Find(
	ctx context.Context,
	userID string,
	key string,
) (*IdempotencyRecord, error) {

	var record IdempotencyRecord

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			user_id,
			idempotency_key,
			status,
			status_code,
			response_body,
			created_at
		FROM idempotency_keys
		WHERE user_id = $1
		  AND idempotency_key = $2
		`,
		userID,
		key,
	).Scan(
		&record.ID,
		&record.UserID,
		&record.IdempotencyKey,
		&record.Status,
		&record.StatusCode,
		&record.ResponseBody,
		&record.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &record, nil
}

func (r *Repository) Save(
	ctx context.Context,
	record IdempotencyRecord,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO idempotency_keys (
			id,
			user_id,
			idempotency_key,
			status,
			status_code,
			response_body,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		record.ID,
		record.UserID,
		record.IdempotencyKey,
		record.Status,
		record.StatusCode,
		record.ResponseBody,
		record.CreatedAt,
	)

	return err
}

func (r *Repository) Reserve(
	ctx context.Context,
	record IdempotencyRecord,
) (*IdempotencyRecord, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO idempotency_keys (
			id,
			user_id,
			idempotency_key,
			status,
			created_at
		)
		VALUES ($1, $2, $3, 'processing', $4)
		ON CONFLICT (user_id, idempotency_key) DO NOTHING
		`,
		record.ID,
		record.UserID,
		record.IdempotencyKey,
		record.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return r.Find(ctx, record.UserID, record.IdempotencyKey)
}

func (r *Repository) Complete(
	ctx context.Context,
	userID string,
	key string,
	statusCode int,
	responseBody []byte,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE idempotency_keys
		SET
			status = 'completed',
			status_code = $3,
			response_body = $4
		WHERE user_id = $1
		  AND idempotency_key = $2
		`,
		userID,
		key,
		statusCode,
		responseBody,
	)

	return err
}

func (r *Repository) Delete(
	ctx context.Context,
	userID string,
	key string,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE FROM idempotency_keys
		WHERE user_id = $1
		  AND idempotency_key = $2
		`,
		userID,
		key,
	)

	return err
}
