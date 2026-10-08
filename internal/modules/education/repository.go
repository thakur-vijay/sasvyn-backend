package education

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

func (r *Repository) Create(ctx context.Context, education Education) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateEducation: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO educations (id, user_id, degree, field_of_study, institution, start_date, end_date, is_pursuing, grade, grade_type, description, sync_version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, education.ID, education.UserID, education.Degree, education.FieldOfStudy, education.Institution, education.StartDate, education.EndDate, education.IsPursuing, education.Grade, education.GradeType, education.Description, education.SyncVersion, education.CreatedAt, education.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]Education, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchEducations: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, degree, field_of_study, institution, start_date, end_date, is_pursuing, grade, grade_type, description, sync_version, created_at, updated_at FROM educations WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var educations []Education
	for rows.Next() {
		var education Education
		if err := rows.Scan(
			&education.ID,
			&education.UserID,
			&education.Degree,
			&education.FieldOfStudy,
			&education.Institution,
			&education.StartDate,
			&education.EndDate,
			&education.IsPursuing,
			&education.Grade,
			&education.GradeType,
			&education.Description,
			&education.SyncVersion,
			&education.CreatedAt,
			&education.UpdatedAt,
		); err != nil {
			return nil, err
		}
		education.CreatedAt = education.CreatedAt.UTC()
		education.UpdatedAt = education.UpdatedAt.UTC()
		educations = append(educations, education)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return educations, nil
}

func (r *Repository) FetchByID(ctx context.Context, id, userID string) (Education, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchEducationByID: %v", time.Since(start))
	}()

	var education Education

	err := r.db.QueryRowContext(ctx, `
	SELECT id, user_id, degree, field_of_study, institution, start_date, end_date, is_pursuing, grade, grade_type, description, sync_version, created_at, updated_at FROM educations WHERE id = $1 AND user_id = $2
	`,
		id,
		userID,
	).Scan(
		&education.ID,
		&education.UserID,
		&education.Degree,
		&education.FieldOfStudy,
		&education.Institution,
		&education.StartDate,
		&education.EndDate,
		&education.IsPursuing,
		&education.Grade,
		&education.GradeType,
		&education.Description,
		&education.SyncVersion,
		&education.CreatedAt,
		&education.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Education{}, sql.ErrNoRows
		}

		return Education{}, err
	}

	education.CreatedAt = education.CreatedAt.UTC()
	education.UpdatedAt = education.UpdatedAt.UTC()

	return education, nil
}

func (r *Repository) Update(ctx context.Context, id string, userID string, request UpdateEducationDTO) (Education, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateEducation: %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if request.Degree != nil {
		set = append(set, fmt.Sprintf("degree = $%d", arg))
		args = append(args, request.Degree)
		arg++
	}

	if request.FieldOfStudy != nil {
		set = append(set, fmt.Sprintf("field_of_study = $%d", arg))
		args = append(args, request.FieldOfStudy)
		arg++
	}

	if request.Institution != nil {
		set = append(set, fmt.Sprintf("institution = $%d", arg))
		args = append(args, request.Institution)
		arg++
	}

	if request.StartDate != nil {
		set = append(set, fmt.Sprintf("start_date = $%d", arg))
		args = append(args, request.StartDate)
		arg++
	}

	if request.EndDate != nil {
		set = append(set, fmt.Sprintf("end_date = $%d", arg))
		args = append(args, request.EndDate)
		arg++
	}

	if request.IsPursuing != nil {
		set = append(set, fmt.Sprintf("is_pursuing = $%d", arg))
		args = append(args, request.IsPursuing)
		arg++
	}

	if request.Grade != nil {
		set = append(set, fmt.Sprintf("grade = $%d", arg))
		args = append(args, request.Grade)
		arg++
	}

	if request.GradeType != nil {
		set = append(set, fmt.Sprintf("grade_type = $%d", arg))
		args = append(args, request.GradeType)
		arg++
	}

	if request.Description != nil {
		set = append(set, fmt.Sprintf("description = $%d", arg))
		args = append(args, request.Description)
		arg++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	set = append(set, "sync_version = sync_version + 1")
	args = append(args, id, userID)

	query := fmt.Sprintf(
		`UPDATE educations
		 SET %s
		 WHERE id = $%d AND user_id = $%d
		 RETURNING
		 id, user_id, degree, field_of_study, institution, start_date, end_date, is_pursuing, grade, grade_type, description, sync_version, created_at, updated_at
		`,
		strings.Join(set, ", "),
		arg,
		arg+1,
	)

	var education Education

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&education.ID,
		&education.UserID,
		&education.Degree,
		&education.FieldOfStudy,
		&education.Institution,
		&education.StartDate,
		&education.EndDate,
		&education.IsPursuing,
		&education.Grade,
		&education.GradeType,
		&education.Description,
		&education.SyncVersion,
		&education.CreatedAt,
		&education.UpdatedAt,
	)
	return education, err
}

func (r *Repository) Delete(ctx context.Context, educationID, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteEducation: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM educations
		WHERE id = $1
		  AND user_id = $2
	`, educationID, userID)
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
