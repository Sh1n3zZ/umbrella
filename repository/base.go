package repository

import (
	"context"

	"github.com/Sh1n3zZ/umbrella/internal/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseRepository provides common database operations
// Following OCP: provides extensible base functionality without requiring modification
type BaseRepository struct {
	db *pgxpool.Pool
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetQueries returns sqlc queries with proper session handling
// Following DRY: centralizes query creation logic
func (r *BaseRepository) GetQueries(ctx context.Context) *sqlc.Queries {
	return GetQueriesWithSession(ctx, r.db)
}

// GetDB returns the database connection
// Following encapsulation: provides controlled access to database
func (r *BaseRepository) GetDB() *pgxpool.Pool {
	return r.db
}
