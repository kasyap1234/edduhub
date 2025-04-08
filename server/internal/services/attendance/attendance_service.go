package attendance

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services/attendance"
)

const (
	Present = "Present"
	Absent  = "Absent"
	Freezed = "Freezed"
)

type AttendanceService interface {
	GenerateQRCode(courseID int, lectureID int) (string, error)
	GetAttendanceByLecture(courseID int, lectureID int) ([]*models.Attendance, error)
	GetAttendanceByCourse(courseID int) ([]*models.Attendance, error)
	GetAttendanceByStudent(studentID int) ([]*models.Attendance, error)
	GetAttendanceByStudentAndCourse(studentID int, courseID int) ([]*models.Attendance, error)
	MarkAttendance(ctx context.Context, studentID int, courseID int, lectureID int) (bool, error)
	UpdateAttendance(ctx context.Context, studentID int, courseID int, lectureID int, currentStatus, updatedStatus string) (bool, error)
	FreezeAttendance(ctx context.Context, studentID int) error
	FreezeStudent(ctx context.Context, RollNo string) error
	UnFreezeStudent(ctx context.Context, RollNo string) error
}

type attendanceService struct {
	repo        repository.AttendanceRepository
	studentRepo repository.StudentRepository
}

func NewAttendanceService(repo repository.AttendanceRepository, userRepo repository.UserRepository) AttendanceService {
	return &attendanceService{
		repo:        repo,
		studentRepo: repository.StudentRepository,
	}
}

func (a *attendanceService) GetAttendanceByLecture(courseID int, lectureID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceByLecture(context.Background(), courseID, lectureID)
}

// to get attendance of all students in a course
func (a *attendanceService) GetAttendanceByCourse(courseID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceByCourse(context.Background(), courseID)
}

func (a *attendanceService) GetAttendanceByStudent(studentID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceStudent(context.Background(), studentID)
}

func (a *attendanceService) GetAttendanceByStudentAndCourse(studentID int, courseID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceStudentInCourse(context.Background(), studentID, courseID)
}

func (a *attendanceService) MarkAttendance(ctx context.Context, studentID, courseID, lectureID int) (bool, error) {
	ok, err := a.repo.MarkAttendance(ctx, studentID, courseID, lectureID)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (a *attendanceService) VerifyStudentEnrollment(ctx context.Context, studentID int, courseID int) (bool, error) {
	enrolled, err := a.repo.VerifyStudentEnrollment(ctx, studentID, courseID)
	if err != nil {
		return false, err
	}
	return enrolled, nil
}

func (a *attendanceService) GenerateAndProcessQRCode(ctx context.Context, studentID int, courseID int, lectureID int) error {
	qrCode, err := a.GenerateQRCode(courseID, lectureID)
	if err != nil {
		return nil
	}
	err = a.ProcessQRCode(ctx, studentID, qrCode)
	if err != nil {
		return err
	}
	return nil
}

func (a *attendanceService) UpdateAttendance(ctx context.Context, studentID int, courseID int, lectureID int, currentStatus, updatedStatus string) (bool, error) {
	switch currentStatus {
	case attendance.Present:
		err := a.repo.UpdateAttendance(ctx, studentID, courseID, lectureID, attendance.Absent)
		if err != nil {
			return false, err
		}
	case attendance.Absent:
		err := a.repo.UpdateAttendance(ctx, studentID, courseID, lectureID, attendance.Present)
		if err != nil {
			return false, err
		}

	}
	return true, nil

}

func (a *attendanceService) FreezeAttendance(ctx context.Context, studentID int) (bool, error) {
	err := a.repo.FreezeAttendance(ctx, studentID)
	if err != nil {
		return false, err
	}
	return true, nil
}
