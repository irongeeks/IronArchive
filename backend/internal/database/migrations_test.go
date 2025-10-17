package database

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMigrationUpSuccess tests that migration up completes without errors
func TestMigrationUpSuccess(t *testing.T) {
	// This test verifies the migration runs successfully
	// The setupTestDatabase function already runs migrations
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// If we got here, migration succeeded
	// Verify by checking if at least one table exists
	var exists bool
	query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')"
	err := db.QueryRow(context.Background(), query).Scan(&exists)
	require.NoError(t, err, "Failed to check if users table exists")
	assert.True(t, exists, "Users table should exist after migration")
}

// TestMigrationIdempotency tests that running migrations multiple times is safe
func TestMigrationIdempotency(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}

	// First migration (done by setupTestDatabase)
	db := setupTestDatabase(t)

	// Check initial state
	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM settings").Scan(&count)
	require.NoError(t, err, "Failed to count settings")
	initialCount := count

	// Try to run migration again (should be idempotent via migrate tool)
	// Note: golang-migrate tracks versions, so re-running is safe
	// This test verifies the schema itself is idempotent

	// Clean up
	teardownTestDatabase(t, db)

	// Second migration
	db2 := setupTestDatabase(t)
	defer teardownTestDatabase(t, db2)

	// Check that settings count is the same
	err = db2.QueryRow(context.Background(), "SELECT COUNT(*) FROM settings").Scan(&count)
	require.NoError(t, err, "Failed to count settings after second migration")
	assert.Equal(t, initialCount, count, "Settings count should be the same after re-migration")
}

// TestMigrationDownSuccess tests that rollback completes without errors
func TestMigrationDownSuccess(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}

	// Set up database with migration
	db := setupTestDatabase(t)

	// Verify tables exist
	var exists bool
	query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')"
	err := db.QueryRow(context.Background(), query).Scan(&exists)
	require.NoError(t, err, "Failed to check if users table exists")
	require.True(t, exists, "Users table should exist before rollback")

	// Run migration down (done by teardownTestDatabase)
	teardownTestDatabase(t, db)

	// Reconnect and verify tables are gone
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://ironarchive:ironarchive_password@localhost:5432/ironarchive?sslmode=disable"
	}

	dbAfter, err := ConnectPostgres(dbURL)
	require.NoError(t, err, "Failed to reconnect to database")
	defer dbAfter.Close()

	// Check if users table still exists (should not)
	err = dbAfter.QueryRow(context.Background(), query).Scan(&exists)
	require.NoError(t, err, "Failed to check if users table exists after rollback")
	assert.False(t, exists, "Users table should not exist after rollback")
}

// TestMigrationRepeatability tests that migrations can be applied, rolled back, and applied again
func TestMigrationRepeatability(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}

	// First migration cycle
	db1 := setupTestDatabase(t)
	var count1 int
	err := db1.QueryRow(context.Background(), "SELECT COUNT(*) FROM settings").Scan(&count1)
	require.NoError(t, err, "Failed to count settings in first cycle")
	teardownTestDatabase(t, db1)

	// Second migration cycle
	db2 := setupTestDatabase(t)
	defer teardownTestDatabase(t, db2)

	var count2 int
	err = db2.QueryRow(context.Background(), "SELECT COUNT(*) FROM settings").Scan(&count2)
	require.NoError(t, err, "Failed to count settings in second cycle")

	// Both cycles should have the same initial data
	assert.Equal(t, count1, count2, "Settings count should be consistent across migration cycles")
}
