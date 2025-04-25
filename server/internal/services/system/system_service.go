// File: internal/services/system/system_service.go
package system

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// SystemService defines the interface for system-related operations.
type SystemService interface {
	CheckHealth(ctx context.Context) error
}

// systemService implements the SystemService interface.
type systemService struct {
	db *bun.DB
}

// NewSystemService creates a new SystemService.
// It requires the database connection pool to perform health checks.
func NewSystemService(db *bun.DB) SystemService {
	if db == nil {
		// Handle nil DB appropriately, maybe return an error or a service that always fails health checks
		// For now, let's assume db is always provided correctly during setup.
		fmt.Println("Warning: NewSystemService received a nil DB connection.")
	}
	return &systemService{
		db: db,
	}
}

// CheckHealth attempts to ping the database to verify connectivity.
func (s *systemService) CheckHealth(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	// Ping the database
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil // Health check passed
}
