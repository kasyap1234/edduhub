package repository

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/services/attendance"
)

type AttendanceRepository interface {
	GetAttendanceByCourse(ctx context.Context, collegeID int, courseID int) ([]*models.Attendance, error)
	MarkAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int) (bool, error)
	UpdateAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int, status string) error
	GetAttendanceStudentInCourse(ctx context.Context, collegeID int, studentID int, courseID int) ([]*models.Attendance, error)
	GetAttendanceStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Attendance, error)
	GetAttendanceByLecture(ctx context.Context, collegeID int, lectureID int, courseID int) ([]*models.Attendance, error)
	FreezeAttendance(ctx context.Context, collegeID int, studentID int) error
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

func (a *attendanceRepository) MarkAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int) (bool, error) {

	attendance := &models.Attendance{
		CollegeID: collegeID,
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
func (a *attendanceRepository) UpdateAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int, status string) error {
	attendance, err := a.db.FindOne(ctx, "college_id=? AND student_id = ? AND course_id=? AND lecture_id=?", collegeID, studentID, courseID, lectureID)
	if err != nil {
		return err
	}
	attendance.Status = status
	return a.db.Update(ctx, attendance)

}

// to get attendance of student in a course a course is a series of lectures wil give student attendance for the entire course for example student 1 has attendance for lecture 1 and lecture2 of course 1

func (a *attendanceRepository) GetAttendanceStudentInCourse(ctx context.Context, collegeID int, studentID int, courseID int) ([]*models.Attendance, error) {

	attendances, err := a.db.FindWhere(ctx, "college_id=? AND student_id=? AND course_id=?", collegeID, studentID, courseID)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// to get attendance of student across all courses and lectures
func (a *attendanceRepository) GetAttendanceStudent(ctx context.Context, collegeID, studentID int) ([]*models.Attendance, error) {
	records, err := a.db.FindWhere(ctx, "student_id = ? AND college_id=?", studentID, collegeID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (a *attendanceRepository) GetAttendanceByLecture(ctx context.Context, collegeID, courseID int, lectureID int) ([]*models.Attendance, error) {
	records, err := a.db.FindWhere(ctx, "course_id=? AND lecture_id=? AND college_id=?", courseID, lectureID, collegeID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// attendance of all students in a course

func (a *attendanceRepository) GetAttendanceByCourse(ctx context.Context, collegeID, courseID int) ([]*models.Attendance, error) {
	records, err := a.db.FindWhere(ctx, "college_id=? AND course_id=?", collegeID, courseID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (a *attendanceRepository) FreezeAttendance(ctx context.Context, collegeID, studentID int) error {
	student, err := a.db.FindOne(ctx, "college_id=? AND student_id=?", collegeID, studentID)
	if err != nil {
		return err
	}
	student.Status = attendance.Freezed
	return a.db.Update(ctx, student)

}
