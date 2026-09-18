package users

type SocialLoginRequest struct {
	AppleID  string `json:"apple_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type User struct {
	ID        string `json:"id"`
	AppleID   string `json:"apple_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SocialLoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
