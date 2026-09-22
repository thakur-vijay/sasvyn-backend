package upload

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	R2Client      *s3.Client
	PresignClient *s3.PresignClient
}

func NewHandler(r2Client *s3.Client, presignClient *s3.PresignClient) *Handler {
	return &Handler{
		R2Client:      r2Client,
		PresignClient: presignClient,
	}
}

func (h *Handler) CreateUploadURL(w http.ResponseWriter, r *http.Request) {
	var request CreateUploadURLRequest

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	switch request.Type {
	case UploadTypeProfileImage, UploadTypeProjectScreenshot, UploadTypeAppIcon:
	default:
		response.Write(w, http.StatusBadRequest, "invalid upload type")
		return
	}

	switch request.ContentType {
	case "image/jpeg", "image/png", "image/webp":
	default:
		response.Write(w, http.StatusBadRequest, "invalid content type")
		return
	}

	userID, _ := auth.UserID(r.Context())

	key := fmt.Sprintf(
		"users/%s/%s",
		userID,
		uuid.New().String(),
	)

	bucket := os.Getenv("R2_BUCKET_NAME")

	presigned, err := h.PresignClient.PresignPutObject(
		r.Context(),
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			ContentType: aws.String(request.ContentType),
		},
		func(options *s3.PresignOptions) {
			options.Expires = 10 * time.Minute
		},
	)

	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to generate upload URL")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Upload URL generated successfully",
		CreateUploadURLResponse{
			UploadURL: presigned.URL,
			ImgKey:    key,
		},
	)
}
