package models

import "time"

type Department struct {
	ID        int       `db:"id" json:"id"`
	CollegeID int       `db:"college_id" json:"college_id"` // Foreign key to colleges table
	Name      string    `db:"name" json:"name"`
	HOD       string    `db:"hod" json:"hod"` // Head of Department (could be a user ID in the future)
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	College *College `db:"-" json:"college,omitempty"`
}
