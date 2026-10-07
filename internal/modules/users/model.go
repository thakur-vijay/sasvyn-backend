package users

import "time"

type User struct {
	ID          string    `json:"id"`
	AppleID     string    `json:"apple_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	DateOfBirth *string   `json:"date_of_birth"`
	ImageUrl    *string   `json:"img_url"`
	ImageKey    *string   `json:"-"`
	SyncVersion *int64    `json:"sync_version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUserDTO struct {
	FullName    *string `json:"full_name" validate:"omitempty,notblank"`
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty,notblank"`
	ImageKey    *string `json:"img_key" validate:"omitempty,notblank"`
}
