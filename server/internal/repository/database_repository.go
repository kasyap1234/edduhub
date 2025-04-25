package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type DatabaseRepository[T any] interface {
	Create(ctx context.Context, model *T) error
	CreateMany(ctx context.Context, models []*T) error
	FindByID(ctx context.Context, id interface{}) (*T, error)
	FindOne(ctx context.Context, query string, args ...interface{}) (*T, error)
	FindAll(ctx context.Context) ([]*T, error)
	FindWhere(ctx context.Context, query string, args ...interface{}) ([]*T, error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, model *T) error
	Count(ctx context.Context, query string, args ...interface{}) (int, error)
	// transaction
	Exists(ctx context.Context, model *T, query string, args ...interface{}) (bool, error)
}

type BaseDatabaseRepository[T any] struct {
	DB *bun.DB
}

func NewBaseRepository[T any](db *bun.DB) DatabaseRepository[T] {
	return &BaseDatabaseRepository[T]{
		DB: db,
	}
}

func (d *BaseDatabaseRepository[T]) Create(ctx context.Context, model *T) error {
	_, err := d.DB.NewInsert().Model(model).Exec(ctx)
	return err
}

func (d *BaseDatabaseRepository[T]) CreateMany(ctx context.Context, models []*T) error {
	_, err := d.DB.NewInsert().Model(&models).Exec(ctx)
	return err
}

func (d *BaseDatabaseRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	model := new(T)
	err := d.DB.NewSelect().Model(model).Where("id= ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (d *BaseDatabaseRepository[T]) FindOne(ctx context.Context, query string, args ...interface{}) (*T, error) {
	model := new(T)
	err := d.DB.NewSelect().Model(model).Where(query, args...).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (d *BaseDatabaseRepository[T]) FindAll(ctx context.Context) ([]*T, error) {
	var models []*T
	err := d.DB.NewSelect().Model(&models).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (d *BaseDatabaseRepository[T]) FindWhere(ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	var models []*T
	err := d.DB.NewSelect().Model(&models).Where(query, args...).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (d *BaseDatabaseRepository[T]) Update(ctx context.Context, model *T) error {
	_, err := d.DB.NewUpdate().Model(model).WherePK().Exec(ctx)
	return err
}

func (d *BaseDatabaseRepository[T]) Count(ctx context.Context, query string, args ...interface{}) (int, error) {
	count, err := d.DB.NewSelect().Model((*T)(nil)).Where(query, args...).Count(ctx)
	return count, err
}

func (d *BaseDatabaseRepository[T]) Delete(ctx context.Context, model *T) error {
	_, err := d.DB.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (d *BaseDatabaseRepository[T]) Exists(ctx context.Context, model *T, query string, args ...interface{}) (bool, error) {
	exists, err := d.DB.NewSelect().Model(&model).Where(query, args...).Exists(ctx)
	return exists, err
}

//Student <-> Course: This is Many-to-Many. You need to create  the enrollments table/model. This is the standard and necessary approach.
//Course <-> Lecture: This is One-to-Many (One Course, Many Lectures). You correctly have course_id in the lectures table. No join table needed. .Relation() works directly.
//Student <-> Attendance: One-to-Many. Correctly have student_id in attendance. No join table needed. .Relation() works directly.
//Course <-> Attendance: One-to-Many. Correctly have course_id in attendance. No join table needed. .Relation() works directly.
//Lecture <-> Attendance: One-to-Many. Correctly have lecture_id in attendance. No join table needed. .Relation() works directly.
// Permissions (Keto): Your assignment_helper.go uses Keto for permissions (e.g., Faculty-Course, Student-Course). Keto acts as an external system managing these relationships (relation tuples). You do not need database join tables for these specific permission relationships because Keto handles them.

// only need to create a dedicated join table (like enrollments) when you have a Many-to-Many relationship between two core entities in your database.

// For One-to-One and One-to-Many relationships, a simple foreign key in one of the tables is sufficient.

// Bun's .Relation() feature is designed to work seamlessly with both scenarios:

// It uses the foreign keys directly for 1:1 and 1:N relationships.
// It traverses through the join table model for M:N relationships.
// Therefore, for other features, analyze the relationship type:

// If it's M:N, create a join table/model like you did for enrollments.
// If it's 1:1 or 1:N, just add the appropriate foreign key column to one of the existing models/tables.
// You don't need to create extra join tables for relationships that are already handled correctly by foreign keys (like Course-Lecture, Student-Attendance, etc.).
