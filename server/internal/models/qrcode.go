package models

import "time"

type QRCode struct {
	ID        int       `db:"id" json:"id"`
	StudentID int       `db:"student_id" json:"student_id"`
	QRCodeID  string    `db:"qr_code_id" json:"qr_code_id"`
	IssuedAt  time.Time `db:"issued_at" json:"issued_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	Student *Student `db:"-" json:"student,omitempty"`
}
