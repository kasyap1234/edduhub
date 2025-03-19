package models 

type Attendance struct {
	ID int `json:"id" bun:"id,autoincrement"`
	StudentID uint `json:"studentId" bun:"student_id"`
	CourseID uint `json:"courseId" bun:"course_id"`
	Date time.Time `json:"date" bun:"date"`
	Status string `json:"status" bun:"status"` 
}

