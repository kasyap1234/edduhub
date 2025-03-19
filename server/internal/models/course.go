// Course represents an individual course/subject
type Course struct {
	CourseID        int `json:"CourseId" bun:"course_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Credits     int    `json:"credits"`
	Description string `json:"description,omitempty"`
	Department  string `json:"department,omitempty"`
	Instructor  string `json:"instructor,omitempty"`
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
