package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4" // For pgx.ErrNoRows
)

const gradeTable = "grades"

var gradeQueryFields = []string{
	"id", "student_id", "course_id", "college_id", "marks_obtained", "total_marks",
	"grade_letter", "semester", "academic_year", "exam_type", "graded_at",
	"comments", "created_at", "updated_at",
}

type GradeRepository interface {
	CreateGrade(ctx context.Context, grade *models.Grade) error
	GetGradeByID(ctx context.Context, gradeID int, collegeID int) (*models.Grade, error)
	UpdateGrade(ctx context.Context, grade *models.Grade) error
	DeleteGrade(ctx context.Context, gradeID int, collegeID int) error
	GetGrades(ctx context.Context, filter models.GradeFilter) ([]*models.Grade, error)
	// GetStudentProgress and GenerateStudentReport might be higher-level service methods
	// or more complex queries. For now, GetGrades with filters can serve many needs.
	GetGradesByCourse(ctx context.Context, collegeID int, courseID int) ([]*models.Grade, error)
	GetGradesByStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Grade, error)
}

type gradeRepository struct {
	DB *DB
}

func NewGradeRepository(db *DB) GradeRepository {
	return &gradeRepository{DB: db}
}

func (r *gradeRepository) CreateGrade(ctx context.Context, grade *models.Grade) error {
	now := time.Now()
	grade.CreatedAt = now
	grade.UpdatedAt = now
	if grade.GradedAt.IsZero() { // Default GradedAt to now if not provided
		grade.GradedAt = now
	}

	query := r.DB.SQ.Insert(gradeTable).
		Columns(
			"student_id", "course_id", "college_id", "marks_obtained", "total_marks",
			"grade_letter", "semester", "academic_year", "exam_type", "graded_at",
			"comments", "created_at", "updated_at",
		).
		Values(
			grade.StudentID, grade.CourseID, grade.CollegeID, grade.MarksObtained, grade.TotalMarks,
			grade.GradeLetter, grade.Semester, grade.AcademicYear, grade.ExamType, grade.GradedAt,
			grade.Comments, grade.CreatedAt, grade.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateGrade: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&grade.ID)
	if err != nil {
		// Consider specific error handling for duplicate entries or foreign key violations
		return fmt.Errorf("CreateGrade: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *gradeRepository) GetGradeByID(ctx context.Context, gradeID int, collegeID int) (*models.Grade, error) {
	grade := &models.Grade{}
	query := r.DB.SQ.Select(gradeQueryFields...).
		From(gradeTable).
		Where(squirrel.Eq{"id": gradeID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetGradeByID: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, grade, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetGradeByID: grade with ID %d for college ID %d not found: %w", gradeID, collegeID, err)
		}
		return nil, fmt.Errorf("GetGradeByID: failed to execute query or scan: %w", err)
	}
	return grade, nil
}

func (r *gradeRepository) UpdateGrade(ctx context.Context, grade *models.Grade) error {
	grade.UpdatedAt = time.Now()

	query := r.DB.SQ.Update(gradeTable).
		Set("student_id", grade.StudentID). // StudentID might not be updatable, depends on rules
		Set("course_id", grade.CourseID).   // CourseID might not be updatable
		Set("marks_obtained", grade.MarksObtained).
		Set("total_marks", grade.TotalMarks).
		Set("grade_letter", grade.GradeLetter).
		Set("semester", grade.Semester).
		Set("academic_year", grade.AcademicYear).
		Set("exam_type", grade.ExamType).
		Set("graded_at", grade.GradedAt).
		Set("comments", grade.Comments).
		Set("updated_at", grade.UpdatedAt).
		Where(squirrel.Eq{"id": grade.ID, "college_id": grade.CollegeID}) // Ensure update is for the correct college

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateGrade: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateGrade: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateGrade: no grade found with ID %d for college ID %d, or no changes made", grade.ID, grade.CollegeID)
	}
	return nil
}

func (r *gradeRepository) DeleteGrade(ctx context.Context, gradeID int, collegeID int) error {
	query := r.DB.SQ.Delete(gradeTable).
		Where(squirrel.Eq{"id": gradeID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteGrade: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteGrade: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteGrade: no grade found with ID %d for college ID %d", gradeID, collegeID)
	}
	return nil
}

func (r *gradeRepository) GetGrades(ctx context.Context, filter models.GradeFilter) ([]*models.Grade, error) {
	query := r.DB.SQ.Select(gradeQueryFields...).From(gradeTable)

	if filter.CollegeID == nil {
		return nil, errors.New("GetGrades: CollegeID filter is required")
	}
	query = query.Where(squirrel.Eq{"college_id": *filter.CollegeID})

	if filter.StudentID != nil {
		query = query.Where(squirrel.Eq{"student_id": *filter.StudentID})
	}
	if filter.CourseID != nil {
		query = query.Where(squirrel.Eq{"course_id": *filter.CourseID})
	}
	if filter.Semester != nil {
		query = query.Where(squirrel.Eq{"semester": *filter.Semester})
	}
	if filter.AcademicYear != nil {
		query = query.Where(squirrel.Eq{"academic_year": *filter.AcademicYear})
	}
	if filter.ExamType != "" {
		query = query.Where(squirrel.Eq{"exam_type": filter.ExamType})
	}

	// For progress tracking, you might want to order by academic_year, semester, graded_at
	query = query.OrderBy("academic_year ASC", "semester ASC", "graded_at ASC") // Default ordering

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetGrades: failed to build query: %w", err)
	}

	var grades []*models.Grade
	err = pgxscan.Select(ctx, r.DB.Pool, &grades, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.Grade{}, nil // Return empty slice if no rows found
		}
		return nil, fmt.Errorf("GetGrades: failed to execute query or scan: %w", err)
	}
	return grades, nil
}

func (r *gradeRepository) GetGradesByStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Grade, error) {
	studentIDStr := string(studentID)
	filter := models.GradeFilter{
		StudentID: &studentIDStr,
		CollegeID: &collegeID,
	}
	return r.GetGrades(ctx, filter)
}

func (r *gradeRepository) GetGradesByCourse(ctx context.Context, collegeID int, courseID int) ([]*models.Grade, error) {
	// courseIDStr := string(courseID)
	filter := models.GradeFilter{
		CourseID:  &courseID,
		CollegeID: &collegeID,
	}
	return r.GetGrades(ctx, filter)
}
