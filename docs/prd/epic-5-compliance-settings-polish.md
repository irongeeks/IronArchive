# Epic 5: Compliance, Settings & Polish

**Goal:** Implement retention policies with automatic deletion, legal hold functionality, comprehensive audit logging, whitelabeling capabilities, theme system, and user profile settings. Deliver the compliance features required for DSGVO/GoBD adherence and the customization options that differentiate IronArchive from competitors. This epic completes the v1.0 feature set, addresses all non-functional requirements, and prepares the application for public release with documentation and polish.

## Story 5.1: Retention Policy Engine Backend

As a developer,
I want retention policy enforcement,
so that emails are automatically deleted after the configured retention period expires.

**Acceptance Criteria:**

1. Retention policy configuration: three-tier override system (Global → Tenant → Mailbox)
2. Global retention policy: stored in settings table (global_retention_policy_days), default templates (6yr=2190 days, 8yr=2920 days, 10yr=3650 days)
3. Tenant retention policy: stored in tenants table (retention_policy_days), overrides global if set
4. Mailbox retention policy: stored in mailboxes table (retention_policy_days), overrides tenant/global if set
5. Retention calculation: email retention_expiry_date = email.sent_at + applicable_retention_policy_days
6. Retention cleanup job: scheduled daily (cron: 2 AM), finds all emails where retention_expiry_date < NOW() AND legal_hold = false
7. Deletion process: soft delete emails (set deleted_at timestamp), remove from Meilisearch index, schedule physical file deletion for 30 days later (grace period)
8. Legal hold check: emails in mailboxes/tenants with legal_hold=true are NEVER deleted regardless of retention policy
9. Audit logging: log each email deletion with email_id, mailbox_id, tenant_id, retention_policy_applied, deleted_by=system
10. Orphaned attachment cleanup: after email deletion, check if attachment SHA-256 hash is referenced by any other email; if not, delete attachment file
11. Job logging: INFO for cleanup job start/completion with counts (emails deleted, attachments cleaned up), WARN if legal hold prevented deletion
12. Configuration endpoint: `GET /api/v1/settings/retention` returns retention policy settings
13. Update endpoint: `PATCH /api/v1/settings/retention` updates global retention policy (requires MSP Admin)
14. Unit tests cover: retention calculation, legal hold enforcement, orphaned attachment detection
15. Integration test: insert old email beyond retention, run cleanup job, verify email soft deleted and removed from search index

## Story 5.2: Retention Policy Settings UI

As an MSP Admin,
I want to configure retention policies,
so that I can comply with DSGVO/GoBD requirements.

**Acceptance Criteria:**

1. Retention policy settings page (`/settings/retention`) displays retention configuration form
2. Global retention policy section: dropdown with templates (6 years, 8 years, 10 years, Custom), custom input (days) if Custom selected
3. Template descriptions: "6 years (GoBD standard for most businesses)", "8 years (GoBD for specific industries)", "10 years (maximum compliance retention)"
4. Warning message: "Shortening retention period will delete older emails during next cleanup (2 AM daily). This action cannot be undone."
5. Per-tenant override: tenant detail page includes "Retention Policy Override" section with same dropdown/custom input
6. Per-mailbox override: mailbox detail page includes "Retention Policy Override" section
7. Effective retention display: shows which policy is currently active (Global/Tenant/Mailbox) with badge
8. "Save Settings" button: updates retention policy via API, shows confirmation message
9. Confirmation modal for shortening retention: "Are you sure? Emails older than [X] years will be deleted. Type 'CONFIRM' to proceed."
10. Form validation: custom days must be >= 1 and <= 10950 (30 years max)
11. Accessibility: form accessible, warnings clearly announced to screen readers
12. Responsive design: mobile, tablet, desktop

## Story 5.3: Legal Hold Backend API

As a developer,
I want legal hold API endpoints,
so that emails can be preserved indefinitely for legal/compliance reasons.

**Acceptance Criteria:**

