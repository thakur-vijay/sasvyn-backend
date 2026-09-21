package skills

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

	var request CreateSkillDTO

	if err := response.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	userID, _ := auth.UserID(r.Context())
	now := time.Now().UTC()
	skill := Skill{
		ID:        uuid.New().String(),
		UserID:    userID,
		Skill:     request.Skill,
		Category:  request.Category,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.repository.Create(r.Context(), skill); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create skill")
	}

	response.Write(w, http.StatusCreated, "Skill added successfully", skill)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	skills, err := h.repository.Fetch(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create skill")
	}

	response.Write(w, http.StatusOK, "Skill fetched successfully", skills)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	var request UpdateSkillDTO

	if err := response.DecodeJSON(r, &request); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	skill := Skill{
		ID:        skillID,
		UserID:    userID,
		Skill:     request.Skill,
		Category:  request.Category,
		UpdatedAt: time.Now().UTC(),
	}

	if err := h.repository.Update(r.Context(), skill); err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "skill was not found")
			return
		}

		if strings.Contains(err.Error(), "23505") {

			response.WriteError(w, http.StatusConflict, "skill already exists")
			return
		}

		log.Printf("Update skill error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "skill could not be updated")
		return
	}

	response.Write(w, http.StatusOK, "skill updated successfully", skill)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	if err := h.repository.Delete(r.Context(), skillID, userID); err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "skill was not found")
			return
		}

		log.Printf("Update skill error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "skill could not be deleted")
		return
	}

	response.Write(w, http.StatusOK, "skill deleted successfully", nil)
}
