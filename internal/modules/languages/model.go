package languages

import (
	"github.com/sasvyn/backend/internal/domain/model"
)

type Language struct {
	model.BaseModel
	UserID       string `json:"user_id"`
	LanguageCode string `json:"language_code"`
	Language     string `json:"language"`
	Proficiency  int16  `json:"proficiency"`
}

type CreateLanguageDTO struct {
	ID           string `json:"id" validate:"required,valid_uuid4"`
	LanguageCode string `json:"language_code" validate:"required"`
	Language     string `json:"language" validate:"required"`
	Proficiency  int16  `json:"proficiency" validate:"required,min=1,max=5"`
}

type UpdateLanguageDTO struct {
	LanguageCode *string `json:"language_code" validate:"omitempty,notblank"`
	Language     *string `json:"language" validate:"omitempty,notblank"`
	Proficiency  *int16  `json:"proficiency" validate:"omitempty,min=1,max=5"`
}
