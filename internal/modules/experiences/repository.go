package experiences

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

func (r *Repository) Create(ctx context.Context, experience Experience) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateExperience: %v", time.Since(start))
	}()
	_, err := r.db.ExecContext(ctx, `INSERT INTO experiences (id, user_id, role, company, start_date, end_date, is_currently_working, location, sync_version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, experience.ID, experience.UserID, experience.Role, experience.Company, experience.StartDate, experience.EndDate, experience.IsCurrentlyWorking, experience.Location, experience.SyncVersion, experience.CreatedAt, experience.UpdatedAt)
	return err
}

func (r *Repository) Fetch(ctx context.Context, userID string) ([]Experience, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchExperiences: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, role, company, start_date, end_date, is_currently_working, location, sync_version, created_at, updated_at FROM experiences WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var experiences []Experience
	for rows.Next() {
		var experience Experience
		if err := rows.Scan(
			&experience.ID,
			&experience.UserID,
			&experience.Role,
			&experience.Company,
			&experience.StartDate,
			&experience.EndDate,
			&experience.IsCurrentlyWorking,
			&experience.Location,
			&experience.SyncVersion,
			&experience.CreatedAt,
			&experience.UpdatedAt,
		); err != nil {
			return nil, err
		}
		experience.CreatedAt = experience.CreatedAt.UTC()
		experience.UpdatedAt = experience.UpdatedAt.UTC()
		experiences = append(experiences, experience)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return experiences, nil
}

func (r *Repository) FetchByID(ctx context.Context, id, userID string) (Experience, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchExperienceByID: %v", time.Since(start))
	}()

	var experience Experience

	err := r.db.QueryRowContext(ctx, `
	SELECT id, user_id, role, company, start_date, end_date, is_currently_working, location, sync_version, created_at, updated_at FROM experiences WHERE id = $1 AND user_id = $2
	`,
		id,
		userID,
	).Scan(
		&experience.ID,
		&experience.UserID,
		&experience.Role,
		&experience.Company,
		&experience.StartDate,
		&experience.EndDate,
		&experience.IsCurrentlyWorking,
		&experience.Location,
		&experience.SyncVersion,
		&experience.CreatedAt,
		&experience.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Experience{}, sql.ErrNoRows
		}

		return Experience{}, err
	}

	experience.CreatedAt = experience.CreatedAt.UTC()
	experience.UpdatedAt = experience.UpdatedAt.UTC()

	return experience, nil
}

func (r *Repository) Update(ctx context.Context, id string, userID string, request UpdateExperienceDTO) (Experience, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateExperience: %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if request.Role != nil {
		set = append(set, fmt.Sprintf("role = $%d", arg))
		args = append(args, request.Role)
		arg++
	}

	if request.Company != nil {
		set = append(set, fmt.Sprintf("company = $%d", arg))
		args = append(args, request.Company)
		arg++
	}

	if request.Location != nil {
		set = append(set, fmt.Sprintf("location = $%d", arg))
		args = append(args, request.Location)
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

	if request.IsCurrentlyWorking != nil {
		set = append(set, fmt.Sprintf("is_currently_working = $%d", arg))
		args = append(args, request.IsCurrentlyWorking)
		arg++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	set = append(set, "sync_version = sync_version + 1")
	args = append(args, id, userID)

	query := fmt.Sprintf(
		`UPDATE experiences
		 SET %s
		 WHERE id = $%d AND user_id = $%d
		 RETURNING
		 id, user_id, role, company, start_date, end_date, is_currently_working, location, sync_version, created_at, updated_at
		`,
		strings.Join(set, ", "),
		arg,
		arg+1,
	)

	var experience Experience

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&experience.ID,
		&experience.UserID,
		&experience.Role,
		&experience.Company,
		&experience.StartDate,
		&experience.EndDate,
		&experience.IsCurrentlyWorking,
		&experience.Location,
		&experience.SyncVersion,
		&experience.CreatedAt,
		&experience.UpdatedAt,
	)
	return experience, err
}

func (r *Repository) Delete(ctx context.Context, id, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteExperience: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM experiences
		WHERE id = $1
		  AND user_id = $2
	`, id, userID)
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

