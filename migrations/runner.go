package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"assessment-management-system/db"
)

// Initialize a migration tracking table
func initMigrationTable(ctx context.Context, db *db.DB) error {
	query := `
                CREATE TABLE IF NOT EXISTS schema_migrations (
                        version VARCHAR(255) PRIMARY KEY,
                        applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
                );
        `
	_, err := db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

// Check if a migration has been applied
func isMigrationApplied(ctx context.Context, db *db.DB, version string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1);`
	err := db.Pool.QueryRow(ctx, query, version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if migration %s is applied: %w", version, err)
	}
	return exists, nil
}

// Mark a migration as applied
func markMigrationApplied(ctx context.Context, db *db.DB, version string) error {
	query := `INSERT INTO schema_migrations (version) VALUES ($1);`
	_, err := db.Pool.Exec(ctx, query, version)
	if err != nil {
		return fmt.Errorf("failed to mark migration %s as applied: %w", version, err)
	}
	return nil
}

// RunMigrations runs all SQL migrations in the migrations directory
func RunMigrations(db *db.DB) error {
	// Find the migrations directory
	_, b, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Dir(b)

	// Execute migrations in order
	ctx := context.Background()

	// Initialize migration tracking table
	if err := initMigrationTable(ctx, db); err != nil {
		return err
	}

	// List of migrations to run in order
	migrations := []string{
		"init.sql",
		"add_refresh_tokens.sql",
	}

	// Execute each migration
	for _, filename := range migrations {
		if err := executeMigrationIfNeeded(ctx, db, migrationsDir, filename); err != nil {
			return err
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// executeMigrationIfNeeded runs a single migration file if it hasn't been applied yet
func executeMigrationIfNeeded(ctx context.Context, db *db.DB, migrationsDir, filename string) error {
	// Get migration version from filename
	version := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Check if migration has already been applied
	applied, err := isMigrationApplied(ctx, db, version)
	if err != nil {
		return err
	}

	if applied {
		log.Printf("Migration %s has already been applied, skipping", filename)
		return nil
	}

	log.Printf("Executing migration: %s", filename)

	// Read the migration file
	migrationFile := filepath.Join(migrationsDir, filename)
	migrationSQL, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	// Execute the migration without a transaction
	// Some migrations might have DDL statements that cannot be run in a transaction
	executeErr := false
	_, err = db.Pool.Exec(ctx, string(migrationSQL))
	if err != nil {
		// Check if error is related to duplicate tables/indexes
		if strings.Contains(err.Error(), "relation") && strings.Contains(err.Error(), "already exists") {
			log.Printf("Objects already exist in migration %s, continuing", filename)
		} else {
			executeErr = true
			log.Printf("Warning: Error executing migration %s: %v", filename, err)
		}
	}

	// Only mark as applied if no execution error or just a duplication error
	if !executeErr {
		// Mark migration as applied
		if err := markMigrationApplied(ctx, db, version); err != nil {
			return err
		}
		log.Printf("Migration %s completed successfully", filename)
	} else {
		return fmt.Errorf("failed to execute migration %s: %w", filename, err)
	}

	return nil
}
