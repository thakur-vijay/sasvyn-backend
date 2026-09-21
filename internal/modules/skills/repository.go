package skills

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, skill Skill) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateSkill: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO skills (id, user_id, skill, category, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`, skill.ID, skill.UserID, skill.Skill, skill.Category, skill.CreatedAt, skill.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]Skill, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchSkills: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, skill, category, created_at, updated_at FROM skills WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var skills []Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(
			&skill.ID,
			&skill.UserID,
			&skill.Skill,
			&skill.Category,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		); err != nil {
			return nil, err
		}
		skill.CreatedAt = skill.CreatedAt.UTC()
		skill.UpdatedAt = skill.UpdatedAt.UTC()
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *Repository) Update(ctx context.Context, skill Skill) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateSkill: %v", time.Since(start))
	}()

	_, err := r.db.ExecContext(ctx, `UPDATE skills SET skill = $1, category = $2, updated_at = $3 WHERE id = $4 AND user_id = $5`, skill.Skill, skill.Category, skill.UpdatedAt, skill.ID, skill.UserID)
	return err
}

func (r *Repository) Delete(ctx context.Context, skillID, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteSkill: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM skills
		WHERE id = $1
		  AND user_id = $2
	`, skillID, userID)
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
