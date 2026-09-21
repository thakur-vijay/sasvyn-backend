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

		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return skills, nil
}
