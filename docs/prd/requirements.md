# Requirements

## Functional

**FR1:** The system shall provide a setup wizard that creates the initial MSP Admin account with secure password requirements during first-time installation

**FR2:** The system shall implement JWT-based authentication with bcrypt password hashing (cost factor 12+), session management, and secure logout

**FR3:** The system shall enforce Role-Based Access Control (RBAC) with three distinct roles: MSP Admin (full system access), Tenant Admin (single tenant access), and User (read-only own mailbox access)

**FR4:** The system shall support Multi-Factor Authentication (MFA) via TOTP and/or Microsoft 365 SSO integration

**FR5:** The system shall provide self-service password reset functionality via email without requiring MSP intervention

**FR6:** The system shall provide a tenant onboarding wizard that automates Azure AD app creation, mailbox discovery via Graph API, visual mailbox selection interface, and immediate backup initiation

**FR7:** The system shall automatically sync all selected mailboxes 4 times daily (6 AM, 12 PM, 6 PM, 12 AM local time) using incremental delta queries

**FR8:** The system shall support manual backup triggers at dashboard-level, tenant-level, and mailbox-level

**FR9:** The system shall display real-time sync progress with live progress bars, estimated time remaining, email count, and auto-refresh (1 second interval)

**FR10:** The system shall provide instant-as-you-type search powered by Meilisearch with results returned in < 200ms, respecting RBAC permissions

**FR11:** The system shall support advanced search filters including date range, sender, recipient, attachment presence, and size range

**FR12:** The system shall display search results with highlighted terms, infinite scroll (50 emails per batch), and sorting options (date, size, sender)

**FR13:** The system shall provide a side panel detail view showing full email body, headers, and attachments without leaving search results

**FR14:** The system shall support multi-select export from search results with checkboxes, "Select All" functionality, and batch operations

**FR15:** The system shall generate EML/ZIP exports for 1-100 emails and PST exports for 100+ emails or full mailboxes

**FR16:** The system shall provide an export modal showing format selection, estimated file size, and generation time

**FR17:** The system shall handle small exports (< 100MB) with immediate download and large exports via background job queue with notification on completion

**FR18:** The system shall generate download links with 7-day expiration and automatic cleanup of expired files

**FR19:** The system shall provide an adaptive dashboard with empty state design ("Add Your First Tenant") that progressively reveals complexity as data populates

**FR20:** The system shall display storage usage widgets with visual charts showing per-tenant breakdown and drill-down navigation (Dashboard → Tenant → Mailbox → Email list)

**FR21:** The system shall provide a task monitoring widget with live-updating list of running/past backup jobs, auto-refresh (1 second), color-coded status, and expandable task details

**FR22:** The system shall support task filtering by status, type (scheduled/manual), tenant, and date range with "Retry" button for failed tasks

**FR23:** The system shall provide global settings (MSP Admin only) for whitelabeling (custom logo, favicon, sanitized CSS), user management, retention policies, and notification channels

**FR24:** The system shall provide profile settings (all users) for theme selection with instant switching (Catppuccin Mocha/Latte, Nord, Cyberpunk, Dracula, Tokyo Night), password change, MFA setup, and display name update

**FR25:** The system shall implement a three-tier retention policy override system (Global → Tenant → Mailbox) with pre-configured DSGVO/GoBD templates (6yr, 8yr, 10yr)

**FR26:** The system shall provide Legal Hold functionality as a simple on/off toggle at tenant or mailbox level for infinite retention

**FR27:** The system shall maintain comprehensive audit logs capturing all user actions (login, search, export, configuration changes) with timestamp, user, IP address, and action details

**FR28:** The system shall implement auto-retry logic for transient errors (rate limits, network timeouts) with 3 retries and exponential backoff

**FR29:** The system shall handle permanent errors (mailbox deleted, permission denied) without retry and create high-priority alerts requiring MSP Admin acknowledgment

**FR30:** The system shall send notifications via multiple channels (Email with HTML tables, Teams with Adaptive Cards, Discord with rich embeds) with channel-specific formatting

**FR31:** The system shall deduplicate attachments using SHA-256 hashing to optimize storage usage

**FR32:** The system shall store emails in a structured filesystem hierarchy: `/archive/tenants/{uuid}/mailboxes/{uuid}/emails/{year}/{month}/`

## Non Functional

**NFR1:** Search queries shall return results in < 200ms for 95th percentile of queries on datasets up to 10M emails

**NFR2:** Dashboard page load time shall be < 2 seconds on modern broadband connections

**NFR3:** The system shall support 50+ concurrent users on recommended hardware specifications

**NFR4:** Email sync throughput shall process 10,000+ emails per hour per tenant during initial backups

**NFR5:** The system shall provide multi-architecture Docker images supporting both ARM64 and AMD64 platforms

**NFR6:** Database schema migrations shall execute automatically on application startup

**NFR7:** All services shall include Docker health checks with automatic container restart on failure

**NFR8:** The system shall target WCAG 2.1 AA accessibility compliance

**NFR9:** The user interface shall be responsive from 320px (mobile) to 4K displays

**NFR10:** The system shall support modern desktop browsers (Chrome/Edge 120+, Firefox 120+, Safari 17+) and mobile browsers (iOS Safari 16+, Android Chrome 120+)

**NFR11:** All external connections shall use TLS 1.3 encryption (Graph API, webhooks, user browsers)

**NFR12:** Password hashing shall use bcrypt with cost factor 12 or higher

**NFR13:** JWT tokens shall use 15-minute access token expiration and 7-day refresh token expiration

**NFR14:** Multi-tenant data isolation shall be enforced via Row-Level Security (RLS) in PostgreSQL or application-enforced tenant filtering on all queries

**NFR15:** Audit log tables shall be immutable (append-only) with database triggers preventing deletion or modification

**NFR16:** All settings shall be configurable via environment variables or `.env` file for deployment flexibility

**NFR17:** The system shall handle Microsoft Graph API rate limits gracefully with backoff and retry strategies

**NFR18:** Backup operations shall ensure zero data loss (all emails retrieved via Graph API are stored)

**NFR19:** The system shall enforce proper RBAC permission boundaries preventing cross-tenant data access

**NFR20:** Email retention policies shall automatically delete emails after configured period expires
