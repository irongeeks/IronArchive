# Epic 3: Search, Retrieval & Export

**Goal:** Implement Meilisearch-powered instant search with advanced filters, email detail side panel, and export functionality supporting EML/ZIP and PST formats. Deliver the core user-facing value proposition—fast, intuitive search across archived emails with seamless export workflows. This epic makes the archive practically useful for end users, tenant admins, and MSP admins performing email discovery and compliance tasks.

## Story 3.1: Search API Endpoint with Meilisearch Integration

As a developer,
I want a search API endpoint,
so that users can query archived emails with instant results.

**Acceptance Criteria:**

1. Endpoint: `GET /api/v1/search?q={query}&tenant_id={id}&mailbox_id={id}&filters={json}` performs search (requires authentication)
2. Query parameter `q`: full-text search across subject, sender, recipients, body_text
3. Optional filters: `tenant_id`, `mailbox_id`, `date_from`, `date_to`, `has_attachments` (boolean), `size_min`, `size_max`
4. Meilisearch query execution with typo tolerance, highlighting enabled
5. RBAC filtering: MSP Admin sees all emails, Tenant Admin sees only assigned tenant emails, User sees only own mailbox emails
6. RBAC implementation: apply tenant/mailbox filters automatically based on user role and permissions
7. Search results return: email id, subject (with highlighted terms), sender, recipients, sent_at, snippet (body preview with highlighted terms), has_attachments, size_bytes
8. Pagination: `limit` (default 50, max 100) and `offset` parameters
9. Sorting: `sort` parameter (date_desc, date_asc, size_desc, size_asc, relevance)
10. Response time: 95% of queries return in < 200ms (verify with logging)
11. Error handling: invalid query syntax returns 400 Bad Request, Meilisearch unavailable returns 503 Service Unavailable
12. Unit tests cover: query building, RBAC filter application, error handling
13. Integration test: index 100 test emails, perform searches with various filters, verify RBAC enforcement

## Story 3.2: Search Page Frontend UI

As a user,
I want a search interface,
so that I can find archived emails quickly.

**Acceptance Criteria:**

1. Search page (`/search`) displays: search bar, filters sidebar, results list, side panel (initially hidden)
2. Search bar: large text input with placeholder "Search emails...", instant search (debounced 300ms after typing stops)
3. Search bar includes search icon and clear button (X) to reset query
4. Filters sidebar includes: Date Range (date pickers for from/to), Sender (text input), Recipient (text input), Has Attachments (checkbox), Size Range (min/max inputs with byte/KB/MB selector)
5. "Apply Filters" button triggers search with filters
6. "Clear Filters" button resets all filters to default
7. Results list displays: email subject (highlighted terms), sender, sent date (relative time: "2 hours ago", "3 days ago"), snippet (body preview with highlighted terms), attachment icon (if has_attachments), size badge
8. Results list includes infinite scroll: loads 50 results initially, loads next 50 on scroll to bottom
9. "Back to Top" button appears after scrolling past 100 results
10. Each result is clickable: opens side panel with email details
11. Empty state: "No results found" message with search tips
12. Loading state: skeleton loaders during search API call
13. Error state: error message with retry button
14. Accessibility: search bar keyboard-accessible, results announced to screen readers, keyboard navigation between results
15. Responsive design: mobile (filters collapse to modal), tablet, desktop

## Story 3.3: Email Detail Side Panel

As a user,
I want to view full email details in a side panel,
so that I can read emails without leaving the search results.

**Acceptance Criteria:**

