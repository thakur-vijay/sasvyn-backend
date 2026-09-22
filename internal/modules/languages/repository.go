package languages

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

func (r *Repository) Create(ctx context.Context, language Language) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateLanguage: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO languages (id, user_id, language_code, language, proficiency, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, language.ID, language.UserID, language.LanguageCode, language.Language, language.Proficiency, language.CreatedAt, language.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]Language, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchLanguages: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, language_code, language, proficiency, created_at, updated_at FROM languages WHERE user_id = $1 ORDER BY proficiency DESC, language ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var languages []Language
	for rows.Next() {
		var language Language
		if err := rows.Scan(
			&language.ID,
			&language.UserID,
			&language.LanguageCode,
			&language.Language,
			&language.Proficiency,
			&language.CreatedAt,
			&language.UpdatedAt,
		); err != nil {
			return nil, err
		}
		language.CreatedAt = language.CreatedAt.UTC()
		language.UpdatedAt = language.UpdatedAt.UTC()
		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return languages, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	userID string,
	request UpdateLanguageDTO,
) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateLanguage: %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if request.LanguageCode != nil {
		set = append(set, fmt.Sprintf("language_code = $%d", arg))
		args = append(args, *request.LanguageCode)
		arg++
	}

	if request.Language != nil {
		set = append(set, fmt.Sprintf("language = $%d", arg))
		args = append(args, *request.Language)
		arg++
	}

	if request.Proficiency != nil {
		set = append(set, fmt.Sprintf("proficiency = $%d", arg))
		args = append(args, *request.Proficiency)
		arg++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	args = append(args, id, userID)

	query := fmt.Sprintf(
		`UPDATE languages
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

func (r *Repository) Delete(ctx context.Context, languageID, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteLanguage: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM languages
		WHERE id = $1
		  AND user_id = $2
	`, languageID, userID)
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
