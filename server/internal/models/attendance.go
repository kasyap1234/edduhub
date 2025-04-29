package models

import (
	"time"
)

type Attendance struct {
	ID        int       `json:"ID"`
	StudentID int       `json:"studentID"`
	CourseID  int       `json:"courseId"`
	CollegeID int       `json:"collegeID"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	ScannedAt time.Time `json:"scannedAt"`
	LectureID int       `json:"lectureID"`
}
