package auth

import (
	"encoding/json"
	"net/http"

	"github.com/sasvyn/backend/internal/sessions"
)

type Handler struct {
	sessionService *sessions.Service
}

func NewHandler(sessionService *sessions.Service) *Handler {
	return &Handler{
		sessionService: sessionService,
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.sessionService.RefreshSession(
		r.Context(),
		request.RefreshToken,
	)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	response := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
