-- ============================================================================
-- Migration Rollback: 000001_initial_schema
-- Description: Drop all tables, indexes, triggers, and extensions
-- Created: 2025-10-17
-- ============================================================================

-- ============================================================================
-- SECTION 1: Drop Triggers
-- ============================================================================

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS audit_log_immutable_trigger ON audit_logs;

-- ============================================================================
-- SECTION 2: Drop Functions
-- ============================================================================

DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS prevent_audit_log_modification();

-- ============================================================================
-- SECTION 3: Drop Tables (in reverse dependency order)
-- ============================================================================

-- Drop tables with no dependencies first
DROP TABLE IF EXISTS settings;

-- Drop tables with dependencies on users
DROP TABLE IF EXISTS audit_logs;

-- Drop jobs table (depends on tenants, mailboxes, users)
DROP TABLE IF EXISTS jobs;

-- Drop attachments (depends on emails)
DROP TABLE IF EXISTS attachments;

-- Drop emails (depends on mailboxes)
DROP TABLE IF EXISTS emails;

-- Drop mailboxes (depends on tenants)
DROP TABLE IF EXISTS mailboxes;

-- Drop users (depends on tenants)
DROP TABLE IF EXISTS users;

-- Drop tenants (no dependencies)
DROP TABLE IF EXISTS tenants;

-- ============================================================================
-- SECTION 4: Drop Extensions
-- ============================================================================

DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "uuid-ossp";

-- ============================================================================
-- Rollback Complete
-- ============================================================================
