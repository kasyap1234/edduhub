package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type AttendanceRepository interface {
	MarkAttendance(ctx context.Context, studentID int, courseID int, lectureID int) error
	UpdateAttendance(ctx context.Context, studentID int, courseID int, lectureID int,status string) error
	GetAttendanceStudentInCourse(ctx context.Context, studentID int, courseID int)([]*models.Attendance,error)
	GetAttendanceStudent(ctx context.Context, studentID int) error
}

type attendanceRepository struct {
	dbRepo DatabaseRepository[models.Attendance]
}

func NewAttendanceRepository(db DatabaseRepository[models.Attendance]) AttendanceRepository {
	return &attendanceRepository{
		dbRepo: db,
	}
}

func (a *attendanceRepository) MarkAttendance(ctx context.Context, studentID int, courseID int, lectureID int) error {

	attendance := &models.Attendance{
		StudentID: studentID,
		CourseID:  courseID,
		LectureID: lectureID,
	}
	return a.dbRepo.Create(ctx, attendance)
}

func (a *attendanceRepository) UpdateAttendance(ctx context.Context, studentID int, courseID int, lectureID int,status string) error {
attendance,err :=a.dbRepo.FindOne(ctx,"student_id = ? AND course_id=? AND lecture_id=?",studentID,courseID,lectureID) 
if err !=nil {
	return err
}
attendance.Status=status
return a.dbRepo.Update(ctx,attendance)

}
func (a *attendanceRepository) GetAttendanceStudentInCourse(ctx context.Context, studentID int, courseID int) ([]*models.Attendance,error) {

 attendances,err :=a.dbRepo.FindWhere(ctx,"student_id=? AND course_id=?",studentID,courseID)
 if err !=nil {
	return nil,err; 
 }
 return attendances,nil 
}

func (a *attendanceRepository) GetAttendanceStudent(ctx context.Context, studentID int) error {
 var records []*models.Attendance
//  records,err := a.dbRepo.FindWhere(ctx,"student_id = ?",studentID)
records,err :=a.dbRepo.FindWhere()

}
