package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type StudentRepository interface {
	GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error)
	GetStudentByID(ctx context.Context, studentID int) (*models.Student, error)
	UpdateStudent(ctx context.Context, model *models.Student) error
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

func (s *studentRepository) GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error) {
	student, err := s.db.FindOne(ctx, "rollNo = ?", rollNo)
	if err != nil {
		return nil, err
	}
	return student, nil

}

func (s *studentRepository) GetStudentByID(ctx context.Context, studentID int) (*models.Student, error) {
	student, err := s.db.FindOne(ctx, "student_id=?", studentID)
	return student, err
}

func (s *studentRepository) UpdateStudent(ctx context.Context, model *models.Student) error {
	err := s.db.Update(ctx, model)
	return err

}

func (s *studentRepository) FreezeStudent(ctx context.Context, RollNo string) error {
	student, err := s.GetStudentByRollNo(ctx, RollNo)
	if err != nil {
		student.IsActive = false
		return s.db.Update(ctx, student)
	}
	return err
}

func (s *studentRepository) UnFreezeStudent(ctx context.Context, RollNo string) error {
	student, err := s.GetStudentByRollNo(ctx, RollNo)
	if err != nil {
		student.IsActive = true
		return s.db.Update(ctx, student)
	}
	return err

}
