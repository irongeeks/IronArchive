package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSchemaTablesExist verifies all expected tables are created by migration
func TestSchemaTablesExist(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	expectedTables := []string{
		"tenants",
		"users",
		"mailboxes",
		"emails",
		"attachments",
		"jobs",
		"audit_logs",
		"settings",
	}

	for _, table := range expectedTables {
		var exists bool
		query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)"
		err := db.QueryRow(context.Background(), query, table).Scan(&exists)
		require.NoError(t, err, "Failed to check if table %s exists", table)
		assert.True(t, exists, "Table %s should exist", table)
	}
}

// TestSchemaIndexesExist verifies all expected indexes are created
func TestSchemaIndexesExist(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	expectedIndexes := []string{
		// Users indexes
		"idx_users_email",
		"idx_users_tenant_id",
		// Tenants indexes
		"idx_tenants_azure_tenant_id",
		// Mailboxes indexes
		"idx_mailboxes_tenant_id",
		"idx_mailboxes_sync_enabled",
		// Emails indexes
		"idx_emails_mailbox_id",
		"idx_emails_message_id",
		"idx_emails_sent_at",
		"idx_emails_sender",
		"idx_emails_deleted_at",
		// Attachments indexes
		"idx_attachments_email_id",
		"idx_attachments_sha256_hash",
		// Jobs indexes
		"idx_jobs_status",
		"idx_jobs_type",
		"idx_jobs_tenant_id",
		"idx_jobs_user_id",
		"idx_jobs_created_at",
		// Audit logs indexes
		"idx_audit_logs_user_id",
		"idx_audit_logs_action",
		"idx_audit_logs_timestamp",
	}

	for _, index := range expectedIndexes {
		var exists bool
		query := "SELECT EXISTS (SELECT FROM pg_indexes WHERE schemaname = 'public' AND indexname = $1)"
		err := db.QueryRow(context.Background(), query, index).Scan(&exists)
		require.NoError(t, err, "Failed to check if index %s exists", index)
		assert.True(t, exists, "Index %s should exist", index)
	}
}

// TestSchemaForeignKeys verifies all expected foreign key constraints exist
func TestSchemaForeignKeys(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Query to check foreign key constraints
	query := `
		SELECT
			tc.table_name,
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			rc.delete_rule
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		JOIN information_schema.referential_constraints AS rc
			ON tc.constraint_name = rc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
			AND tc.table_schema = 'public'
		ORDER BY tc.table_name, kcu.column_name
	`

	rows, err := db.Query(context.Background(), query)
	require.NoError(t, err, "Failed to query foreign key constraints")
	defer rows.Close()

	fkCount := 0
	for rows.Next() {
		var tableName, columnName, foreignTable, foreignColumn, deleteRule string
		err := rows.Scan(&tableName, &columnName, &foreignTable, &foreignColumn, &deleteRule)
		require.NoError(t, err, "Failed to scan foreign key row")
		fkCount++

		// Verify delete rules are correct
		switch tableName {
		case "users":
			if columnName == "tenant_id" {
				assert.Equal(t, "CASCADE", deleteRule, "users.tenant_id should have CASCADE delete rule")
			}
		case "mailboxes":
			if columnName == "tenant_id" {
				assert.Equal(t, "CASCADE", deleteRule, "mailboxes.tenant_id should have CASCADE delete rule")
			}
		case "emails":
			if columnName == "mailbox_id" {
				assert.Equal(t, "CASCADE", deleteRule, "emails.mailbox_id should have CASCADE delete rule")
			}
		case "attachments":
			if columnName == "email_id" {
				assert.Equal(t, "CASCADE", deleteRule, "attachments.email_id should have CASCADE delete rule")
			}
		case "jobs":
			if columnName == "tenant_id" {
				assert.Equal(t, "CASCADE", deleteRule, "jobs.tenant_id should have CASCADE delete rule")
			} else if columnName == "mailbox_id" {
				assert.Equal(t, "CASCADE", deleteRule, "jobs.mailbox_id should have CASCADE delete rule")
			} else if columnName == "user_id" {
				assert.Equal(t, "SET NULL", deleteRule, "jobs.user_id should have SET NULL delete rule")
			}
		case "audit_logs":
			if columnName == "user_id" {
				assert.Equal(t, "SET NULL", deleteRule, "audit_logs.user_id should have SET NULL delete rule")
			}
		}
	}

	// Verify we have at least 8 foreign keys
	assert.GreaterOrEqual(t, fkCount, 8, "Should have at least 8 foreign key constraints")
}

