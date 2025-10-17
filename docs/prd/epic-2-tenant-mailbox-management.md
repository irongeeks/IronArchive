# Epic 2: Tenant & Mailbox Management

**Goal:** Build the tenant onboarding wizard, integrate with Microsoft Graph API to discover mailboxes, implement the email sync engine with incremental delta queries, and establish the scheduling system for automatic 4x daily backups. Deliver the ability to add M365 tenants, select mailboxes, and successfully archive emails with filesystem storage and Meilisearch indexing. This epic creates the core archiving functionality that is the product's raison d'être.

## Story 2.1: Tenant CRUD API Endpoints

As a developer,
I want tenant management API endpoints,
so that MSP Admins can create, read, update, and delete tenants.

**Acceptance Criteria:**

1. Endpoint: `POST /api/v1/tenants` creates new tenant (requires MSP Admin role)
2. Create tenant accepts: name, azure_tenant_id, azure_app_id, azure_app_secret
3. Create tenant encrypts azure_app_secret before storage
4. Create tenant returns tenant object with UUID, created_at, validation errors if input invalid
5. Endpoint: `GET /api/v1/tenants` lists all tenants (MSP Admin sees all, Tenant Admin sees only assigned tenants)
6. List returns: id, name, mailbox_count, last_sync_at, storage_size_bytes, retention_policy_days, legal_hold
7. Endpoint: `GET /api/v1/tenants/:id` returns single tenant details with mailbox list summary
8. Endpoint: `PATCH /api/v1/tenants/:id` updates tenant name, retention_policy_days, legal_hold (requires MSP Admin or assigned Tenant Admin)
9. Endpoint: `DELETE /api/v1/tenants/:id` deletes tenant and cascades to mailboxes, emails, attachments (requires MSP Admin, confirmation required)
10. RBAC enforced: Tenant Admins can only access assigned tenants, Users cannot access tenant management
11. Unit tests cover validation, encryption/decryption, RBAC checks
12. Integration tests cover CRUD operations, authorization failures

## Story 2.2: Microsoft Graph API Client Wrapper

As a developer,
I want a Graph API client wrapper,
so that I can interact with Microsoft 365 APIs consistently with proper error handling and rate limiting.

**Acceptance Criteria:**

1. Graph API client initialized with OAuth2 client credentials flow
2. Token acquisition method: exchanges tenant azure_app_id + azure_app_secret for access token
3. Token caching implemented: stores token in memory with expiration tracking, refreshes automatically before expiration
4. HTTP client configured with: timeouts (30s), retry logic (3 attempts with exponential backoff), user agent header
5. Rate limit handling: detects 429 status code, reads Retry-After header, waits and retries
6. Error handling: distinguishes transient errors (network, rate limit) from permanent errors (invalid credentials, mailbox not found)
7. Methods implemented: `ListMailboxes(tenantID)`, `GetMailboxDeltaQuery(mailboxID, deltaToken)`, `GetEmailMessage(messageID)`, `GetAttachment(messageID, attachmentID)`
8. Logging: INFO for successful calls, WARN for retries, ERROR for permanent failures with request/response details
9. Unit tests with mocked HTTP responses cover: successful token acquisition, token refresh, rate limit handling, error scenarios
10. Integration test with real M365 test tenant validates end-to-end connectivity (manual test, not CI)

## Story 2.3: Mailbox Discovery and Selection Backend

As a developer,
I want mailbox discovery API endpoints,
so that tenants' mailboxes can be discovered and selected for backup.

**Acceptance Criteria:**

1. Endpoint: `POST /api/v1/tenants/:id/discover-mailboxes` triggers mailbox discovery via Graph API (requires MSP Admin)
2. Discovery fetches all mailboxes from M365 tenant via Graph API `users` endpoint
3. Discovery stores discovered mailboxes in `mailboxes` table with: email_address, display_name, mailbox_type (User/Shared/Room/Equipment), last_activity_date, size_bytes, sync_enabled=false (default)
4. Discovery updates existing mailboxes if already present (upsert logic based on email_address)
5. Discovery returns list of discovered mailboxes with metadata: id, email, display_name, type, size, last_activity, sync_enabled
6. Endpoint: `GET /api/v1/tenants/:id/mailboxes` returns list of mailboxes for tenant
7. Endpoint: `PATCH /api/v1/mailboxes/:id` updates mailbox sync_enabled flag (requires MSP Admin or Tenant Admin)
8. Endpoint: `PATCH /api/v1/tenants/:id/mailboxes/bulk-enable` accepts array of mailbox IDs, sets sync_enabled=true (requires MSP Admin or Tenant Admin)
9. Error handling: if Graph API fails, return error with actionable message ("Invalid credentials", "Tenant not found")
10. RBAC enforced: only MSP Admin or assigned Tenant Admin can discover/modify mailboxes
11. Unit tests cover upsert logic, RBAC checks
12. Integration tests cover discovery flow, bulk enable

