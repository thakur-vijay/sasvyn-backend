package sociallinks

import "time"

type LinkType string

const (
	LinkTypeGitHub    LinkType = "github"
	LinkTypeLinkedIn  LinkType = "linkedin"
	LinkTypeX         LinkType = "x"
	LinkTypeInstagram LinkType = "instagram"
	LinkTypeYouTube   LinkType = "youtube"
	LinkTypeDribbble  LinkType = "dribbble"
	LinkTypeBehance   LinkType = "behance"
	LinkTypeMedium    LinkType = "medium"
	LinkTypeWebsite   LinkType = "website"
)

type SocialLink struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        LinkType  `json:"type"`
	Url         string    `json:"url"`
	SyncVersion *int64    `json:"sync_version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSocialLinkDTO struct {
	Type LinkType `json:"type" validate:"required,oneof=github linkedin x instagram youtube dribbble behance medium website"`
	Url  string   `json:"url" validate:"required,url"`
}

type UpdateSocialLinkDTO struct {
	Type *LinkType `json:"type" validate:"omitempty,oneof=github linkedin x instagram youtube dribbble behance medium website"`
	Url  *string   `json:"url" validate:"omitempty,notblank,url"`
}