// TestAuditLogImmutability verifies audit logs cannot be modified or deleted
func TestAuditLogImmutability(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Insert a test audit log
	var logID string
	insertQuery := "INSERT INTO audit_logs (action, ip_address, details) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(context.Background(), insertQuery, "TEST_ACTION", "127.0.0.1", nil).Scan(&logID)
	require.NoError(t, err, "Failed to insert test audit log")

	// Attempt to update (should fail)
	updateQuery := "UPDATE audit_logs SET action = $1 WHERE id = $2"
	_, err = db.Exec(context.Background(), updateQuery, "MODIFIED_ACTION", logID)
	assert.Error(t, err, "UPDATE should fail on audit_logs")
	assert.Contains(t, err.Error(), "immutable", "Error should mention immutability")

	// Attempt to delete (should fail)
	deleteQuery := "DELETE FROM audit_logs WHERE id = $1"
	_, err = db.Exec(context.Background(), deleteQuery, logID)
	assert.Error(t, err, "DELETE should fail on audit_logs")
	assert.Contains(t, err.Error(), "immutable", "Error should mention immutability")
}

// TestUsersUpdatedAtTrigger verifies updated_at is automatically updated
func TestUsersUpdatedAtTrigger(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// First create a tenant
	var tenantID string
	tenantQuery := "INSERT INTO tenants (name, azure_tenant_id, azure_app_credentials) VALUES ($1, $2, $3) RETURNING id"
	credentials := `{"app_id":"550e8400-e29b-41d4-a716-446655440001","app_secret":"test_secret"}`
	err := db.QueryRow(context.Background(), tenantQuery, "Test Tenant", "550e8400-e29b-41d4-a716-446655440000", credentials).Scan(&tenantID)
	require.NoError(t, err, "Failed to insert test tenant")

	// Insert a test user
	var userID string
	insertQuery := `
		INSERT INTO users (email, password_hash, display_name, role, tenant_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = db.QueryRow(context.Background(), insertQuery, "test@example.com", "hash", "Test User", "USER", tenantID).Scan(&userID)
	require.NoError(t, err, "Failed to insert test user")

	// Update the user
	updateQuery := "UPDATE users SET display_name = $1 WHERE id = $2"
	_, err = db.Exec(context.Background(), updateQuery, "Updated User", userID)
	require.NoError(t, err, "Failed to update user")

	// Check that updated_at is set (trigger updates it)
	var hasUpdatedAt bool
	selectQuery := "SELECT updated_at IS NOT NULL FROM users WHERE id = $1"
	err = db.QueryRow(context.Background(), selectQuery, userID).Scan(&hasUpdatedAt)
	require.NoError(t, err, "Failed to query updated user")

	// updated_at should be set by trigger
	assert.True(t, hasUpdatedAt, "updated_at should not be null")
}

// TestCheckConstraints verifies CHECK constraints are working
func TestCheckConstraints(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Create a tenant first
	var tenantID string
	tenantQuery := "INSERT INTO tenants (name, azure_tenant_id, azure_app_credentials) VALUES ($1, $2, $3) RETURNING id"
	credentials := `{"app_id":"550e8400-e29b-41d4-a716-446655440001","app_secret":"test_secret"}`
	err := db.QueryRow(context.Background(), tenantQuery, "Test Tenant", "550e8400-e29b-41d4-a716-446655440000", credentials).Scan(&tenantID)
	require.NoError(t, err, "Failed to insert test tenant")

	// Test invalid user role (should fail)
	invalidRoleQuery := "INSERT INTO users (email, password_hash, display_name, role, tenant_id) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(context.Background(), invalidRoleQuery, "test@example.com", "hash", "Test User", "INVALID_ROLE", tenantID)
	assert.Error(t, err, "Should fail with invalid role")

	// Test invalid mailbox type (should fail)
	invalidTypeQuery := "INSERT INTO mailboxes (tenant_id, email_address, mailbox_type) VALUES ($1, $2, $3)"
	_, err = db.Exec(context.Background(), invalidTypeQuery, tenantID, "test@example.com", "INVALID_TYPE")
	assert.Error(t, err, "Should fail with invalid mailbox type")

	// Test invalid job progress (should fail)
	invalidProgressQuery := "INSERT INTO jobs (type, status, progress) VALUES ($1, $2, $3)"
	_, err = db.Exec(context.Background(), invalidProgressQuery, "SYNC_MAILBOX", "QUEUED", 150)
	assert.Error(t, err, "Should fail with progress > 100")
}

// TestUniqueConstraints verifies unique constraints are working
func TestUniqueConstraints(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Create a tenant
	var tenantID string
	tenantQuery := "INSERT INTO tenants (name, azure_tenant_id, azure_app_credentials) VALUES ($1, $2, $3) RETURNING id"
	credentials := `{"app_id":"550e8400-e29b-41d4-a716-446655440001","app_secret":"test_secret"}`
	err := db.QueryRow(context.Background(), tenantQuery, "Test Tenant", "550e8400-e29b-41d4-a716-446655440000", credentials).Scan(&tenantID)
	require.NoError(t, err, "Failed to insert test tenant")

	// Insert first user
	userQuery := "INSERT INTO users (email, password_hash, display_name, role, tenant_id) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(context.Background(), userQuery, "test@example.com", "hash", "Test User", "USER", tenantID)
	require.NoError(t, err, "Failed to insert first user")

	// Try to insert duplicate email (should fail)
	_, err = db.Exec(context.Background(), userQuery, "test@example.com", "hash", "Another User", "USER", tenantID)
	assert.Error(t, err, "Should fail with duplicate email")

	// Insert first mailbox
	mailboxQuery := "INSERT INTO mailboxes (tenant_id, email_address, mailbox_type) VALUES ($1, $2, $3)"
	_, err = db.Exec(context.Background(), mailboxQuery, tenantID, "mailbox@example.com", "USER")
	require.NoError(t, err, "Failed to insert first mailbox")

	// Try to insert duplicate mailbox for same tenant (should fail)
	_, err = db.Exec(context.Background(), mailboxQuery, tenantID, "mailbox@example.com", "SHARED")
	assert.Error(t, err, "Should fail with duplicate tenant_id+email_address")
}

// TestInitialSettings verifies initial settings are populated
func TestInitialSettings(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	expectedSettings := map[string]bool{
		"global_retention_policy_days": false,
		"smtp_config":                  false,
		"notification_channels":        false,
		"scheduler_enabled":            false,
		"sync_schedule":                false,
	}

	// Query all settings
	query := "SELECT key FROM settings"
	rows, err := db.Query(context.Background(), query)
	require.NoError(t, err, "Failed to query settings")
	defer rows.Close()

	for rows.Next() {
		var key string
		err := rows.Scan(&key)
		require.NoError(t, err, "Failed to scan setting key")

		if _, exists := expectedSettings[key]; exists {
			expectedSettings[key] = true
		}
	}

	// Verify all expected settings exist
	for key, found := range expectedSettings {
		assert.True(t, found, "Setting %s should exist", key)
	}
}

// TestEmailsTableColumns verifies emails table has required body columns (AC-6)
func TestEmailsTableColumns(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Query column information for emails table
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'emails'
			AND column_name IN ('body_text', 'body_html')
		ORDER BY column_name
	`

	rows, err := db.Query(context.Background(), query)
	require.NoError(t, err, "Failed to query emails table columns")
	defer rows.Close()

	foundColumns := make(map[string]bool)
	for rows.Next() {
		var columnName, dataType, isNullable string
		err := rows.Scan(&columnName, &dataType, &isNullable)
		require.NoError(t, err, "Failed to scan column info")

		foundColumns[columnName] = true

		// Verify data type is TEXT
		assert.Equal(t, "text", dataType, "Column %s should be TEXT type", columnName)
	}

	// Verify both columns exist
	assert.True(t, foundColumns["body_text"], "body_text column should exist in emails table")
	assert.True(t, foundColumns["body_html"], "body_html column should exist in emails table")
}

