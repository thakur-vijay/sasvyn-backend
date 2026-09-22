package users

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository    *Repository
	PresignClient *s3.PresignClient
}

func NewHandler(
	repository *Repository,
	presignClient *s3.PresignClient,
) *Handler {
	return &Handler{
		repository:    repository,
		PresignClient: presignClient,
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

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
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
