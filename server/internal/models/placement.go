package models

import "time"

// Placement represents a student's job placement record.
type Placement struct {
	ID            int       `db:"id" json:"id"`
	CollegeID     int       `db:"college_id" json:"college_id"`
	StudentID     int       `db:"student_id" json:"student_id"`
	CompanyName   string    `db:"company_name" json:"company_name"`
	JobTitle      string    `db:"job_title" json:"job_title"`
	Package       float64   `db:"package" json:"package"` // Assuming salary package
	PlacementDate time.Time `db:"placement_date" json:"placement_date"`
	Status        string    `db:"status" json:"status"` // e.g., Offered, Accepted, Rejected, On-Hold
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	Student *Student `db:"-" json:"student,omitempty"`
}
