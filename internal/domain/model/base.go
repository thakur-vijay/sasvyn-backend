package model

import "time"

type BaseModel struct {
	ID          string    `json:"id"`
	SyncVersion int64     `json:"sync_version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
