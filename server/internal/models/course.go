// Course represents an individual course/subject
package models 

type Course struct {
	ID        int `json:"id" bun:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Credits     int    `json:"credits"`
	Description string `json:"description,omitempty"`
	Department  string `json:"department,omitempty"`
	Instructor  string `json:"instructor,omitempty"`
}

type Lecture struct {
	bun.BaseModel `bun:"table:lectures"`
	ID int `bun:"pk,autoincrement"`
	CourseID int `bun:"not_null"`
	Course *Course `rel:"belongs-to,join:course_id=id"`
	QRCode *QRCode `json:"qrcode" bun:"qr_code"`

}

// Courses represents a collection of courses
type Courses struct {
	Items []Course `json:"items"`

}

// Subjects represents the courses a student is enrolled in
type Subjects struct {
	Current  Courses `json:"current"`  // Currently enrolled courses
	Previous Courses `json:"previous"` // Previously completed courses
	Optional Courses `json:"optional"` // Optional/elective courses
}
