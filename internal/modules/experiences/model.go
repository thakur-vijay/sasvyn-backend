package experiences

import (
	"time"

	"github.com/sasvyn/backend/internal/domain/model"
)

type Experience struct {
	model.BaseModel
	UserID             string     `json:"user_id"`
	Role               string     `json:"role"`
	Company            string     `json:"company"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            *time.Time `json:"end_date"`
	IsCurrentlyWorking bool       `json:"is_currently_working"`
	Location           string     `json:"location"`
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

type CreateExperienceResponsibilityDTO struct {
	ID             string `json:"id" validate:"required,valid_uuid4"`
	ExperienceID   string `json:"experience_id" validate:"required,valid_uuid4"`
	Responsibility string `json:"responsibility" validate:"required"`
	Order          int    `json:"order" validate:"gte=0"`
}
