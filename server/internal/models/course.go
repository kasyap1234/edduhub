package models

import "github.com/uptrace/bun"

// Course represents an individual course/subject
type Course struct {
	bun.BaseModel `bun:"table:courses"`
	ID            int        `json:"id" bun:",pk,autoincrement"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Credits       int        `json:"credits"`
	Description   string     `json:"description,omitempty"`
	Department    string     `json:"department,omitempty"`
	Instructor    string     `json:"instructor,omitempty"`
	Lectures      []*Lecture `bun:"rel:has-many,join:course_id=id"` // Fixed join syntax
}

// Lecture represents an individual class session

type Lecture struct {
	bun.BaseModel `bun:"table:lectures"`
	ID            int     `json:"id" bun:",pk,autoincrement"`
	CourseID      int     `json:"course_id" bun:",notnull"`
	QRCodeID      int     `json:"qr_code_id"`
	
	Course  *Course `bun:"rel:belongs-to,join:course_id=id"`
	QRCode  *QRCode `bun:"rel:belongs-to,join:qr_code_id=id"`
}

// QRCode represents a unique QR code for each lecture

// Courses represents a collection of courses
type Courses struct {
	Items []*Course `bun:"rel:has-many,join:id=course_id"` // Fixed syntax
}

// Subjects represents the courses a student is enrolled in
type Subjects struct {
	Current  Courses `json:"current"`  // Currently enrolled courses
	Previous Courses `json:"previous"` // Previously completed courses
	Optional Courses `json:"optional"` // Optional/elective courses
}
