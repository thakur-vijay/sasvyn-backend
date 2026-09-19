package auth

import "github.com/sasvyn/backend/internal/modules/users"

type SocialLoginRequest struct {
	AppleID  string `json:"apple_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type SocialLoginResponse struct {
	User         users.User `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
