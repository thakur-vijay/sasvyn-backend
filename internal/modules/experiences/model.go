package experiences

import (
	"time"

	"github.com/sasvyn/backend/internal/domain/model"
)

type Experience struct {
	model.BaseModel
	UserID             string                     `json:"user_id"`
	Role               string                     `json:"role"`
	Company            string                     `json:"company"`
	StartDate          time.Time                  `json:"start_date"`
	EndDate            *time.Time                 `json:"end_date"`
	IsCurrentlyWorking bool                       `json:"is_currently_working"`
	Location           string                     `json:"location"`
	Responsibilities   []ExperienceResponsibility `json:"responsibilities"`
}

type ExperienceResponsibility struct {
	model.BaseModel
	ExperienceID   string `json:"experience_id"`
	Responsibility string `json:"responsibility"`
	Order          int    `json:"order"`
}

type CreateExperienceDTO struct {
	ID                 string     `json:"id" validate:"required,valid_uuid4"`
	Role               string     `json:"role" validate:"required"`
	Company            string     `json:"company" validate:"required"`
	Location           string     `json:"location" validate:"required"`
	StartDate          time.Time  `json:"start_date" validate:"required"`
	EndDate            *time.Time `json:"end_date" validate:"omitempty"`
	IsCurrentlyWorking bool       `json:"is_currently_working"`
}

type UpdateExperienceDTO struct {
	Role               *string    `json:"role" validate:"omitempty,notblank"`
	Company            *string    `json:"company" validate:"omitempty,notblank"`
	Location           *string    `json:"location" validate:"omitempty,notblank"`
	StartDate          *time.Time `json:"start_date" validate:"omitempty"`
	EndDate            *time.Time `json:"end_date" validate:"omitempty"`
	IsCurrentlyWorking *bool      `json:"is_currently_working" validate:"omitempty"`
}

type CreateExperienceResponsibilityDTO struct {
	ID             string `json:"id" validate:"required,valid_uuid4"`
	Responsibility string `json:"responsibility" validate:"required"`
	Order          int    `json:"order" validate:"gte=0"`
}

type UpdateExperienceResponsibilityDTO struct {
	Responsibility *string `json:"responsibility" validate:"omitempty,notblank"`
	Order          *int    `json:"order" validate:"omitempty,gte=0"`
}
