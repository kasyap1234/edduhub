package repository

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/mocks"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindByKratosID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.StudentRepository)
	kratosID := "123"
	student := &models.Student{}
	mockRepo.On("FindByKratosID", ctx, kratosID).Return(student, nil)
	result, err := mockRepo.FindByKratosID(ctx, kratosID)
	assert.NoError(t, err)
	assert.Equal(t, student, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateStudent(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.StudentRepository)
	student := &models.Student{
		StudentID:        1,
		KratosIdentityID: "identity",
		CollegeID:        3,
		RollNo:           "SE20UARI73",
		Batch:            2020,
		Year:             5,
		Sem:              2,
		Subjects:         models.Subjects{Current: models.Courses{Items: []*models.Course{}}},
		IsActive:         true,
	}
	mockRepo.On("CreateStudent", ctx, student).Return(nil)
	err := mockRepo.CreateStudent(ctx, student)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Add more unit tests for other methods
