package models

import "time"

// Course represents an individual course/subject
type Course struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	CollegeID    int       `db:"college_id" json:"college_id"`
	Description  string    `db:"description" json:"description"`
	Credits      int       `db:"credits" json:"credits"`
	InstructorID int       `db:"instructor_id" json:"instructor_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

	// Relations - not stored in DB
	Instructor  *User         `db:"-" json:"instructor,omitempty"`
	Enrollments []*Enrollment `db:"-" json:"enrollments,omitempty"`
}

// Lecture represents an individual class session

type Lecture struct {
	ID          int       `db:"id" json:"id"`                               // Primary Key
	CourseID    int       `db:"course_id" json:"course_id"`                 // Foreign key to courses table
	CollegeID   int       `db:"college_id" json:"college_id"`               // Denormalized, Foreign key to colleges table
	Title       string    `db:"title" json:"title"`                         // Title of the lecture
	Description string    `db:"description" json:"description,omitempty"`   // Optional description
	StartTime   time.Time `db:"start_time" json:"start_time"`               // Start time of the lecture
	EndTime     time.Time `db:"end_time" json:"end_time"`                   // End time of the lecture
	MeetingLink string    `db:"meeting_link" json:"meeting_link,omitempty"` // For online lectures
	CreatedAt   time.Time `db:"created_at" json:"created_at"`               // Timestamp of creation
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`               // Timestamp of last update
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
