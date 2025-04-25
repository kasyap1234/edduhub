package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Enrollment struct {
	bun.BaseModel  `bun:"table:enrollments,alias:e"` // Added alias 'e'
	ID             int                               `json:"id" bun:",pk,autoincrement"`
	StudentID      int                               `json:"student_id" bun:"student_id,notnull"` // Foreign Key to students
	CourseID       int                               `json:"course_id" bun:"course_id,notnull"`   // Foreign Key to courses
	EnrollmentDate time.Time                         `json:"enrollment_date" bun:",default:current_timestamp"`
	Status         string                            `json:"status" bun:"status,default:'active'"` // e.g., 'active', 'dropped', 'completed'

	// --- Relationships ---
	Student *Student `json:"student,omitempty" bun:"rel:belongs-to,join:student_id=student_id"`
	Course  *Course  `json:"course,omitempty" bun:"rel:belongs-to,join:course_id=id"`
}
