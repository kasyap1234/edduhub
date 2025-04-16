package models

type User struct {
	ID int `json:"ID" bun:"id,autoincrement"`

	Name     string   `json:"name" bun:"name"`
	Role     string   `json:"role" bun:"role"`
	Student  *Student `bun:"rel:has-one,join:id=user_id,nullzero"`
	IsActive bool     `bun:"is_active"`
}
