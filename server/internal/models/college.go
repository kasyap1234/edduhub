package models

import "time"

type College struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Address   string    `db:"address" json:"address"`
	City      string    `db:"city" json:"city"`
	State     string    `db:"state" json:"state"`
	Country   string    `db:"country" json:"country"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	Students []*Student `db:"-" json:"students,omitempty"`
}
