package users

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/google/uuid"
)

func SocialLogin(w http.ResponseWriter, r *http.Request) {
	var request SocialLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
