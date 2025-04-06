package handler

import (
	"eduhub/server/internal/models"
	"eduhub/server/internal/services/attendance"
	"errors"
)

type AttendanceHandler struct {
	attendanceService attendance.AttendanceService
}

func NewAttedanceHandler(attendance attendance.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendance,
	}
}

func (a *AttendanceHandler) MarkAttendance() {

}

func(a *AttendanceHandler)GetAttendanceForStudent(studentID int)([]*models.Attendance){
	attendance, err := a.attendanceService.GetAttendanceByStudent(studentID)
	if err !=nil{
		errors.New("aunable to fetch attendance of student")

	}
	return attendance
}

func(a *AttendanceHandler)GetAttendanceByStudentInCourse(studentID , courseID int)([]*models.Attendance,error){
	attendance,err :=a.attendanceService.GetAttendanceByStudentAndCourse(studentID, courseID)
	if err !=nil{
		return nil,err 
	}
	return attendance,nil 
}


func(a*AttendanceHandler)UpdateAttendance()

