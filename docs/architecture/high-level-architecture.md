## High Level Architecture

### Technical Summary

IronArchive is architected as a **monolithic fullstack application** with clear separation between backend (Go), frontend (SvelteKit), and supporting services (PostgreSQL, Redis, Meilisearch). The system employs a **multi-tenant SaaS architecture** where MSP Admins manage multiple customer tenants, each with their own M365 mailboxes, while maintaining strict data isolation through application-enforced tenant filtering.

The backend leverages **Go 1.24 with Fiber v3** for high-performance HTTP APIs and background job processing, integrating with **Microsoft Graph API** for email synchronization using delta queries for efficiency. Email data is stored in a **hybrid approach**: metadata and relationships in **PostgreSQL 16**, full message bodies on the filesystem with hierarchical directory structure, searchable content in **Meilisearch** for sub-200ms queries, and background tasks orchestrated via **Redis + asynq** job queues.

The frontend utilizes **SvelteKit 2.x** in hybrid SSR+SPA mode, providing server-side rendering for initial page loads and seamless client-side navigation thereafter. **TailwindCSS 4.x** with a custom theme system enables comprehensive whitelabeling for MSPs to brand the platform for their clients.

Deployment follows **Docker Compose** for single-server orchestration with **Traefik** reverse proxy for automatic HTTPS via Let's Encrypt. This architecture targets **cost-efficiency** (75% reduction vs commercial solutions) while maintaining **enterprise-grade compliance** features (DSGVO/GoBD retention policies, legal holds, immutable audit logs) and **superior performance** (instant search, 10K+ emails/hour sync throughput).

The monolithic v1.0 approach prioritizes **rapid development and operational simplicity** for the initial release, with a clear separation path to microservices post-v1.0 if scale demands (API server, sync workers, job processors as separate containers).

### Platform and Infrastructure Choice

**Platform:** Self-Hosted Docker Compose (Multi-Cloud Compatible)

**Key Services:**
- **Application Container:** Go binary (API server + background workers + job queue consumer)
- **Database:** PostgreSQL 16 with pgcrypto, pg_cron extensions
- **Cache & Job Queue:** Redis 7
- **Search Engine:** Meilisearch 1.6
- **Reverse Proxy:** Traefik 2.10+ (automatic HTTPS, Docker label-based routing)
- **File Storage:** Local filesystem with configurable mount point (supports NFS/SAN for large deployments)

**Deployment Hosts and Regions:**
- **Recommended:** Hetzner Cloud (Nuremberg, Helsinki) or OVH (Gravelines) for European data residency
- **Supported:** AWS EC2, Azure VMs, Google Compute Engine, bare metal servers
- **Multi-Region Strategy:** Each MSP deploys a single instance in their preferred region

### Repository Structure

**Structure:** Monorepo
**Monorepo Tool:** Native Git monorepo (no external tool)
**Package Organization:** Directory-based separation

```
/ironarchive
├── /backend          # Go application (API, workers, sync engine)
│   ├── /cmd          # Application entrypoints
│   ├── /internal     # Private application code
│   └── /pkg          # Public shared packages
├── /frontend         # SvelteKit application
│   ├── /src          # Source code (routes, components, stores)
│   └── /static       # Static assets
├── /docker           # Dockerfile, docker-compose.yml
├── /docs             # Documentation (architecture, PRD, guides)
├── /migrations       # Database schema migrations (golang-migrate)
├── /scripts          # Utility scripts (build, deploy, test)
└── README.md
```

### High Level Architecture Diagram

