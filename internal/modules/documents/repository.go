package documents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, document Document) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateDocument: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO documents (id, user_id, name, category, file_size, key, sync_version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, document.ID, document.UserID, document.Name, document.Category, document.FileSize, document.Key, document.SyncVersion, document.CreatedAt, document.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]Document, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchDocuments: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, category, file_size, key, sync_version, created_at, updated_at FROM documents WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var documents []Document
	for rows.Next() {
		var document Document
		if err := rows.Scan(
			&document.ID,
			&document.UserID,
			&document.Name,
			&document.Category,
			&document.FileSize,
			&document.Key,
			&document.SyncVersion,
			&document.CreatedAt,
			&document.UpdatedAt,
		); err != nil {
			return nil, err
		}
		document.CreatedAt = document.CreatedAt.UTC()
		document.UpdatedAt = document.UpdatedAt.UTC()
		r.attachDocumentURL(&document)
		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return documents, nil
}

func (r *Repository) FetchByID(ctx context.Context, id, userID string) (Document, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchDocumentByID: %v", time.Since(start))
	}()

	var document Document

	err := r.db.QueryRowContext(ctx, `
	SELECT id, user_id, name, category, file_size, key, sync_version, created_at, updated_at FROM documents WHERE id = $1 AND user_id = $2
	`,
		id,
		userID,
	).Scan(
		&document.ID,
		&document.UserID,
		&document.Name,
		&document.Category,
		&document.FileSize,
		&document.Key,
		&document.SyncVersion,
		&document.CreatedAt,
		&document.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Document{}, sql.ErrNoRows
		}

		return Document{}, err
	}

	document.CreatedAt = document.CreatedAt.UTC()
	document.UpdatedAt = document.UpdatedAt.UTC()
	r.attachDocumentURL(&document)
	return document, nil
}

func (r *Repository) Delete(ctx context.Context, documentID, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteDocument: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM documents
		WHERE id = $1
		  AND user_id = $2
	`, documentID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) attachDocumentURL(document *Document) {
	if document.Key == nil {
		return
	}

	url := fmt.Sprintf(
		"%s/%s",
		os.Getenv("R2_PUBLIC_URL"),
		*document.Key,
	)

	document.Url = &url
}
