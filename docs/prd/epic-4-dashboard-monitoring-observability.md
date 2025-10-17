# Epic 4: Dashboard, Monitoring & Observability

**Goal:** Build the adaptive dashboard with storage usage widgets, real-time task monitoring interface, and multi-channel notification system (Email, Teams, Discord). Deliver comprehensive visibility into system operations, backup job status, and proactive error alerting. This epic transforms IronArchive from a functional tool into a production-ready, observable, maintainable system that MSPs can confidently deploy for clients.

## Story 4.1: Dashboard Storage Usage Widgets Backend API

As a developer,
I want API endpoints for dashboard storage statistics,
so that the frontend can display storage usage visualizations.

**Acceptance Criteria:**

1. Endpoint: `GET /api/v1/dashboard/storage-summary` returns overall storage statistics (requires authentication)
2. Overall stats include: total_storage_bytes, total_emails_count, total_tenants_count, total_mailboxes_count
3. Endpoint: `GET /api/v1/dashboard/storage-by-tenant` returns storage breakdown per tenant
4. Per-tenant stats include: tenant_id, tenant_name, storage_bytes, email_count, mailbox_count, last_sync_at
5. Endpoint: `GET /api/v1/tenants/:id/storage-by-mailbox` returns storage breakdown per mailbox for tenant
6. Per-mailbox stats include: mailbox_id, email_address, display_name, storage_bytes, email_count, last_sync_at
7. RBAC: MSP Admin sees all tenants, Tenant Admin sees only assigned tenants
8. Calculations: storage_bytes summed from emails.size_bytes + attachments.size_bytes (deduplicated), email_count from COUNT(emails)
9. Caching: results cached in Redis for 5 minutes to reduce database load
10. Response time: <500ms for overall stats, <1s for per-tenant breakdown
11. Unit tests cover calculation logic, RBAC checks, caching
12. Integration test: verify stats accurate after syncing test emails

## Story 4.2: Dashboard Storage Widgets Frontend UI

As an MSP Admin,
I want storage usage widgets on the dashboard,
so that I can see storage consumption at a glance.

**Acceptance Criteria:**

1. Dashboard storage widgets: Overall Storage Card (total storage, total emails), Storage by Tenant Chart (bar or pie chart), Top 10 Largest Mailboxes List
2. Overall Storage Card displays: total storage (formatted: GB/TB), total emails (formatted: 1.2M), total tenants, total mailboxes
3. Storage by Tenant Chart: bar chart with tenant names on X-axis, storage size on Y-axis, color-coded bars
4. Chart interactive: click bar to drill down to tenant detail page
5. Top 10 Largest Mailboxes List: shows email address, storage size, email count; click to navigate to mailbox detail page
6. Empty state (no data): displays "No tenants added yet" message with "Add Tenant" button
7. Progressive disclosure: storage widgets only appear after first emails archived (hide when total_emails_count = 0)
8. Loading state: skeleton loaders while fetching data
9. Error state: error message with retry button if API fails
10. Auto-refresh: data refreshes every 60 seconds (configurable)
11. Accessibility: charts have accessible labels, keyboard navigation
12. Responsive design: mobile (charts stack vertically), tablet, desktop

## Story 4.3: Task Monitoring API Endpoints

As a developer,
I want API endpoints for job monitoring,
so that the frontend can display backup job status and history.

**Acceptance Criteria:**

1. Endpoint: `GET /api/v1/jobs?status={}&type={}&tenant_id={}&limit={}&offset={}` lists jobs with filters (requires authentication)
2. Status filter options: all, queued, running, completed, failed
3. Type filter options: all, sync_mailbox, sync_tenant, sync_all, export
4. Job list returns: job_id, type, status, tenant_name, mailbox_email (if applicable), started_at, completed_at, duration, emails_synced (if sync job), error_message (if failed)
5. Pagination: default limit=50, max=100, offset for pagination
6. Sorting: default sort by started_at DESC (most recent first)
7. Endpoint: `GET /api/v1/jobs/:id` returns detailed job information
8. Job details include: all list fields + full error_message, metadata JSONB (progress details, email IDs for exports, etc.)
9. Endpoint: `POST /api/v1/jobs/:id/retry` re-enqueues failed job (requires MSP Admin or Tenant Admin for tenant jobs)
10. RBAC: MSP Admin sees all jobs, Tenant Admin sees only assigned tenant jobs, Users see only their export jobs
11. Real-time updates: consider WebSocket endpoint or SSE for live job status updates (optional enhancement)
12. Unit tests cover filtering logic, RBAC checks
13. Integration test: create jobs with various statuses, verify filtering and RBAC enforcement

