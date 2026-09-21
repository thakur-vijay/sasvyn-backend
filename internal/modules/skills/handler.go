package skills

import (
	"net/http"
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
