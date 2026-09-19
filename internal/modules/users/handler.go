package users

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository     *Repository
	sessionService *sessions.Service
}

func NewHandler(repository *Repository, sessionService *sessions.Service) *Handler {
	return &Handler{
		repository:     repository,
		sessionService: sessionService,
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	user, err := h.repository.GetByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "user was not found")
			return
		}

		log.Printf("GetByID error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "user could not be retrieved")
		return
	}

	response.Write(w, http.StatusOK, "user retrieved successfully", user)
}
