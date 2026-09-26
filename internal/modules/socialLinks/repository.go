package sociallinks

import (
	"context"
	"database/sql"
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

func (r *Repository) Create(ctx context.Context, link SocialLink) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateSocialLink: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO social_links (id, user_id, type, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`, link.ID, link.UserID, link.Type, link.Url, link.CreatedAt, link.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]SocialLink, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchSocialLinks: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, type, url, created_at, updated_at FROM social_links WHERE user_id = $1`, userID)
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
) error {
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

	args = append(args, id, userID)

	query := fmt.Sprintf(
		`UPDATE social_links
		 SET %s
		 WHERE id = $%d AND user_id = $%d`,
		strings.Join(set, ", "),
		arg,
		arg+1,
	)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
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
