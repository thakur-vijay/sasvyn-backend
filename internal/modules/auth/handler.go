package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	start := time.Now()

	defer func() {
		log.Printf("[SocialLogin] TOTAL: %v", time.Since(start))
	}()

	var request SocialLoginRequest

	if err := response.DecodeJSON(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	user, accessToken, refreshToken, err := h.service.SocialLogin(r.Context(), request)
	if err != nil {
		log.Printf("Social login error: %v", err)
		response.Write(w, http.StatusInternalServerError, "login could not be completed")
		return
	}

	payload := SocialLoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.WriteItem(w, http.StatusOK, "login successful", payload)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Write(w, http.StatusBadRequest, "request body is invalid")
		return
	}

	accessToken, refreshToken, err := h.sessionService.RefreshSession(
		r.Context(),
		request.RefreshToken,
	)
	if err != nil {
		response.Write(w, http.StatusUnauthorized, "refresh token is invalid")
		return
	}

	payload := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.WriteItem(w, http.StatusOK, "token refreshed successfully", payload)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := SessionID(r.Context())
	if !ok {
		response.Write(w, http.StatusUnauthorized, "authentication required")
		return
	}

	if err := h.sessionService.Logout(r.Context(), sessionID); err != nil {
		log.Printf("Logout error: %v", err)
		response.Write(w, http.StatusInternalServerError, "logout failed")
		return
	}

	response.Write(w, http.StatusOK, "logout successful")
}
