package repository

type Repository struct {
	AttendanceRepository AttendanceRepository
	StudentRepository    StudentRepository
	UserRepository       UserRepository
	EnrollmentRepository EnrollmentRepository
	PlacementRepository  PlacementRepository  // Added Placement
	QuizRepository       QuizRepository       // Added Quiz
	DepartmentRepository DepartmentRepository // Added Department
}

// NewRepository creates a new repository with all required sub-repositories
// It needs a bun.DB instance to create the base repositories
func NewRepository(DB *DB) *Repository {
	// Create type-specific database repositories
	attendanceRepo := NewAttendanceRepository(DB)
	studentRepo := NewStudentRepository(DB)
	userRepo := NewUserRepository(DB)
	enrollmentRepo := NewEnrollmentRepository(DB)
	placementRepo := NewPlacementRepository(DB)   // Instantiate Placement
	quizRepo := NewQuizRepository(DB)             // Instantiate Quiz
	departmentRepo := NewDepartmentRepository(DB) // Instantiate Department
	return &Repository{
		AttendanceRepository: attendanceRepo,
		StudentRepository:    studentRepo,
		UserRepository:       userRepo,
		EnrollmentRepository: enrollmentRepo,
		PlacementRepository:  placementRepo,
		QuizRepository:       quizRepo,
		DepartmentRepository: departmentRepo,
	}
}
