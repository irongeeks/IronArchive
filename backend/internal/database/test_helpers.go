package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestDatabase creates a test database connection and runs migrations
func setupTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://ironarchive:ironarchive_password@localhost:5432/ironarchive?sslmode=disable"
	}

	// Run migrations up
	err := runMigrations(dbURL, "up")
	if err != nil {
		t.Fatalf("Failed to run migrations up: %v", err)
	}

	// Connect to database
	db, err := ConnectPostgres(dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// teardownTestDatabase rolls back migrations and closes the database connection
func teardownTestDatabase(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	// Close database connection
	if db != nil {
		db.Close()
	}

	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://ironarchive:ironarchive_password@localhost:5432/ironarchive?sslmode=disable"
	}

	// Run migrations down
	err := runMigrations(dbURL, "down")
	if err != nil {
		t.Logf("Warning: Failed to run migrations down: %v", err)
		// Don't fail the test on teardown errors, just log them
	}
}

// runMigrations executes database migrations using the migrate CLI tool
func runMigrations(dbURL, direction string) error {
	// Find the migrations directory relative to the project root
	// When running tests, we're in backend/internal/database/
	// Migrations are at project root /migrations/
	migrationsPath := "../../../migrations"

	var cmd *exec.Cmd
	if direction == "up" {
		cmd = exec.Command("migrate", "-path", migrationsPath, "-database", dbURL, "up")
	} else if direction == "down" {
		cmd = exec.Command("migrate", "-path", migrationsPath, "-database", dbURL, "down", "-all")
	} else {
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("migration %s failed: %v\nOutput: %s", direction, err, string(output))
	}

	return nil
}
