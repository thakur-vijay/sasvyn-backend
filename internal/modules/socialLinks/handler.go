package sociallinks

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

	var request CreateSocialLinkDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}
	userID, _ := auth.UserID(r.Context())
	now := time.Now().UTC()
	link := SocialLink{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      request.Type,
		Url:       request.Url,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.repository.Create(r.Context(), link); err != nil {
		response.Write(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteItem(w, http.StatusCreated, "Social Link added successfully", link)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	socialLinks, err := h.repository.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch social links")
		return
	}

	response.WriteList(w, http.StatusOK, "Social Links fetched successfully", socialLinks)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	socialLinkId := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	var request UpdateSocialLinkDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repository.Update(
		r.Context(),
		socialLinkId,
		userID,
		request,
	); err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "social link was not found")
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(w, http.StatusConflict, "social link already exists")
			return
		}

		log.Printf("Update language error: %v", err)
		response.Write(w, http.StatusInternalServerError, "social link could not be updated")
		return
	}

	response.Write(w, http.StatusOK, "social link updated successfully")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	socialLinkId := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	if err := h.repository.Delete(r.Context(), socialLinkId, userID); err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "social link was not found")
			return
		}

		log.Printf("Update social link error: %v", err)
		response.Write(w, http.StatusInternalServerError, "social link could not be deleted")
		return
	}

	response.Write(w, http.StatusOK, "social link deleted successfully")
}
