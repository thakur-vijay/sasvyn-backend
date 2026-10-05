package skills

import (
	"github.com/sasvyn/backend/internal/domain/model"
)

type Skill struct {
	model.BaseModel
	UserID   string `json:"user_id"`
	Skill    string `json:"skill"`
	Category string `json:"category"`
}

type CreateSkillDTO struct {
	ID       string `json:"id" validate:"required,valid_uuid4"`
	Skill    string `json:"skill" validate:"required"`
	Category string `json:"category" validate:"required"`
}

type UpdateSkillDTO struct {
	Skill    *string `json:"skill" validate:"omitempty,notblank"`
	Category *string `json:"category" validate:"omitempty,notblank"`
}
