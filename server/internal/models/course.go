package models

import "time"

// Course represents an individual course/subject
type Course struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Credits     int       `db:"credits" json:"credits"`
	InstructorID int      `db:"instructor_id" json:"instructor_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	
	// Relations - not stored in DB
	Instructor  *User     `db:"-" json:"instructor,omitempty"`
	Enrollments []*Enrollment `db:"-" json:"enrollments,omitempty"`
}

// Lecture represents an individual class session

type Lecture struct {
	ID        int `json:"id"`
	CollegeID int `json:"college_id"`
	CourseID  int `json:"course_id"`
	QRCodeID  int `json:"qr_code_id"`
}

// QRCode represents a unique QR code for each lecture

// Courses represents a collection of courses
type Courses struct {
	Items []*Course // Fixed syntax
}

// Subjects represents the courses a student is enrolled in
type Subjects struct {
	Current  Courses `json:"current"`  // Currently enrolled courses
	Previous Courses `json:"previous"` // Previously completed courses
	Optional Courses `json:"optional"` // Optional/elective courses
}
