package users

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

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

	var request UpdateUserDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			user, err := h.repository.Update(
				r.Context(),
				userID,
				request,
			)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusOK,
				Message:    "user updated successfully",
				Data:       user,
			}, nil
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(
				w,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"user was not found",
			)
			return
		}

		log.Printf("Update error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"user could not be updated",
		)
		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)
}
