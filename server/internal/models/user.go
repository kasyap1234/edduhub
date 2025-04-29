// filepath: /home/tgt/Desktop/edduhub/server/internal/models/user.go
package models

import "time"

type User struct {
	ID               int       `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Role             string    `db:"role" json:"role"`
	Email            string    `db:"email" json:"email"`
	KratosIdentityID string    `db:"kratos_identity_id" json:"kratos_identity_id"`
	IsActive         bool      `db:"is_active" json:"is_active"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	Student *Student `db:"-" json:"student,omitempty"`
}