1. Endpoint: `PATCH /api/v1/tenants/:id/legal-hold` toggles legal hold for tenant (requires MSP Admin)
2. Accepts: `{legal_hold: true|false}`
3. Legal hold enabled: sets tenants.legal_hold=true, prevents deletion of ALL emails in tenant regardless of retention policy
4. Endpoint: `PATCH /api/v1/mailboxes/:id/legal-hold` toggles legal hold for specific mailbox (requires MSP Admin)
5. Accepts: `{legal_hold: true|false}`
6. Legal hold disabled: resumes normal retention policy enforcement on next cleanup job
7. Audit logging: log legal hold changes with user_id, tenant_id/mailbox_id, legal_hold value (true/false), timestamp
8. Legal hold indicator: GET tenant/mailbox endpoints return legal_hold status in response
9. Retention job integration: modify retention cleanup job to skip emails where tenant.legal_hold=true OR mailbox.legal_hold=true
10. RBAC: only MSP Admin can toggle legal hold (not Tenant Admin or User)
11. Unit tests cover: legal hold enforcement in retention logic, audit logging, RBAC checks
12. Integration test: enable legal hold, run retention cleanup, verify old emails NOT deleted

## Story 5.4: Legal Hold UI

As an MSP Admin,
I want to enable/disable legal hold,
so that I can preserve emails indefinitely for ongoing legal matters.

**Acceptance Criteria:**

1. Tenant detail page: "Legal Hold" section with toggle switch and status badge (Active/Inactive)
2. Legal hold description: "Legal Hold prevents automatic deletion of archived emails regardless of retention policy. Use for ongoing litigation, investigations, or compliance audits."
3. Toggle switch: click to enable/disable, shows confirmation modal
4. Confirmation modal: "Enable Legal Hold?" / "Disable Legal Hold?" with explanation of consequences
5. Enable confirmation: "All emails in this tenant will be preserved indefinitely until Legal Hold is disabled. Confirm?"
6. Disable confirmation: "Emails will resume normal retention policy enforcement. Old emails may be deleted during next cleanup. Confirm?"
7. Status badge: green "Active" badge if legal_hold=true, gray "Inactive" if false
8. Mailbox detail page: same Legal Hold section with toggle for mailbox-level hold
9. Visual indicator: tenants/mailboxes with legal hold active show shield icon in lists
10. Loading state during API call
11. Success/error messaging after toggle
12. Accessibility: toggle accessible, status changes announced
13. Responsive design: mobile, tablet, desktop

## Story 5.5: Audit Logging Backend Implementation

As a developer,
I want comprehensive audit logging,
so that all user actions are tracked for compliance and security auditing.

**Acceptance Criteria:**

