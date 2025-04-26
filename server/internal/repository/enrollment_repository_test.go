package repository

import (
	"context"
	"testing"
	"time"

	"eduhub/server/internal/models"
	"eduhub/server/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnrollment(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(mocks.EnrollmentRepository)
	enrollment := &models.Enrollment{
		ID:             1,
		StudentID:      1,
		CourseID:       1,
		EnrollmentDate: time.Now(),
		Status:         "active",
	}

	mockRepo.On("CreateEnrollment", ctx, enrollment).Return(nil)
	err := mockRepo.CreateEnrollment(ctx, enrollment)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestIsStudentEnrolled(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.EnrollmentRepository)
	mockRepo.On("IsStudentEnrolled", ctx,1,2,3).Return(true, nil)
	ok, err := mockRepo.IsStudentEnrolled(ctx, 1, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, ok, true)
	mockRepo.AssertExpectations(t)

}
