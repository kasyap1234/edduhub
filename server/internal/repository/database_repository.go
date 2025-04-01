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
}

func (d *BaseDatabaseRepository[T]) Delete(ctx context.Context, model *T) error {
	_, err := d.DB.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

// func(d*BaseDatabaseRepository[T])DeleteByID(ctx context.Context,id interface{})error {
// 	_,err :=d.DB.NewDelete()
// }
