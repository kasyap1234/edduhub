package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type AttendanceRepository interface {
	GetAttendanceByCourse(ctx context.Context, courseID int) ([]*models.Attendance, error)
	MarkAttendance(ctx context.Context, studentID int, courseID int, lectureID int) (bool, error)
	UpdateAttendance(ctx context.Context, studentID int, courseID int, lectureID int, status string) error
	GetAttendanceStudentInCourse(ctx context.Context, studentID int, courseID int) ([]*models.Attendance, error)
	GetAttendanceStudent(ctx context.Context, studentID int) ([]*models.Attendance, error)
	GetAttendanceByLecture(ctx context.Context, lectureID int, courseID int) ([]*models.Attendance, error)
}

type attendanceRepository struct {
	db DatabaseRepository[models.Attendance]
}

func NewAttendanceRepository(db DatabaseRepository[models.Attendance]) AttendanceRepository {
	return &attendanceRepository{
		db: db,
	}
}

// mark attendance

func (a *attendanceRepository) MarkAttendance(ctx context.Context, studentID int, courseID int, lectureID int) (bool, error) {

	attendance := &models.Attendance{
		StudentID: studentID,
		CourseID:  courseID,
		LectureID: lectureID,
	}
	err := a.db.Create(ctx, attendance)
	if err != nil {
		return false, err
	}
	return true, nil
}

// to update student attendance for a single lecture in a course
func (a *attendanceRepository) UpdateAttendance(ctx context.Context, studentID int, courseID int, lectureID int, status string) error {
	attendance, err := a.db.FindOne(ctx, "student_id = ? AND course_id=? AND lecture_id=?", studentID, courseID, lectureID)
	if err != nil {
		return err
	}
	attendance.Status = status
	return a.db.Update(ctx, attendance)

}

// to get attendance of student in a course a course is a series of lectures wil give student attendance for the entire course for example student 1 has attendance for lecture 1 and lecture2 of course 1

func (a *attendanceRepository) GetAttendanceStudentInCourse(ctx context.Context, studentID int, courseID int) ([]*models.Attendance, error) {

	attendances, err := a.db.FindWhere(ctx, "student_id=? AND course_id=?", studentID, courseID)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// to get attendance of student across all courses
// to get attendance of student across all lectures

func (a *attendanceRepository) GetAttendanceStudent(ctx context.Context, studentID int) ([]*models.Attendance, error) {
	records, err := a.db.FindWhere(ctx, "student_id = ?", studentID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (a *attendanceRepository) GetAttendanceByLecture(ctx context.Context, courseID int, lectureID int) ([]*models.Attendance, error) {
	records, err := a.db.FindWhere(ctx, "course_id=? AND lecture_id=?", courseID, lectureID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (a *attendanceRepository) GetAttendanceByCourse(ctx context.Context, courseID int) ([]*models.Attendance, error) {
	records, err := a.db.FindWhere(ctx, "course_id=?", courseID)
	if err != nil {
		return nil, err
	}
	return records, nil
}
