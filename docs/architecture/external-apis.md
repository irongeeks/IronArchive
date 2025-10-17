## External APIs

IronArchive integrates with Microsoft Graph API for email synchronization and supports multiple webhook channels for notifications.

### Microsoft Graph API

- **Purpose:** Email synchronization, mailbox discovery, attachment retrieval
- **Documentation:** https://learn.microsoft.com/en-us/graph/api/overview
- **Base URL(s):** `https://graph.microsoft.com/v1.0/`, `https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token` (auth)
- **Authentication:** OAuth2 Client Credentials Flow (app-only authentication)
- **Rate Limits:** ~10,000 requests per 10 minutes per tenant, throttling headers: `Retry-After`

**Key Endpoints Used:**
- `POST /oauth2/v2.0/token` - Obtain access token with client credentials
- `GET /users` - List all mailboxes in tenant
- `GET /users/{id}/messages/delta` - Fetch emails with delta query (initial and incremental sync)
- `GET /users/{id}/messages/{messageId}` - Retrieve specific email message
- `GET /users/{id}/messages/{messageId}/attachments/{attachmentId}/$value` - Download attachment content

**Integration Notes:**
- Token caching: Access tokens valid for 60 minutes, cache in memory with expiration tracking
- Rate limiting: Implement exponential backoff on 429 responses, queue-based throttling for high-volume sync
- Delta token expiration: Delta tokens expire after 30 days, automatic fallback to initial sync if expired
- Error handling: Distinguish transient (503, 429) from permanent errors (401, 404), retry only transient
- Concurrent requests: Limit to 5 concurrent requests per tenant to avoid rate limits

### SMTP Server (Email Notifications)

- **Purpose:** Send email notifications for sync failures, export completions, alerts
- **Documentation:** Standard SMTP protocol (RFC 5321)
- **Base URL(s):** Configurable via env var `SMTP_HOST:SMTP_PORT`
- **Authentication:** SMTP AUTH (username/password via env vars)
- **Rate Limits:** Depends on configured SMTP provider (typically 100-1000 emails/hour for transactional email services)

**Key Operations:**
- Send HTML emails with embedded tables for job summaries
- Support for TLS/STARTTLS encryption
- Template-based email generation with job details

**Integration Notes:** Use `net/smtp` standard library, implement connection pooling for efficiency, retry failed sends up to 3 times

### Microsoft Teams Webhooks

- **Purpose:** Send rich notifications to Teams channels for MSP admin monitoring
- **Documentation:** https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook
- **Base URL(s):** User-provided webhook URL (e.g., `https://outlook.office.com/webhook/...`)
- **Authentication:** None (webhook URL is secret)
- **Rate Limits:** ~20 requests per second per webhook

**Key Integration:**
- POST JSON payload with Adaptive Card format
- Include sync status, email counts, error details, tenant name
- Visual formatting: color-coded cards (green success, red failure, blue in-progress)

**Integration Notes:** Implement 3-retry logic for network failures, validate webhook URL on settings save

### Discord Webhooks

- **Purpose:** Send notifications to Discord channels for community-focused MSPs
- **Documentation:** https://discord.com/developers/docs/resources/webhook
- **Base URL(s):** User-provided webhook URL (e.g., `https://discord.com/api/webhooks/{id}/{token}`)
- **Authentication:** None (webhook URL is secret)
- **Rate Limits:** 30 requests per minute per webhook

**Key Integration:**
- POST JSON payload with rich embeds
- Include job details, timestamps, status indicators
- Color coding: green (success), red (failure), blue (in-progress)

**Integration Notes:** Rate limit handling with local queue if exceeding 30/min, retry logic for 5xx errors

