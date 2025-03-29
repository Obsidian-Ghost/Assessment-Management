package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"assessment-management-system/db"
)

// RunMigrations runs all SQL migrations in the migrations directory
func RunMigrations(db *db.DB) error {
	// Find the migrations directory
	_, b, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Dir(b)

	// Read the init.sql file
	migrationFile := filepath.Join(migrationsDir, "init.sql")
	migrationSQL, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute the migration
	ctx := context.Background()
	_, err = db.Pool.Exec(ctx, string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Migration completed successfully")
	return nil
}