1. Side panel slides in from right when email result clicked
2. Side panel header: "Close" button (X), "Export" button dropdown (export this email)
3. Side panel content sections: Email Metadata (subject, sender, recipients, sent date, size), Email Body (HTML or plain text rendering), Attachments List (filename, size, download button per attachment)
4. Email body rendering: if body_html exists, render HTML (sanitized to prevent XSS); otherwise render body_text in preformatted block
5. HTML sanitization: strip script tags, iframes, form elements; allow safe HTML (p, div, span, a, img with src restrictions)
6. Attachments list: each attachment shows icon (based on file type), filename, size, "Download" button
7. Attachment download: `GET /api/v1/attachments/:id/download` returns attachment file with correct content-type header
8. Side panel close: click "Close" button, click outside panel (overlay), or press Escape key
9. Side panel persists during infinite scroll (doesn't close when more results load)
10. Loading state while fetching email details (if not already loaded)
11. Error state: if email details fail to load, show error message in panel
12. Accessibility: panel focus trapped while open, close button keyboard-accessible, focus returns to clicked result on close
13. Responsive design: mobile (panel becomes full-screen modal), tablet (panel width 50%), desktop (panel width 40%)

## Story 3.4: Multi-Select and Export UI

As a user,
I want to select multiple emails and export them,
so that I can download emails in bulk for offline access or legal discovery.

**Acceptance Criteria:**

1. Each email result includes checkbox (left side of result card)
2. "Select All" checkbox in results header selects all visible results
3. Selection counter appears when any email selected: "X emails selected" with "Clear Selection" link
4. Bulk actions toolbar appears when any email selected: "Export Selected" button
5. "Export Selected" button opens export modal
6. Export modal includes: format selector (EML/ZIP for <=100 emails, PST for >100 emails or if user chooses), estimated file size display, "Start Export" button
7. Export modal shows format descriptions: EML/ZIP (individual .eml files in ZIP archive, compatible with all email clients), PST (Outlook PST format, for large exports)
8. Format availability: EML/ZIP always available, PST available only if >100 emails selected OR user selects it (with warning if <100 emails: "PST recommended for large exports")
9. "Start Export" button triggers export job API call: `POST /api/v1/exports` with payload `{email_ids: [], format: "eml_zip"|"pst"}`
10. Export API returns job ID, modal updates to show "Export in progress..." with job ID and polling status
11. Small exports (<100MB): API returns download link immediately when job completes, modal shows "Download Ready" button
12. Large exports (>=100MB): background job processes export, user receives notification (email/Teams/Discord) when ready with download link
13. Export history accessible: "Export History" link in user profile menu shows past exports with download links (unexpired)
14. Download link expiration: 7 days, displayed to user ("Expires in 6 days, 23 hours")
15. Accessibility: checkboxes keyboard-accessible, selection state announced to screen readers, export modal keyboard-navigable
16. Responsive design: export modal adapts to mobile/tablet/desktop

## Story 3.5: Export Job Backend (EML/ZIP Generation)

As a developer,
I want EML/ZIP export functionality,
so that users can download emails in standard format.

**Acceptance Criteria:**

1. Export job processor: `ProcessExportJob(jobID)` handles export job execution
2. Job fetches email records from database based on email_ids array
3. For each email: generate .eml file (RFC 822 format) with headers, body, and embedded attachments (Base64 encoded)
4. EML generation includes: standard email headers (From, To, Subject, Date, Message-ID), MIME multipart structure for attachments, proper content encoding
5. ZIP archive creation: creates ZIP file containing all .eml files, named with timestamp: `ironarchive-export-{job_id}-{timestamp}.zip`
6. ZIP file stored in temporary exports directory: `/tmp/exports/{job_id}.zip`
7. File size calculation: estimate before generation (sum of email sizes), actual size after generation
8. Small file (<100MB): job completes synchronously, returns download URL immediately
9. Large file (>=100MB): job processes asynchronously in background worker, sends notification on completion
10. Download endpoint: `GET /api/v1/exports/:job_id/download` serves ZIP file with correct headers (Content-Type: application/zip, Content-Disposition: attachment)
11. Download link: `https://ironarchive.example.com/api/v1/exports/{job_id}/download` (valid for 7 days)
12. Cleanup job: daily cron job deletes export files older than 7 days from `/tmp/exports/`
13. RBAC: users can only download their own exports (job owner check)
14. Error handling: if email records not found, return 404; if generation fails, mark job as failed with error message
15. Logging: INFO for export start/completion with file size, ERROR for failures
16. Unit tests cover: EML generation, ZIP creation, RBAC checks
17. Integration test: export 10 emails, download ZIP, verify .eml files valid and openable in email client

## Story 3.6: Export Job Backend (PST Generation)

As a developer,
I want PST export functionality,
so that users can download large email sets in Outlook-compatible format.

**Acceptance Criteria:**

1. PST export job processor: `ProcessPSTExportJob(jobID)` handles PST generation
2. Job fetches email records from database based on email_ids array
3. PST library integration: use `go-pst` library (or fallback to shell-out to `readpst` if library insufficient)
4. PST generation: creates Outlook PST file with folder structure mirroring mailbox (Inbox, Sent, etc. if metadata available; otherwise single "Archived Emails" folder)
5. For each email: add message to PST with headers, body (HTML + plain text), attachments
6. PST file stored in temporary exports directory: `/tmp/exports/{job_id}.pst`
7. File size: PST format is compressed, typically 60-70% of raw email size
8. PST generation is resource-intensive: always processed as background job (no synchronous generation)
9. Progress tracking: update job metadata with progress percentage (emails_processed / total_emails * 100)
10. Notification on completion: send email/Teams/Discord notification with download link
11. Download endpoint: `GET /api/v1/exports/:job_id/download` serves PST file (Content-Type: application/vnd.ms-outlook, Content-Disposition: attachment)
12. Error handling: if PST generation fails (library error), mark job as failed with error message and suggest trying EML/ZIP format
13. Logging: INFO for export start/progress/completion, ERROR for failures
14. Unit tests cover: PST generation (if library supports it), error handling
15. Integration test: export 100 emails to PST, download file, verify openable in Microsoft Outlook or other PST-compatible client

## Story 3.7: Export History Page

As a user,
I want to view my export history,
so that I can re-download previous exports or see their status.

**Acceptance Criteria:**

1. Export history page (`/exports`) displays list of user's past exports
2. Exports list includes columns: Export Date, Format (EML/ZIP or PST), Email Count, File Size, Status (Completed/In Progress/Failed), Download Link (if available), Expires In (countdown timer)
3. Status badges: Completed (green), In Progress (blue with spinner), Failed (red)
4. Download link: blue button "Download" if export completed and not expired, disabled button "Expired" if expired, "Retry" button if failed
5. Failed exports show error message tooltip on hover
6. In Progress exports show progress bar if progress percentage available
7. Auto-refresh: page auto-refreshes every 5 seconds to update in-progress export statuses
8. Pagination: 20 exports per page, "Load More" button at bottom
9. Empty state: "No exports yet" message with link to search page
10. Export deletion: "Delete" button (icon) for each export, confirmation modal before deletion
11. Delete export: removes job record and deletes file from filesystem (if exists)
12. RBAC: users see only their own exports
13. Accessibility: table accessible, status changes announced, keyboard navigation
14. Responsive design: mobile (table becomes cards), tablet, desktop

---
