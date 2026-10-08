package education

import (
	"time"

	"github.com/sasvyn/backend/internal/domain/model"
)

type Education struct {
	model.BaseModel
	UserID       string     `json:"user_id"`
	Degree       string     `json:"degree"`
	FieldOfStudy string     `json:"field_of_study"`
	Institution  string     `json:"institution"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsPursuing   bool       `json:"is_pursuing"`
	Grade        string     `json:"grade"`
	GradeType    string     `json:"grade_type"`
	Description  *string    `json:"description"`
}

type CreateEducationDTO struct {
	ID           string     `json:"id" validate:"required,valid_uuid4"`
	Degree       string     `json:"degree" validate:"required"`
	FieldOfStudy string     `json:"field_of_study" validate:"required"`
	Institution  string     `json:"institution" validate:"required"`
	StartDate    time.Time  `json:"start_date" validate:"required"`
	EndDate      *time.Time `json:"end_date" validate:"omitempty"`
	IsPursuing   bool       `json:"is_pursuing"`
	Grade        string     `json:"grade" validate:"required"`
	GradeType    string     `json:"grade_type" validate:"required"`
	Description  *string    `json:"description" validate:"omitempty,notblank"`
}

type UpdateEducationDTO struct {
	Degree       *string    `json:"degree" validate:"omitempty,notblank"`
	FieldOfStudy *string    `json:"field_of_study" validate:"omitempty,notblank"`
	Institution  *string    `json:"institution" validate:"omitempty,notblank"`
	StartDate    *time.Time `json:"start_date" validate:"omitempty"`
	EndDate      *time.Time `json:"end_date" validate:"omitempty"`
	IsPursuing   *bool      `json:"is_pursuing" validate:"omitempty,notblank"`
	Grade        *string    `json:"grade" validate:"omitempty,notblank"`
	GradeType    *string    `json:"grade_type" validate:"omitempty,notblank"`
	Description  *string    `json:"description" validate:"omitempty,notblank"`
}
