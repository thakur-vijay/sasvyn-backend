package sociallinks

import (
	"github.com/sasvyn/backend/internal/domain/model"
)

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
	model.BaseModel
	UserID string   `json:"user_id"`
	Type   LinkType `json:"type"`
	Url    string   `json:"url"`
}

type CreateSocialLinkDTO struct {
	ID   string   `json:"id" validate:"required,valid_uuid4"`
	Type LinkType `json:"type" validate:"required,oneof=github linkedin x instagram youtube dribbble behance medium website"`
	Url  string   `json:"url" validate:"required,url"`
}

type UpdateSocialLinkDTO struct {
	Type *LinkType `json:"type" validate:"omitempty,oneof=github linkedin x instagram youtube dribbble behance medium website"`
	Url  *string   `json:"url" validate:"omitempty,notblank,url"`
}
