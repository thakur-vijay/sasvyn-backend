package users

type User struct {
	ID        string `json:"id"`
	AppleID   string `json:"apple_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
