package repository

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/mocks"
	"testing"
	"time"

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