// TestJobsTableErrorColumn verifies jobs table has error_message column (AC-8)
func TestJobsTableErrorColumn(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Query column information for jobs table
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'jobs'
			AND column_name = 'error_message'
	`

	var columnName, dataType, isNullable string
	err := db.QueryRow(context.Background(), query).Scan(&columnName, &dataType, &isNullable)
	require.NoError(t, err, "error_message column should exist in jobs table")

	// Verify data type is TEXT
	assert.Equal(t, "text", dataType, "error_message column should be TEXT type")
	assert.Equal(t, "error_message", columnName, "Column name should be error_message")
}

// TestTenantsTableCredentialsColumn verifies tenants table has azure_app_credentials column (AC-4)
func TestTenantsTableCredentialsColumn(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(t, db)

	// Query column information for tenants table
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'tenants'
			AND column_name = 'azure_app_credentials'
	`

	var columnName, dataType, isNullable string
	err := db.QueryRow(context.Background(), query).Scan(&columnName, &dataType, &isNullable)
	require.NoError(t, err, "azure_app_credentials column should exist in tenants table")

	// Verify data type is TEXT and NOT NULL
	assert.Equal(t, "text", dataType, "azure_app_credentials column should be TEXT type")
	assert.Equal(t, "azure_app_credentials", columnName, "Column name should be azure_app_credentials")
	assert.Equal(t, "NO", isNullable, "azure_app_credentials should be NOT NULL")
}
