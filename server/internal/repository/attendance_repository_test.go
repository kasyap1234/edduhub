package repository

import (
	"context"
	"testing"

	"eduhub/server/internal/models"
	"eduhub/server/mocks"

	"github.com/stretchr/testify/assert"
)

func TestMarkAttendance(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.AttendanceRepository)
	mockRepo.On("MarkAttendance", ctx, 1, 2, 3, 4).Return(true, nil)
	result, err := mockRepo.MarkAttendance(ctx, 1, 2, 3, 4)
	assert.NoError(t, err)
	assert.True(t, result)
	mockRepo.AssertExpectations(t)
}

// TestUpdateAttendance tests the UpdateAttendance method of the AttendanceRepository
// by verifying that the method can successfully update an attendance record
// with the given parameters without returning an error.
func TestUpdateAttendance(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.AttendanceRepository)
	mockRepo.On("UpdateAttendance", ctx, 1, 2, 3, 4, "marked").Return(nil)
	err := mockRepo.UpdateAttendance(ctx, 1, 2, 3, 4, "marked")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAttendanceStudentInCourse(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.AttendanceRepository)
	expected := []*models.Attendance{}
	mockRepo.On("GetAttendanceStudentInCourse", ctx, 1, 2, 3).Return(expected, nil)
	result, err := mockRepo.GetAttendanceStudentInCourse(ctx, 1, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAttendanceStudent(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.AttendanceRepository)
	expected := []*models.Attendance{}
	mockRepo.On("GetAttendanceStudent", ctx, 1, 2).Return(expected, nil)
	result, err := mockRepo.GetAttendanceStudent(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAttendanceByLecture(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.AttendanceRepository)
	expected := []*models.Attendance{}
	mockRepo.On("GetAttendanceByLecture", ctx, 1, 2, 3).Return(expected, nil)
	result, err := mockRepo.GetAttendanceByLecture(ctx, 1, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAttendanceByCourse(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.AttendanceRepository)
	expected := []*models.Attendance{}
	mockRepo.On("GetAttendanceByCourse", ctx, 1, 2).Return(expected, nil)
	result, err := mockRepo.GetAttendanceByCourse(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
