package users

import "time"

type User struct {
	ID          string    `json:"id"`
	AppleID     string    `json:"apple_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	DateOfBirth *string   `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUserDTO struct {
	FullName    string  `json:"full_name"`
	DateOfBirth *string `json:"date_of_birth"`
}
