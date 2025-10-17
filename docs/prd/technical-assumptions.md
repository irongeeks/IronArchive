# Technical Assumptions

## Repository Structure

**Monorepo Architecture**

IronArchive will use a monorepo structure to simplify development, dependency management, and deployment. All code (backend, frontend, shared utilities) resides in a single Git repository with well-defined subdirectories:

```
/ironarchive
├── /backend          # Go application (API, workers, sync engine)
├── /frontend         # SvelteKit application
├── /docker           # Dockerfile, docker-compose.yml
├── /docs             # Documentation
├── /migrations       # Database schema migrations
├── /scripts          # Utility scripts
└── README.md
```

**Rationale:** Monorepo enables atomic commits across frontend/backend changes, simplifies versioning, and reduces coordination overhead for solo developer. Separate repositories would create unnecessary complexity for v1.0.

## Service Architecture

**Monolithic Start with Planned Separation Path**

**v1.0 Approach:** Single Go binary containing API server, background job processor, and sync workers. All components share same codebase and can communicate via in-memory channels or Redis queues as needed.

**Advantages:**
- Simplified deployment (single Docker container for application logic)
- Reduced operational complexity (no inter-service communication debugging)
- Faster development iteration (no API versioning between services)
- Lower resource usage (shared memory, single process overhead)

**Future Separation Strategy (Post-v1.0):**
If scale demands, split into:
- **API Service:** Handles HTTP requests, authentication, serves frontend
- **Sync Workers:** Dedicated processes for Microsoft Graph API communication and email archiving
- **Job Processor:** Background task execution (exports, retention cleanup)

**Communication:** Redis-based job queue (asynq library) and shared PostgreSQL database

**Stateless Design:** All persistent state in PostgreSQL/Redis enables horizontal scaling when needed

## Testing Requirements

**Balanced Testing Pyramid for MVP**

**Unit Tests (Foundation):**
- Critical business logic (RBAC enforcement, retention policy calculations, attachment deduplication)
- Utility functions (email parsing, filesystem path generation, date handling)
- Target: 60%+ code coverage for backend core packages
- Tools: Go standard `testing` package, `testify` for assertions

**Integration Tests (Selective):**
- API endpoint testing (authentication flows, CRUD operations, search queries)
- Database interactions (migrations, complex queries, transaction behavior)
- Microsoft Graph API client (mocked for CI, real API in manual testing)
- Tools: Go `httptest`, dockertest for PostgreSQL/Redis spin-up

**End-to-End Tests (Critical Paths Only):**
- Setup wizard flow (first MSP Admin creation)
- Tenant onboarding wizard (mocked Azure AD app creation)
- Email search and export workflow
- Tools: Playwright for browser automation (JavaScript/TypeScript)

**Manual Testing:**
- Full Microsoft Graph API integration with real M365 tenant
- Whitelabeling and theme switching across devices
- Performance testing with large datasets (1M+ emails)
- Accessibility testing with screen readers (NVDA, VoiceOver)

**NO Extensive Testing for v1.0:**
- ❌ Load testing / stress testing (defer until real-world usage patterns emerge)
- ❌ Security penetration testing (basic security review only, professional pentest post-v1.0)
- ❌ Multi-browser compatibility testing (test on 2-3 primary browsers only)

**Rationale:** Focus testing resources on correctness and critical functionality. Over-testing slows v1.0 delivery without proportional risk reduction for initial release.

## Additional Technical Assumptions and Requests

**Backend Stack:**
- **Language:** Go 1.24 (latest stable as of development start)
- **Web Framework:** Fiber v3 (fast, Express-like API, built on fasthttp)
- **Database Driver:** pgx (native PostgreSQL driver) or GORM if productivity benefits justify slight performance trade-off
- **Authentication:** JWT via `golang-jwt/jwt` library, bcrypt via `golang.org/x/crypto/bcrypt`
- **Job Queue:** asynq (Redis-based, mature, supports retries and scheduling)
- **Graph API Client:** `msgraph-sdk-go` (official Microsoft library)

