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

const calendarBlockTable = "calendar_blocks"

var calendarBlockQueryFields = []string{
	"id", "college_id", "title", "description", "event_type", "date",
	"created_at", "updated_at",
}

type CalendarRepository interface {
	CreateCalendarBlock(ctx context.Context, block *models.CalendarBlock) error
	GetCalendarBlockByID(ctx context.Context, blockID int, collegeID int) (*models.CalendarBlock, error)
	UpdateCalendarBlock(ctx context.Context, block *models.CalendarBlock) error
	DeleteCalendarBlock(ctx context.Context, blockID int, collegeID int) error
	GetCalendarBlocks(ctx context.Context, filter models.CalendarBlockFilter) ([]*models.CalendarBlock, error)
	CountCalendarBlocks(ctx context.Context, filter models.CalendarBlockFilter) (int, error)
}

type calendarRepository struct {
	DB *DB
}

func NewCalendarRepository(db *DB) CalendarRepository {
	return &calendarRepository{DB: db}
}

func (r *calendarRepository) CreateCalendarBlock(ctx context.Context, block *models.CalendarBlock) error {
	now := time.Now()
	block.CreatedAt = now
	block.UpdatedAt = now

	query := r.DB.SQ.Insert(calendarBlockTable).
		Columns("college_id", "title", "description", "event_type", "date", "created_at", "updated_at").
		Values(block.CollegeID, block.Title, block.Description, block.EventType, block.Date, block.CreatedAt, block.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateCalendarBlock: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&block.ID)
	if err != nil {
		return fmt.Errorf("CreateCalendarBlock: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *calendarRepository) GetCalendarBlockByID(ctx context.Context, blockID int, collegeID int) (*models.CalendarBlock, error) {
	block := &models.CalendarBlock{}
	query := r.DB.SQ.Select(calendarBlockQueryFields...).
		From(calendarBlockTable).
		Where(squirrel.Eq{"id": blockID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetCalendarBlockByID: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, block, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetCalendarBlockByID: block with ID %d for college ID %d not found: %w", blockID, collegeID, err)
		}
		return nil, fmt.Errorf("GetCalendarBlockByID: failed to execute query or scan: %w", err)
	}
	return block, nil
}

func (r *calendarRepository) UpdateCalendarBlock(ctx context.Context, block *models.CalendarBlock) error {
	block.UpdatedAt = time.Now()

	query := r.DB.SQ.Update(calendarBlockTable).
		Set("title", block.Title).
		Set("description", block.Description).
		Set("event_type", block.EventType).
		Set("date", block.Date).
		Set("updated_at", block.UpdatedAt).
		Where(squirrel.Eq{"id": block.ID, "college_id": block.CollegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateCalendarBlock: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateCalendarBlock: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateCalendarBlock: no block found with ID %d for college ID %d, or no changes made", block.ID, block.CollegeID)
	}
	return nil
}

func (r *calendarRepository) DeleteCalendarBlock(ctx context.Context, blockID int, collegeID int) error {
	query := r.DB.SQ.Delete(calendarBlockTable).
		Where(squirrel.Eq{"id": blockID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteCalendarBlock: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteCalendarBlock: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteCalendarBlock: no block found with ID %d for college ID %d", blockID, collegeID)
	}
	return nil
}

func (r *calendarRepository) applyCalendarBlockFilter(query squirrel.SelectBuilder, filter models.CalendarBlockFilter) squirrel.SelectBuilder {
	if filter.CollegeID == nil { // This check should ideally be done before calling this helper
		// Or, if CollegeID is mandatory for all queries using this filter, make it non-pointer in the struct.
		// For now, we assume the caller ensures CollegeID is present if required.
		// However, for GetCalendarBlocks, we'll enforce it.
	} else {
		query = query.Where(squirrel.Eq{"college_id": *filter.CollegeID})
	}

	if filter.EventType != nil {
		query = query.Where(squirrel.Eq{"event_type": *filter.EventType})
	}
	if filter.StartDate != nil {
		query = query.Where(squirrel.GtOrEq{"date": *filter.StartDate})
	}
	if filter.EndDate != nil {
		query = query.Where(squirrel.LtOrEq{"date": *filter.EndDate})
	}
	return query
}

func (r *calendarRepository) GetCalendarBlocks(ctx context.Context, filter models.CalendarBlockFilter) ([]*models.CalendarBlock, error) {
	if filter.CollegeID == nil {
		return nil, errors.New("GetCalendarBlocks: CollegeID filter is required")
	}

	query := r.DB.SQ.Select(calendarBlockQueryFields...).From(calendarBlockTable)
	query = r.applyCalendarBlockFilter(query, filter)
	query = query.OrderBy("date ASC", "created_at ASC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetCalendarBlocks: failed to build query: %w", err)
	}

	var blocks []*models.CalendarBlock
	err = pgxscan.Select(ctx, r.DB.Pool, &blocks, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.CalendarBlock{}, nil
		}
		return nil, fmt.Errorf("GetCalendarBlocks: failed to execute query or scan: %w", err)
	}
	return blocks, nil
}

func (r *calendarRepository) CountCalendarBlocks(ctx context.Context, filter models.CalendarBlockFilter) (int, error) {
	if filter.CollegeID == nil {
		return 0, errors.New("CountCalendarBlocks: CollegeID filter is required")
	}

	query := r.DB.SQ.Select("COUNT(*)").From(calendarBlockTable)
	query = r.applyCalendarBlockFilter(query, filter)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("CountCalendarBlocks: failed to build query: %w", err)
	}

	var count int
	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountCalendarBlocks: failed to execute query or scan: %w", err)
	}
	return count, nil
}
