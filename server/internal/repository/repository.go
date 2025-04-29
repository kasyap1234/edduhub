package repository

type Repository struct {
	AttendanceRepository AttendanceRepository
	StudentRepository    StudentRepository
	UserRepository       UserRepository
	EnrollmentRepository EnrollmentRepository
}

// NewRepository creates a new repository with all required sub-repositories
// It needs a bun.DB instance to create the base repositories
func NewRepository(DB *DB) *Repository {
	// Create type-specific database repositories
	attendanceRepo := NewAttendanceRepository(DB)
	studentRepo := NewStudentRepository(DB)
	userRepo := NewUserRepository(DB)
	enrollmentRepo := NewEnrollmentRepository(DB)
	return &Repository{
		AttendanceRepository: attendanceRepo,
		StudentRepository:    studentRepo,
		UserRepository:       userRepo,
		EnrollmentRepository: enrollmentRepo,
	}
}
