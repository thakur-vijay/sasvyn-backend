package skills

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

func (r *Repository) Create(ctx context.Context, skill Skill) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateSkill: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO skills (id, user_id, skill, category, sync_version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, skill.ID, skill.UserID, skill.Skill, skill.Category, skill.SyncVersion, skill.CreatedAt, skill.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]Skill, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchSkills: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, skill, category, sync_version, created_at, updated_at FROM skills WHERE user_id = $1 ORDER BY created_at DESC`, userID)
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
			&skill.SyncVersion,
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

func (r *Repository) FetchByID(ctx context.Context, id, userID string) (Skill, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchSkillByID: %v", time.Since(start))
	}()

	var skill Skill

	err := r.db.QueryRowContext(ctx, `
	SELECT id, user_id, skill, category, sync_version, created_at, updated_at FROM skills WHERE id = $1 AND user_id = $2
	`,
		id,
		userID,
	).Scan(
		&skill.ID,
		&skill.UserID,
		&skill.Skill,
		&skill.Category,
		&skill.SyncVersion,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Skill{}, sql.ErrNoRows
		}

		return Skill{}, err
	}

	skill.CreatedAt = skill.CreatedAt.UTC()
	skill.UpdatedAt = skill.UpdatedAt.UTC()

	return skill, nil
}

func (r *Repository) Update(ctx context.Context, id string, userID string, request UpdateSkillDTO) (Skill, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateSkill: %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if request.Skill != nil {
		set = append(set, fmt.Sprintf("skill = $%d", arg))
		args = append(args, request.Skill)
		arg++
	}

	if request.Category != nil {
		set = append(set, fmt.Sprintf("category = $%d", arg))
		args = append(args, request.Category)
		arg++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	set = append(set, "sync_version = sync_version + 1")
	args = append(args, id, userID)

	query := fmt.Sprintf(
		`UPDATE skills
		 SET %s
		 WHERE id = $%d AND user_id = $%d
		 RETURNING
		 id, 
		 user_id,
		 skill,
		 category,
		 sync_version,
		 created_at,
		 updated_at
		`,
		strings.Join(set, ", "),
		arg,
		arg+1,
	)

	var updatedSkill Skill

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&updatedSkill.ID,
		&updatedSkill.UserID,
		&updatedSkill.Skill,
		&updatedSkill.Category,
		&updatedSkill.SyncVersion,
		&updatedSkill.CreatedAt,
		&updatedSkill.UpdatedAt,
	)
	return updatedSkill, err
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
