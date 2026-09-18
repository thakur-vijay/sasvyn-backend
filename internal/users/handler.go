package users

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/google/uuid"
	"github.com/sasvyn/backend/internal/auth"
	"github.com/sasvyn/backend/internal/sessions"

	"database/sql"
	"log"
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, err := h.repository.GetByAppleID(r.Context(), request.AppleID)
	if err == nil {
		accessToken, refreshToken, err := h.sessionService.CreateSession(r.Context(), existingUser.ID)
		if err != nil {
			log.Printf("Create Session error: %v", err)
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		response := SocialLoginResponse{
			User:         *existingUser,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if err != sql.ErrNoRows {

		log.Printf("GetByAppleID error: %v", err)
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return

	}

	now := time.Now().UTC()

	user := User{
		ID:        uuid.NewString(),
		AppleID:   request.AppleID,
		FullName:  request.FullName,
		Email:     request.Email,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	if err := h.repository.Create(r.Context(), user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.sessionService.CreateSession(r.Context(), user.ID)
	if err != nil {
		log.Printf("Create Session error: %v", err)
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	response := SocialLoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.repository.GetByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		log.Printf("GetByID error: %v", err)
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
