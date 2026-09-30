package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"uuid"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository         *Repository
	PresignClient      *s3.PresignClient
	idempotencyService *idempotency.Service
}

func NewHandler(
	repository *Repository,
	presignClient *s3.PresignClient,
	idempotencyService *idempotency.Service,
) *Handler {
	return &Handler{
		repository:         repository,
		PresignClient:      presignClient,
		idempotencyService: idempotencyService,
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	user, err := h.repository.GetByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "user was not found")
			return
		}

		log.Printf("GetByID error: %v", err)
		response.Write(w, http.StatusInternalServerError, "user could not be retrieved")
		return
	}
	response.WriteItem(w, http.StatusOK, "user retrieved successfully", user)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	idempotencyKey := r.Header.Get("Idempotency-Key")

	if idempotencyKey == "" {
		response.Write(w, http.StatusBadRequest, "Idempotency-Key header is required")
		return
	}
	var request UpdateUserDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	record := idempotency.IdempotencyRecord{
		ID:             uuid.New().String(),
		UserID:         userID,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now().UTC(),
	}

	reserved, err := h.idempotencyService.Reserve(
		r.Context(),
		record,
	)
	if err != nil {
		log.Printf("Idempotency reserve error: %v", err)
		response.Write(w, http.StatusInternalServerError, "request could not be processed")
		return
	}

	if reserved.Status == "completed" {
		response.WriteRaw(
			w,
			*reserved.StatusCode,
			reserved.ResponseBody,
		)
		return
	}

	if reserved.Status == "processing" {
		response.Write(
			w,
			http.StatusConflict,
			"request with this Idempotency-Key is already being processed",
		)
		return
	}

	user, err := h.repository.Update(r.Context(), userID, request)

	if err != nil {
		_ = h.idempotencyService.Delete(
			r.Context(),
			userID,
			idempotencyKey,
		)

		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "user was not found")
			return
		}

		log.Printf("Update error: %v", err)
		response.Write(w, http.StatusInternalServerError, "user could not be updated")
		return
	}

	responseBody, err := json.Marshal(
		response.Response{
			StatusCode: http.StatusOK,
			Message:    "user updated successfully",
			Data:       user,
		},
	)
	if err != nil {
		log.Printf("Response marshal error: %v", err)
		response.Write(w, http.StatusInternalServerError, "user could not be updated")
		return
	}

	if err := h.idempotencyService.Complete(
		r.Context(),
		userID,
		idempotencyKey,
		http.StatusOK,
		responseBody,
	); err != nil {
		log.Printf("Idempotency complete error: %v", err)
		response.Write(w, http.StatusInternalServerError, "user could not be updated")
		return
	}
	response.WriteItem(w, http.StatusOK, "user updated successfully", user)
}
