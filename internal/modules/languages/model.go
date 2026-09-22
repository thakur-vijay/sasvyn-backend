package languages

import "time"

type Language struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	LanguageCode string    `json:"language_code"`
	Language     string    `json:"language"`
	Proficiency  int16     `json:"proficiency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateLanguageDTO struct {
	LanguageCode string `json:"language_code" validate:"required"`
	Language     string `json:"language" validate:"required"`
	Proficiency  int16  `json:"proficiency" validate:"required,min=1,max=5"`
}

type UpdateLanguageDTO struct {
	LanguageCode *string `json:"language_code" validate:"omitempty,notblank"`
	Language     *string `json:"language" validate:"omitempty,notblank"`
	Proficiency  *int16  `json:"proficiency" validate:"omitempty,min=1,max=5"`
}
