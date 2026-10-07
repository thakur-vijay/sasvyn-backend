package documents

import (
	"github.com/sasvyn/backend/internal/domain/model"
)

type Document struct {
	model.BaseModel
	UserID   string  `json:"user_id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	FileSize int64   `json:"file_size"`
	Url      *string `json:"url"`
	Key      *string `json:"-"`
}

type CreateDocumentDTO struct {
	ID       string `json:"id" validate:"required,valid_uuid4"`
	Name     string `json:"name" validate:"required"`
	Category string `json:"category" validate:"required"`
	FileSize string `json:"file_size" validate:"required"`
	Key      string `json:"key" validate:"required"`
}
