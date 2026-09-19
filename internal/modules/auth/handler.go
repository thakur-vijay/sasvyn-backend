package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	service        *Service
	sessionService *sessions.Service
}

func NewHandler(service *Service, sessionService *sessions.Service) *Handler {
	return &Handler{
		service:        service,
		sessionService: sessionService,
	}
}

func (h *Handler) SocialLogin(w http.ResponseWriter, r *http.Request) {
	var request SocialLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.WriteError(w, http.StatusBadRequest, "request body is invalid")
		return
	}

	user, accessToken, refreshToken, err := h.service.SocialLogin(r.Context(), request)
	if err != nil {
		log.Printf("Social login error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "login could not be completed")
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

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := SessionID(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	if err := h.sessionService.Logout(r.Context(), sessionID); err != nil {
		log.Printf("Logout error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "logout failed")
		return
	}

	response.Write(w, http.StatusOK, "logout successful", nil)
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
