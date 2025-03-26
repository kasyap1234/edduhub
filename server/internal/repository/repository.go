package repository

import (
	"eduhub/server/internal/models"
	"github.com/uptrace/bun"
)

type Repository struct {
	AttendanceRepository AttendanceRepository
	StudentRepository    StudentRepository
	UserRepository       UserRepository
	DatabaseRepository   DatabaseRepository[any]
}

// NewRepository creates a new repository with all required sub-repositories
// It needs a bun.DB instance to create the base repositories
func NewRepository(db *bun.DB) *Repository {
	// Create type-specific database repositories
	attendanceDB := NewBaseRepository[models.Attendance](db)
	studentDB := NewBaseRepository[models.Student](db)
	userDB := NewBaseRepository[models.User](db)
	
	// Create the specific repositories using the typed database repositories
	attendanceRepo := NewAttendanceRepository(attendanceDB)
	studentRepo := NewStudentRepository(studentDB)
	userRepo := NewUserRepository(userDB)
	
	// Create a generic database repository for any other needs
	genericDB := NewBaseRepository[any](db)
	
	return &Repository{
		AttendanceRepository: attendanceRepo,
		StudentRepository:    studentRepo,
		UserRepository:       userRepo,
		DatabaseRepository:   genericDB,
	}
}