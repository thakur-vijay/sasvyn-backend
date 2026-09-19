package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (h *Handler) SocialLogin(w http.ResponseWriter, r *http.Request) {
	var request SocialLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.WriteError(w, http.StatusBadRequest, "request body is invalid")
		return
	}

	existingUser, err := h.repository.GetByAppleID(r.Context(), request.AppleID)
	if err == nil {
		accessToken, refreshToken, err := h.sessionService.CreateSession(r.Context(), existingUser.ID)
		if err != nil {
			log.Printf("Create Session error: %v", err)
			response.WriteError(w, http.StatusInternalServerError, "session could not be created")
			return
		}
		payload := SocialLoginResponse{
			User:         *existingUser,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		response.Write(w, http.StatusOK, "login successful", payload)
		return
	}

	if err != sql.ErrNoRows {
		log.Printf("GetByAppleID error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "user could not be retrieved")
		return
	}

	now := time.Now().UTC()

	user := UserResponse{
		ID:        uuid.NewString(),
		AppleID:   request.AppleID,
		FullName:  request.FullName,
		Email:     request.Email,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	if err := h.repository.Create(r.Context(), user); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "user could not be created")
		return
	}

	accessToken, refreshToken, err := h.sessionService.CreateSession(r.Context(), user.ID)
	if err != nil {
		log.Printf("Create Session error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "session could not be created")
		return
	}

	payload := SocialLoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.Write(w, http.StatusOK, "login successful", payload)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.WriteError(w, http.StatusBadRequest, "request body is invalid")
		return
	}

	accessToken, refreshToken, err := h.sessionService.RefreshSession(
		r.Context(),
		request.RefreshToken,
	)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "refresh token is invalid")
		return
	}

	payload := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.Write(w, http.StatusOK, "token refreshed successfully", payload)
}

// func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := UserID(r.Context())
// 	if !ok {
// 		response.WriteError(w, http.StatusUnauthorized, "authentication is required")
// 		return
// 	}

// 	user, err := h.repository.GetByID(r.Context(), userID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			response.WriteError(w, http.StatusNotFound, "user was not found")
// 			return
// 		}

// 		log.Printf("GetByID error: %v", err)
// 		response.WriteError(w, http.StatusInternalServerError, "user could not be retrieved")
// 		return
// 	}

// 	response.Write(w, http.StatusOK, "user retrieved successfully", user)
// }