## Story 2.4: Email Sync Engine (Initial Sync)

As a developer,
I want the email sync engine to perform initial full sync for a mailbox,
so that all existing emails are archived.

**Acceptance Criteria:**

1. Sync function: `SyncMailbox(mailboxID)` performs initial full sync using Graph API delta query
2. Initial sync: fetches all messages from mailbox using `GET /users/{id}/messages/delta` with no delta token
3. For each message: extracts metadata (message_id, subject, sender, recipients, sent_at, body_text, body_html, has_attachments, size_bytes)
4. For each message: saves email record to `emails` table
5. For each attachment: downloads attachment content, calculates SHA-256 hash for deduplication
6. Attachment deduplication: if SHA-256 hash exists in `attachments` table, reuse existing file_path; otherwise save to filesystem
7. Filesystem storage: save email body to `/archive/tenants/{tenant_uuid}/mailboxes/{mailbox_uuid}/emails/{year}/{month}/{message_id}.json`
8. Filesystem storage: save attachments to `/archive/tenants/{tenant_uuid}/attachments/{sha256_hash}.{extension}`
9. Meilisearch indexing: after email saved, index searchable fields (subject, sender, recipients, body_text snippet) in Meilisearch
10. Delta token: after sync completes, save delta token from response to `mailboxes.last_delta_token`
11. Sync progress tracking: update `jobs` table with progress (emails_synced count, status: running/completed/failed)
12. Error handling: transient errors retry, permanent errors mark job as failed with error message
13. Transaction safety: database writes use transactions, rollback on failure
14. Logging: INFO for sync start/completion, DEBUG for each email processed, ERROR for failures
15. Unit tests cover: metadata extraction, deduplication logic, error handling
16. Integration test: sync small mailbox (10 emails) successfully, verify database and filesystem state

## Story 2.5: Email Sync Engine (Incremental Delta Sync)

As a developer,
I want the email sync engine to perform incremental delta syncs,
so that only new/changed emails are archived on subsequent syncs.

**Acceptance Criteria:**

1. Incremental sync function: `SyncMailboxIncremental(mailboxID)` uses saved delta token from previous sync
2. Incremental sync: fetches only new/changed messages using `GET /users/{id}/messages/delta?$deltatoken={token}`
3. For new messages: same processing as initial sync (save email, attachments, index in Meilisearch)
4. For changed messages (e.g., flag updates): update email record metadata if changed
5. For deleted messages: mark email as deleted in database (soft delete, do not physically remove yet—retention policy handles this)
6. Delta token update: save new delta token after successful incremental sync
7. Sync efficiency: only processes delta results, not entire mailbox (verify via logging: "Processed 5 new emails" vs. "Processed 50,000 emails")
8. Error handling: if delta token expired (>30 days), fall back to initial sync
9. Concurrency safety: prevent multiple simultaneous syncs for same mailbox (use database lock or Redis lock)
10. Logging: INFO for incremental sync start/completion with email counts (new, changed, deleted)
11. Unit tests cover: delta token handling, soft delete logic, fallback to initial sync
12. Integration test: perform initial sync, add new email to M365 mailbox, perform incremental sync, verify only new email processed

## Story 2.6: Manual Backup Trigger API and Job Queue

As a developer,
I want manual backup trigger endpoints and job queue infrastructure,
so that backups can be initiated on-demand and processed asynchronously.

**Acceptance Criteria:**

1. Redis-based job queue configured using asynq library with workers
2. Job types defined: `SyncMailboxJob`, `SyncTenantJob`, `SyncDashboardJob` (all mailboxes all tenants)
3. Endpoint: `POST /api/v1/mailboxes/:id/sync` enqueues `SyncMailboxJob` for specific mailbox (requires MSP Admin or Tenant Admin)
4. Endpoint: `POST /api/v1/tenants/:id/sync` enqueues `SyncTenantJob` for all enabled mailboxes in tenant (requires MSP Admin or Tenant Admin)
5. Endpoint: `POST /api/v1/sync-all` enqueues `SyncDashboardJob` for all tenants (requires MSP Admin only)
6. Job queue worker: dequeues jobs, executes sync engine, updates `jobs` table with status/progress
7. Job record created on enqueue: type, status=queued, tenant_id, mailbox_id (if applicable), created_at
8. Job record updated during execution: status=running, started_at
9. Job record updated on completion: status=completed/failed, completed_at, error_message (if failed)
10. Endpoints return job ID immediately (202 Accepted) allowing frontend to poll for status
11. Job status endpoint: `GET /api/v1/jobs/:id` returns job details (type, status, progress, error_message)
12. Concurrency limits: max 5 concurrent sync jobs to prevent resource exhaustion
13. RBAC enforced: Tenant Admins can only trigger syncs for assigned tenants
14. Unit tests cover job creation, RBAC checks
15. Integration test: enqueue job, verify worker processes it, check job status transitions (queued → running → completed)

