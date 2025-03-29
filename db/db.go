package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"assessment-management-system/config"
)

// DB represents a database connection
type DB struct {
	Pool *pgxpool.Pool
}

// Connect establishes a connection to the database
func Connect(dbConfig config.DatabaseConfig) (*DB, error) {
	var connectionString string

	// Use DATABASE_URL if provided, otherwise build connection string from individual parts
	if dbConfig.URL != "" {
		connectionString = dbConfig.URL
	} else {
		connectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	}

	ctx := context.Background()
	PoolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Optimize pool settings
	PoolConfig.MaxConns = 50
	PoolConfig.MinConns = 10
	PoolConfig.MaxConnLifetime = time.Hour
	PoolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.ConnectConfig(ctx, PoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// ExecuteTransaction executes a function within a transaction
func (db *DB) ExecuteTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // Re-throw panic after rolling back
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
