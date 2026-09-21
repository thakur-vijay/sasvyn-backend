package skills

import "time"

type Skill struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Skill     string    `json:"skill"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSkillDTO struct {
	Skill    string `json:"skill"`
	Category string `json:"category"`
}

type UpdateSkillDTO struct {
	Skill    string `json:"skill"`
	Category string `json:"category"`
}
