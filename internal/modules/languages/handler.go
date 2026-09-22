package languages

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var request CreateLanguageDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}
	userID, _ := auth.UserID(r.Context())
	now := time.Now().UTC()
	language := Language{
		ID:           uuid.New().String(),
		UserID:       userID,
		LanguageCode: request.LanguageCode,
		Language:     request.Language,
		Proficiency:  request.Proficiency,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := h.repository.Create(r.Context(), language); err != nil {
		response.Write(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteItem(w, http.StatusCreated, "Language added successfully", language)
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

	if request.LanguageCode == nil && request.Language == nil && request.Proficiency == nil {
		response.Write(w, http.StatusBadRequest, "at least one field is required")
		return
	}

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
