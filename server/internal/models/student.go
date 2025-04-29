package models

import "time"

type Student struct {
	ID               int       `db:"id" json:"id"`
	UserID           int       `db:"user_id" json:"user_id"`
	CollegeID        int       `db:"college_id" json:"college_id"`
	KratosIdentityID string    `db:"kratos_identity_id" json:"kratos_identity_id"`
	EnrollmentYear   int       `db:"enrollment_year" json:"enrollment_year"`
	RollNo           string    `db:"roll_no" json:"roll_no"`     // Added this field
	IsActive         bool      `db:"is_active" json:"is_active"` // Added this field
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB (add db:"-" tag)
	// College     *College      `db:"-" json:"college,omitempty"`
	// Enrollments []*Enrollment `db:"-" json:"enrollments,omitempty"`
	// QRCodes     []*QRCode     `db:"-" json:"qr_codes,omitempty"`
}
