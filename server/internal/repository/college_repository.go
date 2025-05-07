package repository

import (
	"context"
	"eduhub/server/internal/models"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type CollegeRepository interface {
	CreateCollege(ctx context.Context, college *models.College) error
	GetCollegeByID(ctx context.Context, id int) (*models.College, error)
	GetCollegeByName(ctx context.Context, name string) (*models.College, error)
	UpdateCollege(ctx context.Context, college *models.College) error
	DeleteCollege(ctx context.Context, id int) error
	ListColleges(ctx context.Context, limit, offset uint64) ([]*models.College, error)
}

type collegeRepository struct {
	DB *DB
}

const collegeTable = "college"

func NewCollegeRepository(DB *DB) CollegeRepository {
	return &collegeRepository{
		DB: DB,
	}
}

// type College struct {
// 	ID        int       `db:"id" json:"id"`
// 	Name      string    `db:"name" json:"name"`
// 	Address   string    `db:"address" json:"address"`
// 	City      string    `db:"city" json:"city"`
// 	State     string    `db:"state" json:"state"`
// 	Country   string    `db:"country" json:"country"`
// 	CreatedAt time.Time `db:"created_at" json:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

// 	// Relations - not stored in DB
// 	Students []*Student `db:"-" json:"students,omitempty"`
// }

func (c *collegeRepository) CreateCollege(ctx context.Context, college *models.College) error {

	now := time.Now()
	college.CreatedAt = now
	college.UpdatedAt = now

	query := c.DB.SQ.Insert(collegeTable).Columns("name", "address", "city", "state", "country", "created_at", "updated_at").Values(college.Name, college.Address, college.City, college.State, college.Country, college.CreatedAt, college.UpdatedAt).Suffix("RETURNING id")
	sql, args, err := query.ToSql()
	if err != nil {
		return errors.New("create college: failed  to build query")

	}
	err = c.DB.Pool.QueryRow(ctx, sql, args...).Scan(&college.ID)
	if err != nil {
		return errors.New("unable to create college")
	}
	return nil
}

func (c *collegeRepository) GetCollegeByID(ctx context.Context, id int) (*models.College, error) {
	query := c.DB.SQ.Select("id", "name", "address", "city", "state", "country", "created_at", "updated_at").From(collegeTable).Where(squirrel.Eq{
		"id": id,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.New("failed to build query")
	}
	college := &models.College{}
	findErr := pgxscan.Get(ctx, c.DB.Pool, college, sql, args...)
	if findErr != nil {
		if findErr == pgx.ErrNoRows {
			return nil, errors.New("GetCollegeByid id not found")
		}
		return nil, errors.New("GetCollegeID: failed to execute query")
	}
	return college, nil
}

func (c *collegeRepository) GetCollegeByName(ctx context.Context, name string) (*models.College, error) {
	query := c.DB.SQ.Select("id", "name", "address", "city", "state", "country", "created_at", "updated_at").From(collegeTable).Where(squirrel.Eq{
		"name": name,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.New("failed to build query")
	}
	college := &models.College{}
	findErr := pgxscan.Get(ctx, c.DB.Pool, college, sql, args...)
	if findErr != nil {
		if findErr == pgx.ErrNoRows {
			return nil, errors.New("GetCollegeByid id not found")
		}
		return nil, errors.New("GetCollegeID: failed to execute query")
	}
	return college, nil
}

// type College struct {
// 	ID        int       `db:"id" json:"id"`
// 	Name      string    `db:"name" json:"name"`
// 	Address   string    `db:"address" json:"address"`
// 	City      string    `db:"city" json:"city"`
// 	State     string    `db:"state" json:"state"`
// 	Country   string    `db:"country" json:"country"`
// 	CreatedAt time.Time `db:"created_at" json:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

// 	// Relations - not stored in DB
// 	Students []*Student `db:"-" json:"students,omitempty"`
// }

func (c *collegeRepository) UpdateCollege(ctx context.Context, college *models.College) error {
	college.UpdatedAt = time.Now()
	query := c.DB.SQ.Update(collegeTable).Set("name", college.Name).Set("address", college.Address).Set("city", college.City).Set("updated_at", college.UpdatedAt).Where(squirrel.Eq{
		"id": college.ID,
	})
	sql, args, err := query.ToSql()
	if err != nil {
		return errors.New("update college: failed to build query")
	}
	commandTag, err := c.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.New("update college: faile to execute query")
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("update college: no college found with the requested id")

	}

	return nil

}

func (c *collegeRepository) DeleteCollege(ctx context.Context, id int) error {
	query := c.DB.SQ.Delete(collegeTable).Where(squirrel.Eq{
		"id": id,
	})
	sql, args, err := query.ToSql()
	if err != nil {
		return errors.New("failed to build query")
	}
	commandTag, err := c.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.New("failed to execute query")
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no enrollment found")
	}
	return nil
}

func (c *collegeRepository) ListColleges(ctx context.Context, limit, offset uint64) ([]*models.College, error) {

	query := c.DB.SQ.Select("id", "name", "address", "city", "state", "country", "created_at", "updated_at").From(collegeTable)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.New("list college: failed to execute query")
	}
	colleges := []*models.College{}
	err = pgxscan.Select(ctx, c.DB.Pool, &colleges, sql, args...)
	if err != nil {
		return nil, errors.New("list colleges: failed to execute query")
	}
	return colleges, nil
}
