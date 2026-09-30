package idempotency

import "time"

type IdempotencyRecord struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Status         string    `json:"status"`
	StatusCode     *int      `json:"status_code"`
	ResponseBody   []byte    `json:"response_body"`
	CreatedAt      time.Time `json:"created_at"`
}
