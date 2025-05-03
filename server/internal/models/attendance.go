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

// StudentAttendanceStatus is used for bulk attendance marking requests.
type StudentAttendanceStatus struct {
	StudentID int    `json:"student_id" validate:"required,gt=0"`
	Status    string `json:"status" validate:"required,oneof=Present Absent"` // Ensure status is either Present or Absent
}
