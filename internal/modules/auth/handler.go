package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	service            *Service
	sessionService     *sessions.Service
	clientDisconnecter ClientDisconnecter
}

type ClientDisconnecter interface {
	DisconnectClient(userID, clientID string) error
}

func NewHandler(
	service *Service,
	sessionService *sessions.Service,
	clientDisconnecter ClientDisconnecter,
) *Handler {
	return &Handler{
		service:            service,
		sessionService:     sessionService,
		clientDisconnecter: clientDisconnecter,
	}
}

// SocialLogin godoc
// @Summary Sign in with Apple
// @Description Authenticate a user using Apple Sign In credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SocialLoginRequest true "Apple sign-in details"
// @Success 200 {object} SocialLoginResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Login failed"
// @Router /auth/socialLogin [post]
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

// Refresh godoc
// @Summary Refresh authentication tokens
// @Description Generate new access and refresh tokens using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Router /auth/refresh [post]
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

// Logout godoc
// @Summary Log out
// @Description Invalidate the current session and disconnect the specified client
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Param X-Client-ID header string true "Client ID"
// @Success 200 {string} string "logout successful"
// @Failure 400 {object} map[string]interface{} "Client ID is required"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 500 {object} map[string]interface{} "Logout failed"
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	clientID := r.Header.Get("X-Client-ID")
	sessionID, _ := middleware.SessionID(r.Context())

	if err := h.sessionService.Logout(r.Context(), sessionID); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := h.clientDisconnecter.DisconnectClient(userID, clientID); err != nil {
		log.Printf("Logout websocket disconnect error: %v", err)
	}

	response.Write(w, http.StatusOK, "logout successful")
}
