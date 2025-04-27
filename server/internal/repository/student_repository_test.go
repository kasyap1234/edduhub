// package repository

// import (
// 	"context"
// 	"eduhub/server/internal/models"
// 	"eduhub/server/mocks"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestFindByKratosID(t *testing.T) {
// 	// Arrange
// 	ctx := context.Background()
// 	mockRepo := new(mocks.StudentRepository)
// 	kratosID := "123"
// 	expectedStudent := &models.Student{
// 		StudentID:        1,
// 		KratosIdentityID: kratosID,
// 		CollegeID:        3,
// 		RollNo:           "SE20UARI73",
// 		Batch:            2020,
// 		Year:             5,
// 		Sem:              2,
// 		IsActive:         true,
// 	}

// 	// Mock the repository method
// 	mockRepo.On("FindByKratosID", ctx, kratosID).Return(expectedStudent, nil)

// 	// Act
// 	result, err := mockRepo.FindByKratosID(ctx, kratosID)

// 	// Assert
// 	assert.NoError(t, err, "Expected no error when finding student by Kratos ID")
// 	assert.Equal(t, expectedStudent, result, "Expected the returned student to match the mock")
// 	mockRepo.AssertExpectations(t)
// }

// func TestCreateStudent(t *testing.T) {
// 	// Arrange
// 	ctx := context.Background()
// 	mockRepo := new(mocks.StudentRepository)
// 	newStudent := &models.Student{
// 		StudentID:        1,
// 		KratosIdentityID: "identity",
// 		CollegeID:        3,
// 		RollNo:           "SE20UARI73",
// 		Batch:            2020,
// 		Year:             5,
// 		Sem:              2,
// 		Subjects:         models.Subjects{Current: models.Courses{Items: []*models.Course{}}},
// 		IsActive:         true,
// 	}

// 	// Mock the repository method
// 	mockRepo.On("CreateStudent", ctx, newStudent).Return(nil)

// 	// Act
// 	err := mockRepo.CreateStudent(ctx, newStudent)

// 	// Assert
// 	assert.NoError(t, err, "Expected no error when creating a new student")
// 	mockRepo.AssertExpectations(t)
// }

// Add more unit tests for other methods like table driven test cases 

package repository

import (
    "context"
    "errors"
    "eduhub/server/internal/models"
    "eduhub/server/mocks"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestFindByKratosID(t *testing.T) {
    // Arrange
    ctx := context.Background()
    mockRepo := new(mocks.StudentRepository)

    // Define test cases
    testCases := []struct {
        name           string
        kratosID       string
        mockReturn     *models.Student
        mockError      error
        expectedError  bool
        expectedResult *models.Student
    }{
        {
            name:           "Success case",
            kratosID:       "123",
            mockReturn:     &models.Student{StudentID: 1, KratosIdentityID: "123"},
            mockError:      nil,
            expectedError:  false,
            expectedResult: &models.Student{StudentID: 1, KratosIdentityID: "123"},
        },
        {
            name:           "Failure case - Not found",
            kratosID:       "999",
            mockReturn:     nil,
            mockError:      errors.New("student not found"),
            expectedError:  true,
            expectedResult: nil,
        },
        {
            name:           "Failure case - Invalid input",
            kratosID:       "",
            mockReturn:     nil,
            mockError:      errors.New("invalid input"),
            expectedError:  true,
            expectedResult: nil,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Mock the repository method
            mockRepo.On("FindByKratosID", ctx, tc.kratosID).Return(tc.mockReturn, tc.mockError)

            // Act
            result, err := mockRepo.FindByKratosID(ctx, tc.kratosID)

            // Assert
            if tc.expectedError {
                assert.Error(t, err, "Expected an error but got none")
                assert.Equal(t, tc.mockError.Error(), err.Error(), "Expected error message does not match")
            } else {
                assert.NoError(t, err, "Expected no error but got one")
            }
            assert.Equal(t, tc.expectedResult, result, "Expected result does not match actual result")

            // Verify mock expectations
            mockRepo.AssertExpectations(t)
        })
    }
}


func TestCreateStudent(t *testing.T) {
    // Arrange
    ctx := context.Background()
    mockRepo := new(mocks.StudentRepository)

    // Define test cases
    testCases := []struct {
        name          string
        student       *models.Student
        mockError     error
        expectedError bool
    }{
        {
            name: "Success case",
            student: &models.Student{
                StudentID:        1,
                KratosIdentityID: "identity",
                CollegeID:        3,
                RollNo:           "SE20UARI73",
                Batch:            2020,
                Year:             5,
                Sem:              2,
                IsActive:         true,
            },
            mockError:     nil,
            expectedError: false,
        },
        {
            name: "Failure case - Duplicate entry",
            student: &models.Student{
                StudentID:        1,
                KratosIdentityID: "identity",
            },
            mockError:     errors.New("duplicate entry"),
            expectedError: true,
        },
        {
            name: "Failure case - Invalid input",
            student: &models.Student{
                StudentID: 0, // Invalid ID
            },
            mockError:     errors.New("invalid input"),
            expectedError: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Mock the repository method
            mockRepo.On("CreateStudent", ctx, tc.student).Return(tc.mockError)

            // Act
            err := mockRepo.CreateStudent(ctx, tc.student)

            // Assert
            if tc.expectedError {
                assert.Error(t, err, "Expected an error but got none")
                assert.Equal(t, tc.mockError.Error(), err.Error(), "Expected error message does not match")
            } else {
                assert.NoError(t, err, "Expected no error but got one")
            }

            // Verify mock expectations
            mockRepo.AssertExpectations(t)
        })
    }
}