1. Audit log triggers: log all state-changing actions (create, update, delete operations)
2. Logged actions: user login, user logout, tenant create/update/delete, mailbox enable/disable, search query, export request, settings changes, legal hold toggle, retention policy changes
3. Audit log record: user_id, action (enum or string), ip_address, timestamp, details (JSONB with action-specific data)
4. Details examples: search query includes {query, filters, result_count}, export includes {email_ids, format}, tenant create includes {tenant_name, azure_tenant_id}
5. Immutability enforcement: database trigger prevents UPDATE or DELETE on audit_logs table (append-only)
6. IP address capture: middleware extracts IP from X-Forwarded-For header (for reverse proxy) or remote address
7. Audit log middleware: intercepts all API requests, logs after successful completion (don't log failed requests due to auth)
8. Endpoint: `GET /api/v1/audit-logs?user_id={}&action={}&date_from={}&date_to={}&limit={}&offset={}` returns audit logs (requires MSP Admin)
9. Audit log list returns: log_id, user_email, action, ip_address, timestamp, details (expandable)
10. RBAC: only MSP Admin can view audit logs
11. Performance: audit logging is asynchronous (uses Redis queue) to not slow down main operations
12. Retention: audit logs retained indefinitely (not subject to email retention policy)
13. Unit tests: verify immutability trigger, middleware logging logic
14. Integration test: perform various actions, verify audit logs created correctly

## Story 5.6: Audit Logs Viewer UI

As an MSP Admin,
I want to view audit logs,
so that I can review user actions for compliance and security auditing.

**Acceptance Criteria:**

1. Audit logs page (`/audit-logs`) displays searchable/filterable audit log viewer (requires MSP Admin role)
2. Filter controls: User dropdown (All/{user names}), Action dropdown (All/Login/Search/Export/Settings Change/etc.), Date Range picker
3. Audit logs table: User Email, Action, IP Address, Timestamp, Details (expandable), Export button
4. Details column: "View Details" button opens modal/panel with full JSON details rendered as formatted key-value pairs
5. Table sorting: click column headers to sort (Timestamp default DESC)
6. Pagination: 100 logs per page, numbered pagination
7. Export functionality: "Export to CSV" button downloads filtered logs as CSV file
8. CSV includes: timestamp, user_email, action, ip_address, details (as JSON string)
9. Search: text input filters by action or user email (client-side or server-side)
10. Empty state: "No audit logs match your filters"
11. Loading state: skeleton loaders during fetch
12. Accessibility: table accessible, details modal keyboard-navigable
13. Responsive design: mobile (cards), tablet, desktop

## Story 5.7: Whitelabeling Backend API

As a developer,
I want whitelabeling configuration API,
so that MSP Admins can customize branding.

**Acceptance Criteria:**

1. Endpoint: `GET /api/v1/settings/whitelabel` returns whitelabel settings (requires authentication)
2. Whitelabel settings return: logo_url, favicon_url, custom_css
3. Endpoint: `PATCH /api/v1/settings/whitelabel` updates whitelabel settings (requires MSP Admin)
4. Logo upload: endpoint `POST /api/v1/settings/whitelabel/logo` accepts file upload (PNG/SVG/JPEG, max 2MB)
5. Logo file validation: check file size, verify MIME type, scan for malicious content (basic check)
6. Logo storage: save to filesystem `/archive/whitelabel/logo.{ext}`, return URL
7. Favicon upload: endpoint `POST /api/v1/settings/whitelabel/favicon` accepts file upload (ICO/PNG, max 500KB, dimensions 16x16, 32x32, 48x48)
8. Custom CSS: stored in database, sanitized to prevent malicious CSS (CSP-safe, no external URLs in url())
9. CSS sanitization: allow color/font/spacing properties, block dangerous properties (behavior, expression, moz-binding)
10. CSS validation: parse CSS, reject if invalid syntax
11. Whitelabel settings stored in settings table (whitelabel_logo_path, whitelabel_favicon_path, whitelabel_custom_css)
12. Settings applied globally for all users
13. Unit tests cover: file validation, CSS sanitization
14. Integration test: upload logo/favicon, verify served correctly

## Story 5.8: Whitelabeling Settings UI

As an MSP Admin,
I want to customize branding,
so that IronArchive matches my MSP's brand identity.

**Acceptance Criteria:**

1. Whitelabeling settings page (`/settings/whitelabel`) displays branding configuration form
2. Logo section: current logo preview, "Upload New Logo" file input (accepts PNG/SVG/JPEG), "Remove Logo" button (resets to default IronArchive logo)
3. Favicon section: current favicon preview, "Upload New Favicon" file input (accepts ICO/PNG)
4. Custom CSS section: code editor textarea with syntax highlighting, "Preview" button, "Reset to Default" button
5. File upload: drag-and-drop zone with "or click to browse" fallback
6. File validation: client-side check for file size/type before upload, shows error if invalid
7. Logo preview: shows uploaded logo immediately after upload, before "Save" button clicked (client-side preview)
8. Custom CSS preview: "Preview" button opens modal showing current page with custom CSS applied (temporary)
9. CSS examples/templates: dropdown with pre-made color scheme examples ("Blue Theme", "Green Theme", "Dark Red Theme") that populate CSS textarea
10. "Save Changes" button: uploads files and saves CSS, shows success message, reloads page to apply changes
11. Warning message: "Custom CSS changes may affect usability. Test thoroughly before applying."
12. Form validation: show file size/type errors before upload
13. Accessibility: form accessible, code editor keyboard-navigable
14. Responsive design: mobile, tablet, desktop

## Story 5.9: Theme System Backend API

As a developer,
I want user theme preferences API,
so that users can select personal themes.

**Acceptance Criteria:**

1. Endpoint: `GET /api/v1/profile/theme` returns user's theme preference (requires authentication)
2. Theme preference return: `{theme: "catppuccin-mocha" | "catppuccin-latte" | "nord" | "cyberpunk" | "dracula" | "tokyo-night"}`
3. Endpoint: `PATCH /api/v1/profile/theme` updates user's theme preference
4. Accepts: `{theme: "<theme_name>"}`
5. Validation: theme name must be one of the six predefined themes
6. Theme stored in users table (theme_preference column, default: "catppuccin-mocha")
7. Theme setting independent of whitelabeling: whitelabel CSS applies globally, theme applies per-user on top of whitelabel
8. Theme applied via CSS classes: frontend adds `theme-{name}` class to root element, CSS variables defined per theme
9. Unit tests cover validation, RBAC (users can only update their own theme)
10. Integration test: update theme, verify setting saved

## Story 5.10: Theme System Frontend Implementation

As a user,
I want to select a personal theme,
so that I can customize the interface to my preference.

**Acceptance Criteria:**

1. Theme selection in profile settings page (`/profile`): "Theme" section with visual theme picker
2. Theme picker displays 6 theme cards: Catppuccin Mocha, Catppuccin Latte, Nord, Cyberpunk, Dracula, Tokyo Night
3. Each theme card shows: theme name, preview swatch (shows primary/secondary/background colors), "Select" button
4. Currently selected theme highlighted with checkmark/border
5. Theme application: clicking "Select" immediately applies theme (no "Save" button needed), updates backend via API
6. Theme switching: smooth transition animation (CSS transitions on color changes)
7. Theme persistence: theme loaded from user preferences on page load, applied before first render (prevent flash of wrong theme)
8. Theme CSS variables defined: --color-primary, --color-secondary, --color-background, --color-text, --color-border, --color-error, --color-success, etc.
9. All components use CSS variables (not hard-coded colors) for theme-ability
10. Dark/light mode compatibility: Catppuccin Mocha, Nord, Cyberpunk, Dracula, Tokyo Night are dark themes; Catppuccin Latte is light theme
11. Accessibility: sufficient contrast ratios in all themes (WCAG 2.1 AA), theme cards keyboard-navigable
12. Responsive design: theme picker adapts to mobile/tablet/desktop

## Story 5.11: Profile Settings Page

As a user,
I want to manage my profile settings,
so that I can update my display name, password, and MFA configuration.

**Acceptance Criteria:**

1. Profile settings page (`/profile`) displays user profile management form
2. Form sections: Display Name, Email (read-only, cannot change), Password Change, Multi-Factor Authentication (MFA), Theme Selection (from Story 5.10)
3. Display Name section: text input, "Update Display Name" button
4. Password Change section: Current Password input, New Password input, Confirm New Password input, "Change Password" button
5. Password change validation: current password correct, new password meets complexity requirements, passwords match
6. MFA section: "MFA Status" badge (Enabled/Disabled), "Enable MFA" / "Disable MFA" button
7. Enable MFA flow: generates TOTP secret, displays QR code, asks user to scan with authenticator app, enter verification code to confirm
8. MFA verification: user enters 6-digit code from authenticator app, backend validates code against generated secret
9. MFA enabled success: stores encrypted TOTP secret in users.mfa_secret, shows "MFA Enabled" badge, displays backup codes (one-time use codes for account recovery)
10. Disable MFA flow: requires current password + TOTP code to disable, removes mfa_secret from database
11. Form validation: display name non-empty, passwords meet requirements
12. API endpoints: `PATCH /api/v1/profile` (update display name), `POST /api/v1/profile/change-password`, `POST /api/v1/profile/mfa/enable`, `POST /api/v1/profile/mfa/disable`
13. Loading states: during API calls
14. Success/error messaging: for each operation
15. Accessibility: form accessible, MFA QR code has alt text with manual entry key
16. Responsive design: mobile, tablet, desktop

## Story 5.12: MFA Login Flow

As a developer,
I want MFA verification in login flow,
so that users with MFA enabled must provide TOTP code to authenticate.

**Acceptance Criteria:**

1. Login flow modification: if user has mfa_secret set, require TOTP verification after password validation
2. Login endpoint logic: POST /api/v1/auth/login validates email/password, if valid AND mfa_enabled, return `{mfa_required: true, mfa_token: "<temporary_token>"}`
3. MFA token: short-lived JWT (5 minutes) containing user_id, used for MFA verification step
4. Frontend login flow: if mfa_required=true, show MFA verification step (input for 6-digit code)
5. MFA verification endpoint: `POST /api/v1/auth/mfa/verify` accepts `{mfa_token, totp_code}`
6. Verification logic: validates mfa_token, extracts user_id, retrieves user's mfa_secret, validates TOTP code using time-based algorithm
7. TOTP validation: accepts codes within 1 time step window (30 seconds before/after current time) to account for clock skew
8. Verification success: returns full access/refresh JWT tokens, user can proceed to dashboard
9. Verification failure: returns 401 Unauthorized with error message "Invalid MFA code"
10. Rate limiting: limit to 5 MFA verification attempts per user per 15 minutes (prevent brute force)
11. Backup codes: during MFA verification, allow entry of 8-character backup code instead of TOTP code (one-time use, mark as used in database)
12. Unit tests cover: TOTP code generation/validation, backup code usage, rate limiting
13. Integration test: enable MFA for test user, login with email/password, verify with TOTP code, access protected endpoint

## Story 5.13: Documentation and README

As a potential user,
I want comprehensive documentation,
so that I can install, configure, and use IronArchive successfully.

**Acceptance Criteria:**

1. README.md includes: project description, key features list, quick start guide, installation instructions (Docker Compose), system requirements, links to detailed documentation
2. Installation guide (`/docs/installation.md`): detailed steps for Docker Compose deployment, environment variable configuration, first-time setup wizard
3. User guide (`/docs/user-guide.md`): covers all user workflows (tenant onboarding, search, export, settings, profile)
4. Admin guide (`/docs/admin-guide.md`): covers MSP Admin tasks (user management, whitelabeling, retention policies, legal hold, audit logs)
5. Architecture documentation (`/docs/architecture.md`): system design overview, component diagram, database schema, technology stack, API architecture
6. API documentation (`/docs/api.md`): all endpoints documented with request/response examples, authentication requirements, error codes
7. Compliance documentation (`/docs/compliance.md`): DSGVO/GoBD compliance mapping, audit features, retention policy explanation, legal hold usage
8. Troubleshooting guide (`/docs/troubleshooting.md`): common issues and solutions, error message explanations, debugging tips
9. Contributing guide (`CONTRIBUTING.md`): how to contribute code, coding standards, pull request process, development environment setup
10. Code of Conduct (`CODE_OF_CONDUCT.md`): community standards and expectations
11. License (`LICENSE`): MIT License text
12. Changelog (`CHANGELOG.md`): version history, feature additions, bug fixes (v1.0 initial release entry)
13. Documentation written in clear, concise language with screenshots/diagrams where helpful
14. All documentation reviewed for accuracy and completeness

## Story 5.14: Final Polish and Testing

As a developer,
I want the application polished and thoroughly tested,
so that v1.0 is production-ready.

**Acceptance Criteria:**

1. Visual design consistency: review all pages for consistent spacing, typography, color usage, button styles
2. Error message review: all error messages are user-friendly, actionable, and grammatically correct
3. Loading states: all async operations have appropriate loading indicators (spinners, skeleton loaders)
4. Empty states: all lists/tables have empty state designs with helpful instructions
5. Responsive design verification: test all pages on mobile (iPhone 12/13/14), tablet (iPad), desktop (1920x1080, 2560x1440)
6. Browser compatibility testing: test on Chrome, Firefox, Safari (macOS), Edge
7. Accessibility audit: run automated accessibility checker (axe, Lighthouse), fix critical issues, verify WCAG 2.1 AA compliance
8. Performance testing: test search with 100K indexed emails, verify <200ms response time; test sync with 10K email mailbox, verify throughput >10K emails/hour
9. Security review: review authentication/authorization logic, check for SQL injection vulnerabilities (parameterized queries), verify CSRF protection, test XSS prevention (HTML sanitization)
10. Integration testing: end-to-end tests for critical workflows (setup wizard → tenant onboarding → sync → search → export)
11. Error handling testing: test failure scenarios (database down, Redis down, Meilisearch down, Graph API errors), verify graceful degradation
12. Docker image build: build multi-arch images (ARM64 + AMD64), test deployment on Ubuntu 22.04 LTS
13. Documentation review: verify all documentation accurate and matches implemented features
14. Release notes: write comprehensive v1.0 release notes listing all features, known limitations, installation instructions
15. GitHub release preparation: tag v1.0.0, publish release with changelog, attach Docker Compose file and .env template

---
