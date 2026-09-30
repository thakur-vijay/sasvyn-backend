package sociallinks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, link SocialLink) (SocialLink, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateSocialLink: %v", time.Since(start))
	}()

	var createdLink SocialLink

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO social_links (
			id,
			user_id,
			type,
			url,
			sync_version,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			user_id,
			type,
			url,
			sync_version,
			created_at,
			updated_at`,
		link.ID,
		link.UserID,
		link.Type,
		link.Url,
		int64(1),
		link.CreatedAt,
		link.UpdatedAt,
	).Scan(
		&createdLink.ID,
		&createdLink.UserID,
		&createdLink.Type,
		&createdLink.Url,
		&createdLink.SyncVersion,
		&createdLink.CreatedAt,
		&createdLink.UpdatedAt,
	)

	if err != nil {
		return SocialLink{}, err
	}

	createdLink.CreatedAt = createdLink.CreatedAt.UTC()
	createdLink.UpdatedAt = createdLink.UpdatedAt.UTC()

	return createdLink, nil
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]SocialLink, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchSocialLinks: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, type, url, sync_version, created_at, updated_at FROM social_links WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var socialLinks []SocialLink
	for rows.Next() {
		var link SocialLink
		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Type,
			&link.Url,
			&link.SyncVersion,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, err
		}
		link.CreatedAt = link.CreatedAt.UTC()
		link.UpdatedAt = link.UpdatedAt.UTC()
		socialLinks = append(socialLinks, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return socialLinks, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	userID string,
	request UpdateSocialLinkDTO,
) (SocialLink, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateSocialLink: %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if request.Type != nil {
		set = append(set, fmt.Sprintf("type = $%d", arg))
		args = append(args, *request.Type)
		arg++
	}

	if request.Url != nil {
		set = append(set, fmt.Sprintf("url = $%d", arg))
		args = append(args, *request.Url)
		arg++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	set = append(set, "sync_version = sync_version + 1")

	args = append(args, id, userID)

	query := fmt.Sprintf(
		`UPDATE social_links
		 SET %s
		 WHERE id = $%d AND user_id = $%d
		 RETURNING
			id,
			user_id,
			type,
			url,
			sync_version,
			created_at,
			updated_at`,
		strings.Join(set, ", "),
		arg,
		arg+1,
	)

	var updatedLink SocialLink

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&updatedLink.ID,
		&updatedLink.UserID,
		&updatedLink.Type,
		&updatedLink.Url,
		&updatedLink.SyncVersion,
		&updatedLink.CreatedAt,
		&updatedLink.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SocialLink{}, sql.ErrNoRows
		}

		return SocialLink{}, err
	}

	updatedLink.CreatedAt = updatedLink.CreatedAt.UTC()
	updatedLink.UpdatedAt = updatedLink.UpdatedAt.UTC()

	return updatedLink, nil
}

func (r *Repository) Delete(ctx context.Context, linkID, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteSocialLink: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM social_links
		WHERE id = $1
		  AND user_id = $2
	`, linkID, userID)
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
