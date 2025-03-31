package attendance

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"

	"github.com/uptrace/bun"
)

type AttendanceService interface {
	GenerateQRCode(courseID int, lectureID int) (string, error)
	GetAttendanceByLecture(courseID int, lectureID int) ([]*models.Attendance, error)
	GetAttendanceByCourse(courseID int) (*models.Attendance, error)
	GetAttendanceByStudent(studentID int) (*models.Attendance, error)
	GetAttendanceByStudentAndCourse(studentID int, courseID int) (*models.Attendance, error)
	MarkAttendance(studentID int, courseID int, lectureID int) (bool, error)
}

type attendanceService struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {
	return &attendanceService{
	repo :repo ,
	}
}

func (a *attendanceService) GetAttendanceByLecture(courseID int, lectureID int) ([]*models.Attendance, error) {
return a.repo.GetAttendanceByLecture(context.Background(),courseID,lectureID)
}

func (a *attendanceService) GetAttendanceByCourse(courseID int) (*models.Attendance, error) {

}

func (a *attendanceService) GetAttendanceByStudent(studentID int) (*models.Attendance, error) {

}

func (a *attendanceService) GetAttendanceByStudentAndCourse(studentID int, courseID int) (*models.Attendance, error) {

}

func (a *attendanceService) MarkAttendance(studentID, courseID, lectureID int) (bool, error) {

}