## Story 4.4: Task Monitoring Widget Frontend UI

As an MSP Admin,
I want a task monitoring widget on the dashboard,
so that I can see backup job status in real-time.

**Acceptance Criteria:**

1. Task monitoring widget on dashboard: displays list of running/recent jobs with status
2. Widget includes: filter dropdown (All/Running/Failed), "View All Jobs" link to dedicated jobs page
3. Job list items display: job type icon, tenant name (or "All Tenants"), mailbox email (if applicable), status badge, progress bar (if running), duration/timestamp, "Retry" button (if failed)
4. Status badges: Queued (gray), Running (blue with spinner), Completed (green), Failed (red)
5. Progress bar: shows percentage if available (e.g., "450/1000 emails synced"), indeterminate spinner if no percentage
6. Auto-refresh: jobs list refreshes every 1 second when any running job present, every 30 seconds otherwise
7. Real-time updates: visual flash/animation when job status changes
8. "Retry" button: calls retry API, shows loading state, updates job status on success
9. Widget shows max 10 most recent jobs, "View All" link for full history
10. Empty state: "No recent backup jobs" message
11. Loading state: skeleton loaders during initial fetch
12. Accessibility: job status changes announced to screen readers, keyboard navigation
13. Responsive design: mobile (compact view), tablet, desktop

## Story 4.5: Detailed Jobs History Page

As an MSP Admin,
I want a dedicated jobs history page,
so that I can review all past backup and export jobs with detailed filtering.

**Acceptance Criteria:**

1. Jobs history page (`/jobs`) displays full jobs list with advanced filtering
2. Filter controls: Status dropdown (All/Queued/Running/Completed/Failed), Type dropdown (All/Sync Mailbox/Sync Tenant/Sync All/Export), Tenant dropdown (All/{tenant names}), Date Range picker
3. Jobs table columns: Type Icon, Tenant, Mailbox (if applicable), Status Badge, Started At, Duration, Emails Processed (for sync), Actions (View Details, Retry if failed)
4. "View Details" button: opens modal/panel with full job information including error messages, metadata, logs
5. Table sorting: click column headers to sort (Started At, Duration, Status)
6. Pagination: 50 jobs per page, numbered pagination at bottom
7. Auto-refresh: refreshes every 5 seconds if any running jobs visible
8. Bulk actions: "Retry All Failed" button (with confirmation) retries all failed jobs matching current filters
9. Export jobs list: CSV export of jobs table for reporting
10. Empty state: "No jobs match your filters" message
11. RBAC: MSP Admin sees all jobs, Tenant Admin sees only assigned tenant jobs
12. Accessibility: table accessible, sortable columns announced, keyboard navigation
13. Responsive design: mobile (cards instead of table), tablet, desktop

## Story 4.6: Notification System Backend Infrastructure

As a developer,
I want notification infrastructure,
so that users can receive alerts via Email, Teams, and Discord.

**Acceptance Criteria:**

