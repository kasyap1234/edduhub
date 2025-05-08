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

const timeTableBlockTable = "timetable_blocks"

var timeTableBlockQueryFields = []string{
	"id", "college_id", "department_id", "course_id", "class_id",
	"day_of_week", "start_time", "end_time", "room_number", "faculty_id",
	"created_at", "updated_at",
}

type TimeTableRepository interface {
	CreateTimeTableBlock(ctx context.Context, block *models.TimeTableBlock) error
	GetTimeTableBlockByID(ctx context.Context, blockID int, collegeID int) (*models.TimeTableBlock, error)
	UpdateTimeTableBlock(ctx context.Context, block *models.TimeTableBlock) error
	DeleteTimeTableBlock(ctx context.Context, blockID int, collegeID int) error
	GetTimeTableBlocks(ctx context.Context, filter models.TimeTableBlockFilter) ([]*models.TimeTableBlock, error)
	CountTimeTableBlocks(ctx context.Context, filter models.TimeTableBlockFilter) (int, error)
}

type timetableRepository struct {
	DB *DB
}

func NewTimeTableRepository(db *DB) TimeTableRepository {
	return &timetableRepository{DB: db}
}

func (r *timetableRepository) CreateTimeTableBlock(ctx context.Context, block *models.TimeTableBlock) error {
	now := time.Now()
	block.CreatedAt = now
	block.UpdatedAt = now

	query := r.DB.SQ.Insert(timeTableBlockTable).
		Columns(
			"college_id", "department_id", "course_id", "class_id",
			"day_of_week", "start_time", "end_time", "room_number", "faculty_id",
			"created_at", "updated_at",
		).
		Values(
			block.CollegeID, block.DepartmentID, block.CourseID, block.ClassID,
			block.DayOfWeek, block.StartTime, block.EndTime, block.RoomNumber, block.FacultyID,
			block.CreatedAt, block.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateTimeTableBlock: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&block.ID)
	if err != nil {
		return fmt.Errorf("CreateTimeTableBlock: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *timetableRepository) GetTimeTableBlockByID(ctx context.Context, blockID int, collegeID int) (*models.TimeTableBlock, error) {
	block := &models.TimeTableBlock{}
	query := r.DB.SQ.Select(timeTableBlockQueryFields...).
		From(timeTableBlockTable).
		Where(squirrel.Eq{"id": blockID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetTimeTableBlockByID: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, block, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetTimeTableBlockByID: block with ID %d for college ID %d not found: %w", blockID, collegeID, err)
		}
		return nil, fmt.Errorf("GetTimeTableBlockByID: failed to execute query or scan: %w", err)
	}
	return block, nil
}

func (r *timetableRepository) UpdateTimeTableBlock(ctx context.Context, block *models.TimeTableBlock) error {
	block.UpdatedAt = time.Now()

	query := r.DB.SQ.Update(timeTableBlockTable).
		Set("department_id", block.DepartmentID).
		Set("course_id", block.CourseID).
		Set("class_id", block.ClassID).
		Set("day_of_week", block.DayOfWeek).
		Set("start_time", block.StartTime).
		Set("end_time", block.EndTime).
		Set("room_number", block.RoomNumber).
		Set("faculty_id", block.FacultyID).
		Set("updated_at", block.UpdatedAt).
		Where(squirrel.Eq{"id": block.ID, "college_id": block.CollegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateTimeTableBlock: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateTimeTableBlock: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateTimeTableBlock: no block found with ID %d for college ID %d, or no changes made", block.ID, block.CollegeID)
	}
	return nil
}

func (r *timetableRepository) DeleteTimeTableBlock(ctx context.Context, blockID int, collegeID int) error {
	query := r.DB.SQ.Delete(timeTableBlockTable).
		Where(squirrel.Eq{"id": blockID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteTimeTableBlock: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteTimeTableBlock: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteTimeTableBlock: no block found with ID %d for college ID %d", blockID, collegeID)
	}
	return nil
}

func (r *timetableRepository) applyTimeTableBlockFilter(query squirrel.SelectBuilder, filter models.TimeTableBlockFilter) squirrel.SelectBuilder {
	query = query.Where(squirrel.Eq{"college_id": filter.CollegeID}) // CollegeID is mandatory

	if filter.DepartmentID != nil {
		query = query.Where(squirrel.Eq{"department_id": *filter.DepartmentID})
	}
	if filter.CourseID != nil {
		query = query.Where(squirrel.Eq{"course_id": *filter.CourseID})
	}
	if filter.ClassID != nil {
		query = query.Where(squirrel.Eq{"class_id": *filter.ClassID})
	}
	if filter.DayOfWeek != nil {
		query = query.Where(squirrel.Eq{"day_of_week": *filter.DayOfWeek})
	}
	if filter.FacultyID != nil {
		query = query.Where(squirrel.Eq{"faculty_id": *filter.FacultyID})
	}
	return query
}

func (r *timetableRepository) GetTimeTableBlocks(ctx context.Context, filter models.TimeTableBlockFilter) ([]*models.TimeTableBlock, error) {
	if filter.CollegeID == 0 { // Or handle as pointer and check for nil
		return nil, errors.New("GetTimeTableBlocks: CollegeID filter is required")
	}

	query := r.DB.SQ.Select(timeTableBlockQueryFields...).From(timeTableBlockTable)
	query = r.applyTimeTableBlockFilter(query, filter)
	query = query.OrderBy("day_of_week ASC", "start_time ASC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetTimeTableBlocks: failed to build query: %w", err)
	}

	var blocks []*models.TimeTableBlock
	err = pgxscan.Select(ctx, r.DB.Pool, &blocks, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.TimeTableBlock{}, nil
		}
		return nil, fmt.Errorf("GetTimeTableBlocks: failed to execute query or scan: %w", err)
	}
	return blocks, nil
}

func (r *timetableRepository) CountTimeTableBlocks(ctx context.Context, filter models.TimeTableBlockFilter) (int, error) {
	if filter.CollegeID == 0 {
		return 0, errors.New("CountTimeTableBlocks: CollegeID filter is required")
	}

	query := r.DB.SQ.Select("COUNT(*)").From(timeTableBlockTable)
	query = r.applyTimeTableBlockFilter(query, filter)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("CountTimeTableBlocks: failed to build query: %w", err)
	}

	var count int
	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountTimeTableBlocks: failed to execute query or scan: %w", err)
	}
	return count, nil
}
