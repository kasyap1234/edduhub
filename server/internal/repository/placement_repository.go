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

type PlacementRepository interface {
	CreatePlacement(ctx context.Context, placement *models.Placement) error
	GetPlacementByID(ctx context.Context, collegeID int, placementID int) (*models.Placement, error)
	UpdatePlacement(ctx context.Context, placement *models.Placement) error
	DeletePlacement(ctx context.Context, collegeID int, placementID int) error

	// Find methods with pagination
	FindPlacementsByStudent(ctx context.Context, collegeID int, studentID int, limit, offset uint64) ([]*models.Placement, error)
	FindPlacementsByCollege(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Placement, error)
	FindPlacementsByCompany(ctx context.Context, collegeID int, companyName string, limit, offset uint64) ([]*models.Placement, error)

	// Count methods
	CountPlacementsByStudent(ctx context.Context, collegeID int, studentID int) (int, error)
	CountPlacementsByCollege(ctx context.Context, collegeID int) (int, error)
}

type placementRepository struct {
	DB *DB
}

func NewPlacementRepository(db *DB) PlacementRepository {
	return &placementRepository{DB: db}
}

const placementTable = "placements"

func (r *placementRepository) CreatePlacement(ctx context.Context, placement *models.Placement) error {
	now := time.Now()
	placement.CreatedAt = now
	placement.UpdatedAt = now

	query := r.DB.SQ.Insert(placementTable).
		Columns("college_id", "student_id", "company_name", "job_title", "package", "placement_date", "status", "created_at", "updated_at").
		Values(placement.CollegeID, placement.StudentID, placement.CompanyName, placement.JobTitle, placement.Package, placement.PlacementDate, placement.Status, placement.CreatedAt, placement.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreatePlacement: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&placement.ID)
	if err != nil {
		return fmt.Errorf("CreatePlacement: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *placementRepository) GetPlacementByID(ctx context.Context, collegeID int, placementID int) (*models.Placement, error) {
	query := r.DB.SQ.Select("id", "college_id", "student_id", "company_name", "job_title", "package", "placement_date", "status", "created_at", "updated_at").
		From(placementTable).
		Where(squirrel.Eq{"id": placementID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetPlacementByID: failed to build query: %w", err)
	}

	placement := &models.Placement{}
	err = pgxscan.Get(ctx, r.DB.Pool, placement, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetPlacementByID: placement with ID %d not found for college ID %d", placementID, collegeID)
		}
		return nil, fmt.Errorf("GetPlacementByID: failed to execute query or scan: %w", err)
	}
	return placement, nil
}

func (r *placementRepository) UpdatePlacement(ctx context.Context, placement *models.Placement) error {
	placement.UpdatedAt = time.Now()

	query := r.DB.SQ.Update(placementTable).
		Set("company_name", placement.CompanyName).
		Set("job_title", placement.JobTitle).
		Set("package", placement.Package).
		Set("placement_date", placement.PlacementDate).
		Set("status", placement.Status).
		Set("updated_at", placement.UpdatedAt).
		Where(squirrel.Eq{"id": placement.ID, "college_id": placement.CollegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdatePlacement: failed to build query: %w", err)
	}

	cmdTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdatePlacement: failed to execute query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdatePlacement: no placement found with ID %d for college ID %d, or no changes made", placement.ID, placement.CollegeID)
	}
	return nil
}

func (r *placementRepository) DeletePlacement(ctx context.Context, collegeID int, placementID int) error {
	query := r.DB.SQ.Delete(placementTable).
		Where(squirrel.Eq{"id": placementID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeletePlacement: failed to build query: %w", err)
	}

	cmdTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeletePlacement: failed to execute query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("DeletePlacement: no placement found with ID %d for college ID %d, or already deleted", placementID, collegeID)
	}
	return nil
}

func (r *placementRepository) findPlacements(ctx context.Context, whereClause squirrel.Sqlizer, limit, offset uint64) ([]*models.Placement, error) {
	query := r.DB.SQ.Select("id", "college_id", "student_id", "company_name", "job_title", "package", "placement_date", "status", "created_at", "updated_at").
		From(placementTable).
		Where(whereClause).
		OrderBy("placement_date DESC", "student_id ASC").
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("findPlacements: failed to build query: %w", err)
	}

	placements := []*models.Placement{}
	err = pgxscan.Select(ctx, r.DB.Pool, &placements, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("findPlacements: failed to execute query or scan: %w", err)
	}
	return placements, nil
}

func (r *placementRepository) FindPlacementsByStudent(ctx context.Context, collegeID int, studentID int, limit, offset uint64) ([]*models.Placement, error) {
	where := squirrel.Eq{"college_id": collegeID, "student_id": studentID}
	return r.findPlacements(ctx, where, limit, offset)
}

func (r *placementRepository) FindPlacementsByCollege(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Placement, error) {
	where := squirrel.Eq{"college_id": collegeID}
	return r.findPlacements(ctx, where, limit, offset)
}

func (r *placementRepository) FindPlacementsByCompany(ctx context.Context, collegeID int, companyName string, limit, offset uint64) ([]*models.Placement, error) {
	// Use ILIKE for case-insensitive search, adjust if case-sensitive is needed
	where := squirrel.And{
		squirrel.Eq{"college_id": collegeID},
		squirrel.ILike{"company_name": fmt.Sprintf("%%%s%%", companyName)}, // Wildcard search
	}
	return r.findPlacements(ctx, where, limit, offset)
}

func (r *placementRepository) countPlacements(ctx context.Context, whereClause squirrel.Sqlizer) (int, error) {
	query := r.DB.SQ.Select("COUNT(*)").From(placementTable).Where(whereClause)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("countPlacements: failed to build query: %w", err)
	}
	var count int
	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("countPlacements: failed to execute query or scan: %w", err)
	}
	return count, nil
}

func (r *placementRepository) CountPlacementsByStudent(ctx context.Context, collegeID int, studentID int) (int, error) {
	where := squirrel.Eq{"college_id": collegeID, "student_id": studentID}
	return r.countPlacements(ctx, where)
}

func (r *placementRepository) CountPlacementsByCollege(ctx context.Context, collegeID int) (int, error) {
	where := squirrel.Eq{"college_id": collegeID}
	return r.countPlacements(ctx, where)
}
