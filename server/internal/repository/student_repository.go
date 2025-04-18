package repository

import (
	"context"
	"database/sql"
	"eduhub/server/internal/models"
	"errors"
)

type StudentRepository interface {
	CreateStudent(ctx context.Context, student *models.Student) error
	GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error)
	GetStudentByID(ctx context.Context, studentID int) (*models.Student, error)
	UpdateStudent(ctx context.Context, model *models.Student) error
	FreezeStudent(ctx context.Context, RollNo string) error
	UnFreezeStudent(ctx context.Context, RollNo string) error
	FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error)
}

type studentRepository struct {
	db DatabaseRepository[models.Student]
}

func NewStudentRepository(db DatabaseRepository[models.Student]) StudentRepository {
	return &studentRepository{
		db: db,
	}
}

func (s *studentRepository) CreateStudent(ctx context.Context, student *models.Student) error {
	return s.db.Create(ctx, student)
}

// Add FindByKratosID implementation
func (s *studentRepository) FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error) {
	// Ensure column name 'kratos_identity_id' matches your DB schema
	student, err := s.db.FindOne(ctx, "kratos_identity_id = ?", kratosID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil, nil for not found - let service/middleware handle logic
		}
		return nil, err // Return other DB errors
	}
	return student, nil
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
