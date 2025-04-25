package repository

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCourse(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.CourseRepository)
	course := &models.Course{
		ID:          0,
		CollegeID:   1,
		Name:        "abc",
		Code:        "111",
		Credits:     4,
		Description: "description",
		Department:  "chemistry",
		Instructor:  "virat",
		Lectures:    []*models.Lecture(nil),
	}
	mockRepo.On("CreateCourse", ctx, course).Return(nil)
	err := mockRepo.CreateCourse(ctx, course)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