## Story 2.7: Scheduled Automatic Sync (4x Daily)

As a developer,
I want automatic sync scheduling,
so that backups run 4 times daily without manual intervention.

**Acceptance Criteria:**

1. Scheduler configured using pg_cron extension or Go cron library (e.g., robfig/cron)
2. Cron jobs scheduled: 6 AM, 12 PM, 6 PM, 12 AM (system local time or configurable via environment variable)
3. Each cron execution: enqueues `SyncDashboardJob` (syncs all enabled mailboxes across all tenants)
4. Scheduler logs: INFO when cron job triggered, number of mailboxes queued for sync
5. Scheduler resilience: if scheduler process crashes, jobs resume on restart without duplication
6. Scheduler configuration: cron schedule configurable via environment variable (default: "0 6,12,18,0 * * *")
7. Manual sync pause: global setting to disable automatic syncs (for maintenance), accessible via API
8. Endpoint: `GET /api/v1/settings/scheduler` returns scheduler status (enabled/disabled, next run times)
9. Endpoint: `PATCH /api/v1/settings/scheduler` toggles scheduler enabled flag (requires MSP Admin)
10. Integration test: mock cron execution, verify jobs enqueued for all enabled mailboxes

## Story 2.8: Tenant Onboarding Wizard Frontend (Step 1: Basic Info)

As an MSP Admin,
I want to add a new tenant via a guided wizard,
so that I can onboard M365 tenants easily.

**Acceptance Criteria:**

1. Tenant onboarding wizard accessible from dashboard "Add Tenant" button
2. Wizard UI: multi-step wizard component (Step 1: Basic Info, Step 2: Azure AD Setup, Step 3: Mailbox Selection)
3. Step 1 form includes: Tenant Name (text input), Azure Tenant ID (text input with validation), "Next" button
4. Form validation: tenant name non-empty, Azure Tenant ID is valid GUID format
5. "Next" button disabled until validation passes
6. Step progress indicator shows current step (1 of 3)
7. Wizard state persisted in Svelte store (survives page refresh if user navigates away)
8. Accessibility: keyboard navigation, form labels, ARIA attributes
9. Responsive design: mobile, tablet, desktop
10. Visual design: consistent with application theme, clear instructional text

## Story 2.9: Tenant Onboarding Wizard Frontend (Step 2: Azure AD Setup)

As an MSP Admin,
I want the wizard to guide me through Azure AD app creation,
so that I can authorize IronArchive to access M365 data.

**Acceptance Criteria:**

1. Step 2 displays: Azure AD app creation instructions with copy-paste commands for Azure CLI
2. Instructions include: required API permissions (Mail.ReadWrite, Mail.Read, User.Read.All), redirect URL placeholder
3. Form includes: Azure App ID (text input), Azure App Secret (password input with show/hide toggle), "Verify Connection" button
4. "Verify Connection" button: calls backend API to test Graph API connectivity with provided credentials
5. Backend verification endpoint: `POST /api/v1/tenants/verify-credentials` accepts azure_tenant_id, azure_app_id, azure_app_secret; returns success/failure
6. Verification success: shows green checkmark, enables "Next" button
7. Verification failure: shows error message ("Invalid credentials", "Insufficient permissions"), disables "Next" button
8. Loading state during verification API call
9. "Back" button returns to Step 1, preserves form data
10. Step progress indicator shows current step (2 of 3)

## Story 2.10: Tenant Onboarding Wizard Frontend (Step 3: Mailbox Selection)

As an MSP Admin,
I want to discover and select mailboxes for backup,
so that I can choose which mailboxes to archive.

**Acceptance Criteria:**

1. Step 3 automatically triggers mailbox discovery API call on load
2. Loading state: shows spinner with "Discovering mailboxes..." message
3. Discovered mailboxes displayed in table/list with columns: Checkbox, Email Address, Display Name, Type (badge: User/Shared/Room/Equipment), Size, Last Activity
4. "Select All" checkbox selects/deselects all mailboxes
5. Individual checkboxes for each mailbox
6. Mailboxes sorted by: Last Activity (most recent first) by default
7. Filter controls: mailbox type filter (All, User, Shared, Room, Equipment), search input (filters by email/display name)
8. Badge colors: User (blue), Shared (green), Room (yellow), Equipment (gray)
9. "Finish" button creates tenant and enables selected mailboxes via API (`POST /api/v1/tenants`, then `PATCH /api/v1/tenants/:id/mailboxes/bulk-enable`)
10. Success: shows success message, redirects to tenant detail page after 2 seconds
11. Error: shows error message, allows retry
12. "Back" button returns to Step 2
13. Step progress indicator shows current step (3 of 3)
14. Accessibility: table accessible to screen readers, checkbox states announced
15. Responsive design: table scrolls horizontally on mobile, remains usable

---
