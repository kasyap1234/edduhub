// package repository

// import (
// 	"github.com/Masterminds/squirrel"
// 	"github.com/jackc/pgx/v4/pgxpool"
// )

// // DB represents the database connection
// type DB struct {
// 	Pool *pgxpool.Pool

// 	SQ squirrel.StatementBuilderType
// }

// // NewDB creates a new database connection
// func NewDB(pool *pgxpool.Pool) *DB {
// 	return &DB{
// 		Pool: pool,
// 		SQ:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
// 	}
// }

// // Close closes the database connection

//	func (db *DB) Close() {
//		if db.Pool != nil {
//			db.Pool.Close()
//		}
//	}
package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn" // Import pgconn for CommandTag
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	// Import pgx for Rows, Row
	// Keep this import if you still use *pgxpool.Pool directly somewhere (e.g. in NewDB)
)

// PoolIface defines the methods of pgxpool.Pool used by our repositories.
// Your production DB struct will now depend on this interface.
// --- Define this interface in your production code ---
type PoolIface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error) // Use pgconn.CommandTag
	Close()                                                                                      // Include Close if used by DB.Close()
	// Add other methods from pgxpool.Pool if *any* repository method calls them
	// e.g., Acquire(ctx context.Context) (*pgxpool.Conn, error)
	//      Stat() *pgxpool.Stat
}

// DB represents the database connection structure used by repositories
type DB struct {
	Pool PoolIface

	SQ squirrel.StatementBuilderType
}

// NewDB creates a new database connection structure
// Now accepts the interface. In production, you will pass the real *pgxpool.Pool instance
// (as *pgxpool.Pool implements PoolIface). In tests, you pass the mock.
func NewDB(pool *pgxpool.Pool) *DB { // <-- *** CHANGE THIS PARAMETER TYPE IN PRODUCTION ***
	// Optional: In production setup (e.g., in your main or config package),
	// you would connect the real pool and then pass it here:
	// realPool, err := pgxpool.Connect(...)
	// appDB := repository.NewDB(realPool) // This is valid because *pgxpool.Pool implements PoolIface

	return &DB{
		Pool: pool, // Store the interface
		SQ:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Close calls the Close method on the interface
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
