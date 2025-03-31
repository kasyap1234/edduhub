package attendance

import (
	"eduhub/server/internal/models"

	"github.com/uptrace/bun"
)

type AttendanceService interface {
	GenerateQRCode(courseID int, lectureID int) (string, error)
	GetAttendanceByLecture(courseID int, lectureID int) (*models.Attendance, error)
	GetAttendanceByCourse(courseID int) (*models.Attendance, error)
	GetAttendanceByStudent(studentID int) (*models.Attendance, error)
	GetAttendanceByStudentAndCourse(studentID int, courseID int) (*models.Attendance, error)
	MarkAttendance(studentID int, courseID int, lectureID int) (bool, error)
}

type attendanceService struct {
	db *bun.DB
}

func NewAttendanceService(db *bun.DB) AttendanceService {
	return &attendanceService{
		db: db,
	}
}

func (a *attendanceService) GetAttendanceByLecture(courseID int, lectureID int) (*models.Attendance, error) {

}

func (a *attendanceService) GetAttendanceByCourse(courseID int) (*models.Attendance, error) {

}

func (a *attendanceService) GetAttendanceByStudent(studentID int) (*models.Attendance, error) {

}

func (a *attendanceService) GetAttendanceByStudentAndCourse(studentID int, courseID int) (*models.Attendance, error) {

}

func (a *attendanceService) MarkAttendance(studentID, courseID, lectureID int) (bool, error) {

}
