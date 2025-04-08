package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type StudentRepository interface {
	GetStudentByRollNo(rollNo string) (*models.Student, error)
	GetStudentByID(studentID int) (*models.Student, error)
	UpdateStudent(*models.Student) error
	FreezeStudent(ctx context.Context, RollNo string) error
	UnFreezeStudent(ctx context.Context, RollNo string) error
}

type studentRepository struct {
	db DatabaseRepository[models.Student]
}

func NewStudentRepository(db DatabaseRepository[models.Student]) StudentRepository {
	return &studentRepository{
		db: db,
	}
}

func (s *studentRepository) GetStudentByRollNo(rollNo string) (*models.Student, error) {

}

func (s *studentRepository) GetStudentByID(studentID int) (*models.Student, error) {

}

func (s *studentRepository) UpdateStudent(*models.Student) error {

}

func (s *studentRepository) FreezeStudent(ctx context.Context, RollNo string) error {

}

func (s *studentRepository) UnFreezeStudent(ctx context.Context, RollNo string) error {

}
