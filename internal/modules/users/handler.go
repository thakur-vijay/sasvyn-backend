package users

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/sasvyn/backend/internal/response"
	"github.com/sasvyn/backend/internal/validation"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{
		repository: repository,
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

	var request UpdateUserDTO

	if err := response.DecodeJSON(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, "request body is invalid")
		return
	}

	if request.DateOfBirth != nil && !validation.IsISODate(*request.DateOfBirth) {
		response.Write(
			w,
			http.StatusBadRequest,
			"date_of_birth must be a valid date in YYYY-MM-DD format",
		)
		return
	}

	user, err := h.repository.Update(r.Context(), userID, request)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "user was not found")
			return
		}

		log.Printf("Update error: %v", err)
		response.Write(w, http.StatusInternalServerError, "user could not be updated")
		return
	}

	response.WriteItem(w, http.StatusOK, "user updated successfully", user)
}
