package upload

type CreateUploadURLRequest struct {
	Type        string `json:"type" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}

type CreateUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
}
