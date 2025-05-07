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

type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, department *models.Department) error
	GetDepartmentByID(ctx context.Context, collegeID int, departmentID int) (*models.Department, error)
	GetDepartmentByName(ctx context.Context, collegeID int, name string) (*models.Department, error)
	UpdateDepartment(ctx context.Context, department *models.Department) error
	DeleteDepartment(ctx context.Context, collegeID int, departmentID int) error
	ListDepartmentsByCollege(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Department, error)
	CountDepartmentsByCollege(ctx context.Context, collegeID int) (int, error)
}

type departmentRepository struct {
	DB *DB
}

const departmentTable = "departments"

func NewDepartmentRepository(DB *DB) DepartmentRepository {
	return &departmentRepository{
		DB: DB,
	}
}

func (r *departmentRepository) CreateDepartment(ctx context.Context, department *models.Department) error {
	now := time.Now()
	department.CreatedAt = now
	department.UpdatedAt = now

	query := r.DB.SQ.Insert(departmentTable).
		Columns("college_id", "name", "hod", "created_at", "updated_at").
		Values(department.CollegeID, department.Name, department.HOD, department.CreatedAt, department.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateDepartment: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&department.ID)
	if err != nil {
		return fmt.Errorf("CreateDepartment: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *departmentRepository) GetDepartmentByID(ctx context.Context, collegeID int, departmentID int) (*models.Department, error) {
	department := &models.Department{}
	query := r.DB.SQ.Select("id", "college_id", "name", "hod", "created_at", "updated_at").
		From(departmentTable).
		Where(squirrel.Eq{"id": departmentID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetDepartmentByID: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, department, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetDepartmentByID: department with ID %d not found for college ID %d", departmentID, collegeID)
		}
		return nil, fmt.Errorf("GetDepartmentByID: failed to execute query or scan: %w", err)
	}
	return department, nil
}

func (r *departmentRepository) GetDepartmentByName(ctx context.Context, collegeID int, name string) (*models.Department, error) {
	department := &models.Department{}
	query := r.DB.SQ.Select("id", "college_id", "name", "hod", "created_at", "updated_at").
		From(departmentTable).
		Where(squirrel.Eq{"name": name, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetDepartmentByName: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, department, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetDepartmentByName: department with name '%s' not found for college ID %d", name, collegeID)
		}
		return nil, fmt.Errorf("GetDepartmentByName: failed to execute query or scan: %w", err)
	}
	return department, nil
}

func (r *departmentRepository) UpdateDepartment(ctx context.Context, department *models.Department) error {
	department.UpdatedAt = time.Now()
	query := r.DB.SQ.Update(departmentTable).
		Set("name", department.Name).
		Set("hod", department.HOD).
		Set("updated_at", department.UpdatedAt).
		Where(squirrel.Eq{"id": department.ID, "college_id": department.CollegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateDepartment: failed to build query: %w", err)
	}

	cmdTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateDepartment: failed to execute query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateDepartment: no department found with ID %d for college ID %d, or no changes made", department.ID, department.CollegeID)
	}
	return nil
}

func (r *departmentRepository) DeleteDepartment(ctx context.Context, collegeID int, departmentID int) error {
	query := r.DB.SQ.Delete(departmentTable).
		Where(squirrel.Eq{"id": departmentID, "college_id": collegeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteDepartment: failed to build query: %w", err)
	}

	cmdTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteDepartment: failed to execute query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteDepartment: no department found with ID %d for college ID %d, or already deleted", departmentID, collegeID)
	}
	return nil
}

func (r *departmentRepository) ListDepartmentsByCollege(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Department, error) {
	departments := []*models.Department{}
	query := r.DB.SQ.Select("id", "college_id", "name", "hod", "created_at", "updated_at").
		From(departmentTable).
		Where(squirrel.Eq{"college_id": collegeID}).
		OrderBy("name ASC").Limit(limit).Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ListDepartmentsByCollege: failed to build query: %w", err)
	}

	err = pgxscan.Select(ctx, r.DB.Pool, &departments, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ListDepartmentsByCollege: failed to execute query or scan: %w", err)
	}
	return departments, nil
}

func (r *departmentRepository) CountDepartmentsByCollege(ctx context.Context, collegeID int) (int, error) {
	var count int
	query := r.DB.SQ.Select("COUNT(*)").From(departmentTable).Where(squirrel.Eq{"college_id": collegeID})
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("CountDepartmentsByCollege: build query: %w", err)
	}
	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountDepartmentsByCollege: exec/scan: %w", err)
	}
	return count, nil
}
