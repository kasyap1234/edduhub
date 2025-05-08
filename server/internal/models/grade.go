// /home/tgt/Desktop/edduhub/server/internal/models/grade.go (Example, create or adjust as needed)
package models

import "time"

type Grade struct {
	ID            int       `db:"id" json:"id"`
	StudentID     string    `db:"student_id" json:"student_id"` // Kratos ID or internal student identifier
	CourseID      int       `db:"course_id" json:"course_id"`
	CollegeID     int       `db:"college_id" json:"college_id"`
	MarksObtained float64   `db:"marks_obtained" json:"marks_obtained"`
	TotalMarks    float64   `db:"total_marks" json:"total_marks"`
	GradeLetter   *string   `db:"grade_letter" json:"grade_letter,omitempty"`
	Semester      int       `db:"semester" json:"semester"`
	AcademicYear  string    `db:"academic_year" json:"academic_year"` // e.g., "2023-2024"
	ExamType      string    `db:"exam_type" json:"exam_type"`         // e.g., "Midterm", "Final", "Assignment"
	GradedAt      time.Time `db:"graded_at" json:"graded_at"`         // When the grade was officially given
	Comments      *string   `db:"comments" json:"comments,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// GradeFilter can be used for querying lists of grades with specific criteria
type GradeFilter struct {
	StudentID    *string `json:"student_id,omitempty"`
	CourseID     *int    `json:"course_id,omitempty"`
	CollegeID    *int    `json:"college_id,omitempty"` // Essential for multi-tenancy
	Semester     *int    `json:"semester,omitempty"`
	AcademicYear *string `json:"academic_year,omitempty"`
	ExamType     *string `json:"exam_type,omitempty"`
	// Add pagination fields if needed
	Limit  uint64 `json:"limit,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}
