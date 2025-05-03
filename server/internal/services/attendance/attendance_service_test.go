package attendance

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupAttendanceServiceTest(t *testing.T) (*mocks.AttendanceRepository, *mocks.StudentRepository, *mocks.EnrollmentRepository, AttendanceService, context.Context) {
	mockAttendanceRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockEnrollmentRepo := new(mocks.EnrollmentRepository)
	service := NewAttendanceService(mockAttendanceRepo, mockStudentRepo, mockEnrollmentRepo)
	ctx := context.Background()
	return mockAttendanceRepo, mockStudentRepo, mockEnrollmentRepo, service, ctx
}

func TestGetAtendanceByLecture_Success(t *testing.T) {
	mockAttendanceRepo, _, _, service, ctx := setupAttendanceServiceTest(t)
	collegeID := 1
	courseID := 101
	lectureID := 201
	expectedAttendances := []*models.Attendance{
		{ID: 1, StudentID: 1, CourseID: courseID, LectureID: lectureID, CollegeID: collegeID, Status: "Present"},
		{ID: 2, StudentID: 2, CourseID: courseID, LectureID: lectureID, CollegeID: collegeID, Status: "Absent"},
	}
	mockAttendanceRepo.On("GetAttendanceByLecture", ctx, collegeID, courseID, lectureID).Return(expectedAttendances, nil).Once()
	attendances, err := service.GetAttendanceByLecture(ctx, collegeID, courseID, lectureID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttendances, attendances)
	mockAttendanceRepo.AssertExpectations(t)

}

func TestGetAttendanceByCourse_Success(t *testing.T) {
	{
		mockAttendanceRepo, _, _, service, ctx := setupAttendanceServiceTest(t)
		collegeID := 1
		courseID := 1
		expectedAttendance := []*models.Attendance{
			{ID: 1, StudentID: 1, CourseID: courseID, CollegeID: collegeID, Status: "Present"},
			{ID: 3, StudentID: 3, CourseID: courseID, CollegeID: collegeID, Status: "Present"},
		}

		mockAttendanceRepo.On("GetAttendanceByCourse", ctx, collegeID, courseID).Return(expectedAttendance, nil).Once()
		attendances, err := service.GetAttendanceByCourse(ctx, collegeID, courseID)
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance, attendances)
		mockAttendanceRepo.AssertExpectations(t)

	}
}
