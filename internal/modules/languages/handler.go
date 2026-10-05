package languages

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository         *Repository
	idempotencyService *idempotency.Service
}

func NewHandler(repository *Repository, idempotencyService *idempotency.Service) *Handler {
	return &Handler{repository: repository, idempotencyService: idempotencyService}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")

	var request CreateLanguageDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		now := time.Now().UTC()
		language := Language{
			ID:           uuid.New().String(),
			UserID:       userID,
			LanguageCode: request.LanguageCode,
			Language:     request.Language,
			Proficiency:  request.Proficiency,
			SyncVersion:  1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := h.repository.Create(r.Context(), language); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusCreated,
			Message:    "Language added successfully",
			Data:       language,
		}, nil

	},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w, http.StatusNotFound,
				"Language was not found",
			)
			return
		}

		log.Printf("Create language error: %v", err)
		response.Write(w, http.StatusInternalServerError, "language could not be created")
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	languages, err := h.repository.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch languages")
		return
	}

	response.WriteList(w, http.StatusOK, "Languages fetched successfully", languages)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	languageID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	var request UpdateLanguageDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	// if request.LanguageCode == nil && request.Language == nil && request.Proficiency == nil {
	// 	response.Write(w, http.StatusBadRequest, "at least one field is required")
	// 	return
	// }

	if err := h.repository.Update(
		r.Context(),
		languageID,
		userID,
		request,
	); err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "language was not found")
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(w, http.StatusConflict, "language already exists")
			return
		}

		log.Printf("Update language error: %v", err)
		response.Write(w, http.StatusInternalServerError, "language could not be updated")
		return
	}

	response.Write(w, http.StatusOK, "language updated successfully")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	languageID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	if err := h.repository.Delete(r.Context(), languageID, userID); err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "language was not found")
			return
		}

		log.Printf("Update language error: %v", err)
		response.Write(w, http.StatusInternalServerError, "language could not be deleted")
		return
	}

	response.Write(w, http.StatusOK, "language deleted successfully")
}
