package sociallinks

import "time"

type SocialLink struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSocialLinkDTO struct {
	Type string `json:"type" validate:"required"`
	Url  string `json:"url" validate:"required"`
}

type UpdateSocialLinkDTO struct {
	Type *string `json:"type" validate:"omitempty,notblank"`
	Url  *string `json:"url" validate:"omitempty,notblank"`
}
