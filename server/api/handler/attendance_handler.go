package handler

import (
	"eduhub/server/internal/services/attendance"
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
