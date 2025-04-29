package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB represents the database connection
type DB struct {
	Pool *pgxpool.Pool

	SQ squirrel.StatementBuilderType
}

// NewDB creates a new database connection
func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		Pool: pool,
		SQ:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
