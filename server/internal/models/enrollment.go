package models

import "time"

type Enrollment struct {
	ID             int       `db:"id" json:"id"`
	StudentID      int       `db:"student_id" json:"student_id"`
	CourseID       int       `db:"course_id" json:"course_id"`
	CollegeID      int       `db:"college_id" json:"course_id"`
	EnrollmentDate time.Time `db:"enrollment_date" json:"enrollment_date"`
	Status         string    `db:"status" json:"status"` // Active, Completed, Dropped
	Grade          string    `db:"grade" json:"grade,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	Student *Student `db:"-" json:"student,omitempty"`
	Course  *Course  `db:"-" json:"course,omitempty"`
}