```mermaid
graph TB
    subgraph "User Access"
        MSPAdmin[MSP Admin<br/>Firefox/Chrome]
        TenantAdmin[Tenant Admin<br/>Firefox/Chrome]
        User[End User<br/>Firefox/Chrome]
    end

    subgraph "Internet"
        M365[Microsoft 365<br/>Graph API]
        SMTP[SMTP Server<br/>Email Notifications]
        Teams[Microsoft Teams<br/>Webhooks]
        Discord[Discord<br/>Webhooks]
    end

    subgraph "IronArchive Infrastructure"
        subgraph "Traefik Reverse Proxy"
            Traefik[Traefik 2.10+<br/>HTTPS / Let's Encrypt]
        end

        subgraph "Application Container"
            API[Fiber API Server<br/>Go 1.24]
            Worker[Sync Workers<br/>Background Jobs]
            Scheduler[Cron Scheduler<br/>4x Daily Syncs]

            API -.In-Process.-> Worker
            Scheduler -.In-Process.-> Worker
        end

        subgraph "Frontend Container"
            Frontend[SvelteKit 2.x<br/>SSR + SPA]
        end

        subgraph "Data Layer"
            Postgres[(PostgreSQL 16<br/>Metadata & Relations)]
            Redis[(Redis 7<br/>Job Queue & Cache)]
            Meilisearch[(Meilisearch 1.6<br/>Search Index)]
            Filesystem[/Filesystem Storage<br/>Email Bodies & Attachments/]
        end
    end

    MSPAdmin -->|HTTPS| Traefik
    TenantAdmin -->|HTTPS| Traefik
    User -->|HTTPS| Traefik

    Traefik -->|Proxy /api/*| API
    Traefik -->|Proxy /*| Frontend

    Frontend -.Fetch API.-> API

    API -->|Read/Write| Postgres
    API -->|Pub/Sub & Cache| Redis
    API -->|Index & Search| Meilisearch
    API -->|Store Files| Filesystem

    Worker -->|OAuth2 Client Credentials| M365
    Worker -->|Read/Write| Postgres
    Worker -->|Store Files| Filesystem
    Worker -->|Index| Meilisearch
    Worker -->|Dequeue Jobs| Redis

    API -->|Send Notifications| SMTP
    API -->|Send Notifications| Teams
    API -->|Send Notifications| Discord

    Scheduler -->|Enqueue Jobs| Redis

    style Traefik fill:#326CE5,stroke:#fff,color:#fff
    style API fill:#00ADD8,stroke:#fff,color:#fff
    style Worker fill:#00ADD8,stroke:#fff,color:#fff
    style Frontend fill:#FF3E00,stroke:#fff,color:#fff
    style Postgres fill:#336791,stroke:#fff,color:#fff
    style Redis fill:#DC382D,stroke:#fff,color:#fff
    style Meilisearch fill:#FF5CAA,stroke:#fff,color:#fff
```

### Architectural Patterns

- **Monolithic Architecture with Clear Boundaries:** Single deployable Go binary containing API server, background workers, and cron scheduler communicating via in-process channels and Redis queues. Simplifies v1.0 deployment while maintaining logical separation for future microservices extraction. _Rationale:_ Reduces operational complexity, eliminates inter-service network latency, and accelerates development velocity for initial release.

- **Multi-Tenant SaaS with Application-Enforced Isolation:** All database queries filtered by `tenant_id` via middleware/repository pattern. Each tenant's data logically isolated in shared database tables and physically separated on filesystem by directory structure. _Rationale:_ Balances cost-efficiency (shared infrastructure) with security (strict isolation prevents cross-tenant data leaks).

- **Event-Driven Background Job Processing:** Long-running operations (email sync, export generation) processed asynchronously via Redis-backed job queues (asynq library) with retry logic and progress tracking. _Rationale:_ Keeps API responses fast (<500ms), handles rate limits gracefully, and provides fault tolerance for external API interactions.

- **Hybrid Data Storage Strategy:** Relational data (users, tenants, mailboxes) in PostgreSQL for ACID transactions; email metadata in PostgreSQL with full bodies on filesystem for cost-efficiency; searchable content in Meilisearch for performance. _Rationale:_ Optimizes for query speed (search), cost (filesystem cheaper than database BLOB storage), and relational integrity (foreign keys).

- **Repository Pattern for Data Access:** Abstract database operations behind repository interfaces (`UserRepository`, `TenantRepository`, `EmailRepository`) with implementations using pgx or GORM. _Rationale:_ Enables unit testing with mocked repositories, centralizes tenant filtering logic, and provides future flexibility to swap database implementations.

- **API Gateway Pattern (Frontend BFF):** SvelteKit server-side routes act as Backend-for-Frontend, aggregating multiple API calls and handling session management. _Rationale:_ Reduces frontend complexity, minimizes round trips, and enables server-side rendering with authenticated data fetching.

- **Component-Based UI with Atomic Design:** Reusable Svelte components organized by atomic design principles (atoms, molecules, organisms, templates, pages) using Shadcn-Svelte primitives. _Rationale:_ Ensures consistency across whitelabeled instances and accelerates UI development with pre-built accessible components.

- **Outbox Pattern for Notifications:** Critical notifications (sync failures, export completions) written to database `notifications` table alongside business logic in same transaction, then sent asynchronously by separate worker. _Rationale:_ Guarantees notification delivery even if webhook endpoint is temporarily unavailable.

