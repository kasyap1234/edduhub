package attendance

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
	"fmt"
)

const (
	Present = "Present"
	Absent  = "Absent"
	Freezed = "Freezed"
)

type AttendanceService interface {
	GenerateQRCode(ctx context.Context, collegeID, courseID int, lectureID int) (string, error)
	GetAttendanceByLecture(ctx context.Context, collegeID, courseID int, lectureID int) ([]*models.Attendance, error)
	GetAttendanceByCourse(ctx context.Context, collegeID, courseID int) ([]*models.Attendance, error)
	GetAttendanceByStudent(ctx context.Context, collegeID, studentID int) ([]*models.Attendance, error)
	GetAttendanceByStudentAndCourse(ctx context.Context, collegeID, studentID int, courseID int) ([]*models.Attendance, error)
	MarkAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int) (bool, error)
	UpdateAttendance(ctx context.Context, collegeID, studentID int, courseID int, lectureID int, currentStatus, updatedStatus string) (bool, error)
	FreezeAttendance(ctx context.Context, collegeID, studentID int) (bool, error)
	VerifyStudentEnrollment(ctx context.Context, collegeID, studentID, courseID int) (bool, error)
	FreezeStudent(ctx context.Context, collegeID int, RollNo string) error
	UnFreezeStudent(ctx context.Context, collegeID int, RollNo string) error
}

type attendanceService struct {
	repo           repository.AttendanceRepository
	studentRepo    repository.StudentRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewAttendanceService(repo repository.AttendanceRepository, studentRepo repository.StudentRepository, enrollmentRepo repository.EnrollmentRepository) AttendanceService {
	return &attendanceService{
		repo:           repo,
		studentRepo:    studentRepo,
		enrollmentRepo: enrollmentRepo,
	}

}

func (a *attendanceService) GetAttendanceByLecture(ctx context.Context, collegeID int, courseID int, lectureID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceByLecture(ctx, collegeID, courseID, lectureID)
}

// to get attendance of all students in a course
func (a *attendanceService) GetAttendanceByCourse(ctx context.Context, collegeID int, courseID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceByCourse(ctx, collegeID, courseID)
}

func (a *attendanceService) GetAttendanceByStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceStudent(ctx, collegeID, studentID)
}

func (a *attendanceService) GetAttendanceByStudentAndCourse(ctx context.Context, collegeID int, studentID int, courseID int) ([]*models.Attendance, error) {
	return a.repo.GetAttendanceStudentInCourse(ctx, collegeID, studentID, courseID)
}

func (a *attendanceService) MarkAttendance(ctx context.Context, collegeID int, studentID, courseID, lectureID int) (bool, error) {
	// First verify student enrollment and active status
	enrolled, err := a.VerifyStudentEnrollment(ctx, collegeID, studentID, courseID)
	if err != nil {
		return false, fmt.Errorf("failed to verify enrollment: %w", err)
	}
	if !enrolled {
		return false, fmt.Errorf("student not enrolled in course")
	}

	// Verify student is active
	student, err := a.studentRepo.GetStudentByID(ctx, collegeID, studentID)
	if err != nil {
		return false, fmt.Errorf("failed to get student: %w", err)
	}
	if !student.IsActive {
		return false, fmt.Errorf("student account is not active")
	}

	// Mark attendance only if all verifications pass
	ok, err := a.repo.MarkAttendance(ctx, collegeID, studentID, courseID, lectureID)
	if err != nil {
		return false, fmt.Errorf("failed to mark attendance: %w", err)
	}

	return ok, nil
}

func (a *attendanceService) VerifyStudentEnrollment(ctx context.Context, collegeID int, studentID int, courseID int) (bool, error) {
	student, err := a.studentRepo.GetStudentByID(ctx, collegeID, studentID)
	if err != nil {
		return false, err
	}
	if !student.IsActive {
		return false, nil
	}
	if student.CollegeID != collegeID {
		return false, nil
	}
	// enrolled, err := a.studentRepo.VerifyStudentEnrollment(ctx, collegeID, studentID, courseID)
	exists, err := a.enrollmentRepo.IsStudentEnrolled(ctx, studentID, courseID)
	return exists, err
}

func (a *attendanceService) GenerateAndProcessQRCode(ctx context.Context, collegeID, studentID int, courseID int, lectureID int) error {
	qrCode, err := a.GenerateQRCode(ctx, collegeID, courseID, lectureID)
	if err != nil {
		return nil
	}
	err = a.ProcessQRCode(ctx, collegeID, studentID, qrCode)
	if err != nil {
		return err
	}
	return nil
}

func (a *attendanceService) UpdateAttendance(ctx context.Context, collegeID, studentID int, courseID int, lectureID int, currentStatus, updatedStatus string) (bool, error) {
	switch currentStatus {
	case Present:
		updatedStatus = Absent

	case Absent:
		updatedStatus = Present

	}
	err := a.repo.UpdateAttendance(ctx, collegeID, studentID, courseID, lectureID, updatedStatus)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (a *attendanceService) FreezeAttendance(ctx context.Context, collegeID, studentID int) (bool, error) {
	err := a.repo.FreezeAttendance(ctx, collegeID, studentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *attendanceService) FreezeStudent(ctx context.Context, collegeID int, RollNo string) error {
	err := a.studentRepo.FreezeStudent(ctx, RollNo)
	return err
}

func (a *attendanceService) UnFreezeStudent(ctx context.Context, collegeID int, RollNo string) error {
	err := a.studentRepo.UnFreezeStudent(ctx, RollNo)
	return err
}
