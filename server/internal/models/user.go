package models 

// Course represents an individual course/subject
type Course struct {
	ID          string `json:"id"`
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

type User struct {
	ID        int      `json:"ID" bun:"id,autoincrement"`
	CollegeID int      `json:"CollegeID" bun:"college_id"`
	RollNo    string   `json:"RollNo" bun:"roll_no,pk"`
	Batch     int      `json:"Batch" bun:"batch"`
	Year      int      `json:"Year" bun:"year"`
	Sem       int      `json:"Sem" bun:"sem"`
	Subjects  Subjects `json:"Subjects" bun:"subjects,json"`
}

