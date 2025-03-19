package models 

import "github.com/uptrace/bun"
type Attendance struct {
	bun.BaseModel `bun:"table:attendance"`
	ID int `json:"id" bun:"pk,id,autoincrement"`
	StudentID uint `json:"studentId" bun:"student_id"`
	CourseID uint `json:"courseId" bun:"course_id"`
	Student *Student `bun:"rel:belongs-to,join:student_id=id"`
	Course *Course    `bun:"rel:belongs-to,join:course_id=id"`
	LectureID int 	`bun:"notnull"`
	Lecture *Lecture  `bun:"rel:belongs-to,join:lecture_id=id"`
	Date time.Time `json:"date" bun:"date"`
	Status string `json:"status" bun:"status"` 
	ScannedAt time.Time `json:"scannedAt" bun:"scanned_at"`
}


