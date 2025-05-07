package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type LectureRepository interface {
	CreateLecture(ctx context.Context, lecture *models.Lecture) error
	GetLectureByID(ctx context.Context, collegeID int, lectureID int) (*models.Lecture, error)
	UpdateLecture(ctx context.Context, lecture *models.Lecture) error
	DeleteLecture(ctx context.Context, collegeID int, lectureID int) error

	// Finder methods
	FindLecturesByCourse(ctx context.Context, collegeID int, courseID int, limit, offset uint64) ([]*models.Lecture, error)
	CountLecturesByCourse(ctx context.Context, collegeID int, courseID int) (int, error)
	// Add more finders as needed, e.g., FindLecturesByDateRange, FindLecturesByInstructor (if lectures are directly linked to instructors)
}

type lectureRepository struct {
	DB *DB
}

func NewLectureRepository(db *DB) LectureRepository {
	return &lectureRepository{DB: db}
}

const lectureTable = "lectures"

func (r *lectureRepository) CreateLecture(ctx context.Context, lecture *models.Lecture) error {
	now := time.Now()
	lecture.CreatedAt = now
	lecture.UpdatedAt = now

	query := r.DB.SQ.Insert(lectureTable).
		Columns("course_id", "college_id", "title", "description", "start_time", "end_time", "meeting_link", "created_at", "updated_at").
		Values(lecture.CourseID, lecture.CollegeID, lecture.Title, lecture.Description, lecture.StartTime, lecture.EndTime, lecture.MeetingLink, lecture.CreatedAt, lecture.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateLecture: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&lecture.ID)
	if err != nil {
		// Consider checking for specific DB errors like foreign key violations
		return fmt.Errorf("CreateLecture: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *lectureRepository) GetLectureByID(ctx context.Context, collegeID int, lectureID int) (*models.Lecture, error) {
	query := r.DB.SQ.Select("id", "course_id", "college_id", "title", "description", "start_time", "end_time", "meeting_link", "created_at", "updated_at").
		From(lectureTable).
		Where(squirrel.Eq{"id": lectureID, "college_id": collegeID}) // Ensure lecture belongs to the specified college

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetLectureByID: failed to build query: %w", err)
	}

	lecture := &models.Lecture{}
	err = pgxscan.Get(ctx, r.DB.Pool, lecture, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetLectureByID: lecture with ID %d not found for college ID %d", lectureID, collegeID)
		}
		return nil, fmt.Errorf("GetLectureByID: failed to execute query or scan: %w", err)
	}
	return lecture, nil
}

func (r *lectureRepository) UpdateLecture(ctx context.Context, lecture *models.Lecture) error {
	lecture.UpdatedAt = time.Now()

	query := r.DB.SQ.Update(lectureTable).
		Set("title", lecture.Title).
		Set("description", lecture.Description).
		Set("start_time", lecture.StartTime).
		Set("end_time", lecture.EndTime).
		Set("meeting_link", lecture.MeetingLink).
		Set("course_id", lecture.CourseID). // Allow course_id to be updated if necessary
		Set("updated_at", lecture.UpdatedAt).
		Where(squirrel.Eq{"id": lecture.ID, "college_id": lecture.CollegeID}) // Ensure update is scoped

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateLecture: failed to build query: %w", err)
	}

	cmdTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateLecture: failed to execute query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateLecture: no lecture found with ID %d for college ID %d, or no changes made", lecture.ID, lecture.CollegeID)
	}
	return nil
}

func (r *lectureRepository) DeleteLecture(ctx context.Context, collegeID int, lectureID int) error {
	query := r.DB.SQ.Delete(lectureTable).
		Where(squirrel.Eq{"id": lectureID, "college_id": collegeID}) // Ensure deletion is scoped

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteLecture: failed to build query: %w", err)
	}

	cmdTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		// Consider foreign key constraint errors (e.g., if attendance records exist)
		// These should ideally be handled at the service layer (e.g., prevent deletion or cascade)
		return fmt.Errorf("DeleteLecture: failed to execute query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteLecture: no lecture found with ID %d for college ID %d, or already deleted", lectureID, collegeID)
	}
	return nil
}

func (r *lectureRepository) FindLecturesByCourse(ctx context.Context, collegeID int, courseID int, limit, offset uint64) ([]*models.Lecture, error) {
	query := r.DB.SQ.Select("id", "course_id", "college_id", "title", "description", "start_time", "end_time", "meeting_link", "created_at", "updated_at").
		From(lectureTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
			"course_id":  courseID,
		}).
		OrderBy("start_time ASC"). // Order lectures chronologically
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindLecturesByCourse: failed to build query: %w", err)
	}

	lectures := []*models.Lecture{}
	err = pgxscan.Select(ctx, r.DB.Pool, &lectures, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FindLecturesByCourse: failed to execute query or scan: %w", err)
	}
	return lectures, nil
}

func (r *lectureRepository) CountLecturesByCourse(ctx context.Context, collegeID int, courseID int) (int, error) {
	query := r.DB.SQ.Select("COUNT(*)").
		From(lectureTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
			"course_id":  courseID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("CountLecturesByCourse: failed to build query: %w", err)
	}

	var count int
	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountLecturesByCourse: failed to execute query or scan: %w", err)
	}
	return count, nil
}