func (r *Repository) CreateResponsibility(ctx context.Context, responsibility ExperienceResponsibility, userID string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] CreateExperienceResponsibility: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO experience_responsibilities (
			id, experience_id, responsibility, "order",
			sync_version, created_at, updated_at
		)
		SELECT $1, e.id, $3, $4, $5, $6, $7
		FROM experiences e
		WHERE e.id = $2 AND e.user_id = $8
	`,
		responsibility.ID,
		responsibility.ExperienceID,
		responsibility.Responsibility,
		responsibility.Order,
		responsibility.SyncVersion,
		responsibility.CreatedAt,
		responsibility.UpdatedAt,
		userID,
	)
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

func (r *Repository) FetchResponsibilities(ctx context.Context, experienceID string) ([]ExperienceResponsibility, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchExperienceResponsibilities: %v", time.Since(start))
	}()

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			experience_id,
			responsibility,
			"order",
			sync_version,
			created_at,
			updated_at
		FROM experience_responsibilities
		WHERE experience_id = $1
		ORDER BY "order" ASC, created_at ASC
	`, experienceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responsibilities []ExperienceResponsibility

	for rows.Next() {
		var responsibility ExperienceResponsibility

		if err := rows.Scan(
			&responsibility.ID,
			&responsibility.ExperienceID,
			&responsibility.Responsibility,
			&responsibility.Order,
			&responsibility.SyncVersion,
			&responsibility.CreatedAt,
			&responsibility.UpdatedAt,
		); err != nil {
			return nil, err
		}

		responsibility.CreatedAt = responsibility.CreatedAt.UTC()
		responsibility.UpdatedAt = responsibility.UpdatedAt.UTC()

		responsibilities = append(responsibilities, responsibility)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return responsibilities, nil
}

func (r *Repository) FetchResponsibilityByID(ctx context.Context, experienceID, responsibilityID string) (ExperienceResponsibility, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] FetchExperienceResponsibilityByID: %v", time.Since(start))
	}()

	var responsibility ExperienceResponsibility

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			experience_id,
			responsibility,
			"order",
			sync_version,
			created_at,
			updated_at
		FROM experience_responsibilities
		WHERE experience_id = $1
		  AND id = $2
	`, experienceID, responsibilityID).Scan(
		&responsibility.ID,
		&responsibility.ExperienceID,
		&responsibility.Responsibility,
		&responsibility.Order,
		&responsibility.SyncVersion,
		&responsibility.CreatedAt,
		&responsibility.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ExperienceResponsibility{}, sql.ErrNoRows
		}

		return ExperienceResponsibility{}, err
	}

	responsibility.CreatedAt = responsibility.CreatedAt.UTC()
	responsibility.UpdatedAt = responsibility.UpdatedAt.UTC()

	return responsibility, nil
}

func (r *Repository) UpdateResponsibility(
	ctx context.Context,
	userID, experienceID, responsibilityID string,
	request UpdateExperienceResponsibilityDTO,
) (ExperienceResponsibility, error) {
	start := time.Now()
	defer func() {
		log.Printf("[DB] UpdateExperienceResponsibility: %v", time.Since(start))
	}()

	set := []string{}
	args := []any{}
	arg := 1

	if request.Responsibility != nil {
		set = append(set, fmt.Sprintf("responsibility = $%d", arg))
		args = append(args, *request.Responsibility)
		arg++
	}

	if request.Order != nil {
		set = append(set, fmt.Sprintf(`"order" = $%d`, arg))
		args = append(args, *request.Order)
		arg++
	}

	if len(set) == 0 {
		return r.FetchResponsibilityByID(ctx, experienceID, responsibilityID)
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", arg))
	args = append(args, time.Now().UTC())
	arg++

	set = append(set, "sync_version = sync_version + 1")

	args = append(args, experienceID, responsibilityID, userID)

	query := fmt.Sprintf(`
		UPDATE experience_responsibilities
		SET %s
		WHERE experience_id = $%d
		  AND id = $%d
		  AND EXISTS (
			  SELECT 1
			  FROM experiences e
			  WHERE e.id = experience_responsibilities.experience_id
			    AND e.user_id = $%d
		  )
		RETURNING
			id,
			experience_id,
			responsibility,
			"order",
			sync_version,
			created_at,
			updated_at
	`, strings.Join(set, ", "), arg, arg+1, arg+2)

	var responsibility ExperienceResponsibility

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&responsibility.ID,
		&responsibility.ExperienceID,
		&responsibility.Responsibility,
		&responsibility.Order,
		&responsibility.SyncVersion,
		&responsibility.CreatedAt,
		&responsibility.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ExperienceResponsibility{}, sql.ErrNoRows
		}
		return ExperienceResponsibility{}, err
	}

	responsibility.CreatedAt = responsibility.CreatedAt.UTC()
	responsibility.UpdatedAt = responsibility.UpdatedAt.UTC()

	return responsibility, nil
}

func (r *Repository) DeleteResponsibility(
	ctx context.Context,
	userID, experienceID, responsibilityID string,
) error {
	start := time.Now()
	defer func() {
		log.Printf("[DB] DeleteExperienceResponsibility: %v", time.Since(start))
	}()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM experience_responsibilities
		WHERE experience_id = $1
		  AND id = $2
		  AND EXISTS (
			  SELECT 1
			  FROM experiences e
			  WHERE e.id = experience_responsibilities.experience_id
			    AND e.user_id = $3
		  )
	`, experienceID, responsibilityID, userID)
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