1. Notification service module: `NotificationService` with methods `SendEmail()`, `SendTeamsMessage()`, `SendDiscordMessage()`
2. Email notification: uses SMTP client configured via environment variables (SMTP host, port, username, password, from address)
3. Email template: HTML email with header, body content, footer; includes IronArchive branding and links
4. Teams notification: sends Adaptive Card via webhook URL configured in settings
5. Adaptive Card template: includes title, severity color (green/yellow/red), fields (tenant, mailbox, error message), action buttons ("View in IronArchive")
6. Discord notification: sends rich embed via webhook URL configured in settings
7. Discord embed template: includes title, severity color, fields, footer with timestamp, thumbnail icon
8. Notification types: SyncSuccess, SyncFailure, ExportReady, ErrorAlert
9. Notification preferences: stored in database (users table: notification_channels JSONB array ["email", "teams", "discord"], notification_webhook_urls JSONB)
10. Notification queue: uses Redis queue to send notifications asynchronously (don't block main operations)
11. Error handling: if notification fails, log error but don't fail main operation; retry 3 times for transient failures
12. Testing: unit tests with mocked SMTP/HTTP clients, integration test sends real notifications to test accounts
13. Logging: INFO for successful sends, WARN for retries, ERROR for permanent failures

## Story 4.7: Notification Settings Backend API

As a developer,
I want notification settings API endpoints,
so that users can configure their notification preferences.

**Acceptance Criteria:**

1. Endpoint: `GET /api/v1/settings/notifications` returns user's notification settings (requires authentication)
2. Settings return: enabled_channels (array: ["email", "teams", "discord"]), email_address, teams_webhook_url, discord_webhook_url, notification_types (array of types to receive)
3. Endpoint: `PATCH /api/v1/settings/notifications` updates notification settings
4. Update accepts: enabled_channels, email_address (optional, defaults to user's account email), teams_webhook_url, discord_webhook_url, notification_types
5. Validation: webhook URLs must be valid HTTPS URLs, notification_types must be valid enum values
6. Endpoint: `POST /api/v1/settings/notifications/test` sends test notification to verify configuration (requires authentication)
7. Test notification: sends sample message to all enabled channels, returns success/failure per channel
8. RBAC: users can only view/update their own notification settings
9. MSP Admin global settings: endpoint `PATCH /api/v1/settings/global-notifications` sets default notification channels for new users
10. Unit tests cover validation, RBAC checks
11. Integration test: update settings, send test notification, verify received

## Story 4.8: Notification Settings Frontend UI

As a user,
I want to configure my notification preferences,
so that I receive alerts via my preferred channels.

**Acceptance Criteria:**

1. Notification settings page (`/profile/notifications`) displays notification configuration form
2. Form sections: Notification Channels (checkboxes: Email, Microsoft Teams, Discord), Channel Configuration (conditional fields based on enabled channels), Notification Types (checkboxes: Sync Success, Sync Failure, Export Ready, Error Alerts)
3. Email configuration: email address input (pre-filled with user account email, editable)
4. Teams configuration: webhook URL input, "How to get webhook URL" link to documentation
5. Discord configuration: webhook URL input, "How to get webhook URL" link to documentation
6. "Test Notifications" button: sends test notification to all enabled channels, shows success/failure per channel
7. Test results: displays green checkmark or red X with error message per channel
8. "Save Settings" button: updates settings via API, shows success message
9. Form validation: webhook URLs must be valid HTTPS URLs, at least one notification type selected
10. Loading states: during save and test operations
11. Accessibility: form accessible, checkboxes keyboard-navigable, validation errors announced
12. Responsive design: mobile, tablet, desktop

## Story 4.9: Notification Triggers Integration

As a developer,
I want notifications triggered at appropriate events,
so that users are alerted to important system activities.

**Acceptance Criteria:**

1. Sync job completion: trigger SyncSuccess notification (send to MSP Admin and assigned Tenant Admin) if job completed successfully
2. Sync job failure: trigger SyncFailure notification (send to MSP Admin and assigned Tenant Admin) immediately with error details
3. Sync job failure includes: tenant name, mailbox email, error message, timestamp, link to job details page
4. Export job completion: trigger ExportReady notification (send to export job owner) with download link and expiration time
5. Permanent errors: trigger ErrorAlert notification (send to MSP Admin) for critical errors requiring manual intervention (e.g., invalid Azure credentials, permissions revoked)
6. Notification batching: if multiple mailboxes fail in same tenant sync, batch into single notification (list all failed mailboxes)
7. Notification throttling: don't send more than 10 SyncFailure notifications per hour per tenant (prevent spam if misconfigured)
8. Integration: modify sync engine and export processor to call NotificationService on completion/failure
9. Logging: INFO when notification sent, WARN if throttled, ERROR if notification fails
10. Unit tests: verify notifications triggered for various scenarios
11. Integration test: trigger sync failure, verify notification received via email/Teams/Discord

---
