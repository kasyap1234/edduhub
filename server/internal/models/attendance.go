package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Attendance struct {
	bun.BaseModel `bun:"table:attendance"`
	ID            int       `json:"id" bun:"pk,autoincrement"`
	StudentID     int       `json:"studentId" bun:"student_id"`
	CourseID      int       `json:"courseId" bun:"course_id"`
	CollegeID     int       `json:"collegeID" bun:"college_id"`
	Date          time.Time `json:"date" bun:"date"`
	Status        string    `json:"status" bun:"status"`
	ScannedAt     time.Time `json:"scannedAt" bun:"scanned_at"`
	LectureID     int       `bun:"notnull"`

	Student *Student `bun:"rel:belongs-to,join:student_id=id"`
	Course  *Course  `bun:"rel:belongs-to,join:course_id=id"`
	Lecture *Lecture `bun:"rel:belongs-to,join:lecture_id=id"`
}
