package repository

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindByKratosID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.StudentRepository)
	kratosID := "123"
	expectedStudent := &models.Student{
		StudentID:        1,
		KratosIdentityID: kratosID,
		CollegeID:        3,
		RollNo:           "SE20UARI73",
		Batch:            2020,
		Year:             5,
		Sem:              2,
		IsActive:         true,
	}

	// Mock the repository method
	mockRepo.On("FindByKratosID", ctx, kratosID).Return(expectedStudent, nil)

	// Act
	result, err := mockRepo.FindByKratosID(ctx, kratosID)

	// Assert
	assert.NoError(t, err, "Expected no error when finding student by Kratos ID")
	assert.Equal(t, expectedStudent, result, "Expected the returned student to match the mock")
	mockRepo.AssertExpectations(t)
}

func TestCreateStudent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.StudentRepository)
	newStudent := &models.Student{
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

	// Mock the repository method
	mockRepo.On("CreateStudent", ctx, newStudent).Return(nil)

	// Act
	err := mockRepo.CreateStudent(ctx, newStudent)

	// Assert
	assert.NoError(t, err, "Expected no error when creating a new student")
	mockRepo.AssertExpectations(t)
}

// Add more unit tests for other methods