**Frontend Stack:**
- **Framework:** SvelteKit 2.x with SSR (server-side rendering) + SPA (single-page app) hybrid mode
- **Styling:** TailwindCSS 4.x with custom theme configuration for color palette system
- **Component Library:** Shadcn-Svelte (accessible components based on Radix UI primitives, Svelte port)
- **Charts:** Chart.js (simple, lightweight, sufficient for storage/usage widgets)
- **HTTP Client:** Native fetch API with SvelteKit load functions for SSR data fetching
- **State Management:** Svelte stores (built-in reactivity, no external state library needed)

**Database:**
- **Primary Database:** PostgreSQL 16+ with pgcrypto and pg_cron extensions
- **Migrations:** golang-migrate (version-controlled schema changes, up/down migrations)
- **Backup Strategy:** Daily pg_dump via cron, WAL archiving for point-in-time recovery (documented for users)
- **Indexing Strategy:** B-tree indexes on foreign keys and frequently queried columns, GIN indexes for JSONB audit log fields

**Infrastructure:**
- **Deployment:** Docker Compose (simple, sufficient for 95% of MSPs)
- **Container Images:** Multi-arch builds (ARM64 + AMD64) via Docker Buildx
- **Reverse Proxy:** Traefik 2.10+ (automatic HTTPS via Let's Encrypt, Docker label-based configuration)
- **File Storage:** Local filesystem with configurable mount point (enables NAS/SAN integration for large deployments)
- **Monitoring:** Prometheus exporters available but not required (users can optionally deploy Grafana dashboards)

**External Integrations:**
- **Microsoft Graph API:** OAuth2 client credentials flow, v1.0 stable endpoints
- **Email SMTP:** Standard library SMTP client for notifications
- **Teams Webhooks:** HTTP POST with Adaptive Card JSON payload
- **Discord Webhooks:** HTTP POST with embed JSON payload

**Security:**
- **Encryption at Rest:** Document recommendation for LUKS/dm-crypt at OS level (PostgreSQL TDE optional)
- **Encryption in Transit:** TLS 1.3 enforced for all external connections
- **Secrets Management:** Environment variables (`.env` file for Docker Compose), no hardcoded secrets
- **Multi-Tenant Isolation:** Application-enforced tenant filtering on all queries (consider PostgreSQL RLS in future)
- **CSRF Protection:** SameSite cookies + CSRF tokens for state-changing requests
- **Input Validation:** Strict validation with Go validator library, parameterized SQL queries

**Development Environment:**
- **Go Version:** 1.24
- **Node.js Version:** 20 LTS
- **Docker Desktop:** Latest stable (or Docker Engine + Docker Compose on Linux)
- **IDE Recommendations:** VSCode with Go and Svelte extensions (documented in CONTRIBUTING.md)
- **Local Development:** Docker Compose stack with hot-reload for backend (air) and frontend (Vite HMR)

**Compliance & Data Handling:**
- **DSGVO/GoBD Features:** Retention policies enforced via scheduled PostgreSQL jobs (pg_cron), legal hold via database flag
- **Audit Logs:** Immutable table with database triggers preventing updates/deletes, includes user_id, ip_address, action, timestamp, details (JSONB)
- **Data Retention:** Configurable at global/tenant/mailbox levels with automatic deletion after expiration
- **Right to Erasure:** Manual process via MSP Admin (export user data, delete from archive, document in audit log)

**Known Technical Constraints:**
- **Microsoft Graph API Rate Limits:** ~10,000 requests per 10 minutes per tenant; sync engine must implement backoff and queuing
- **Meilisearch Memory Usage:** Approximately 1-2 GB RAM per 1M indexed emails; document scaling recommendations
- **Filesystem Performance:** Millions of small files may impact performance on certain filesystems (recommend ext4 or XFS over NFS)
- **PST Export Library:** go-pst exists but may need evaluation; fallback to shell-out to readpst if necessary
