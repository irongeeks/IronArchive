## Tech Stack

The following technology stack is the single source of truth for IronArchive development. All AI agents and developers MUST use these exact versions and technologies.

### Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|-----------|---------|---------|-----------|
| Frontend Language | TypeScript | 5.3+ | Type-safe frontend development | Prevents runtime errors, enhances IDE support, aligns with Svelte ecosystem |
| Frontend Framework | SvelteKit | 2.x | Fullstack framework with SSR/SPA hybrid | Best-in-class developer experience, minimal bundle size, native reactivity without virtual DOM |
| UI Component Library | Shadcn-Svelte | Latest | Accessible, customizable UI components | Built on Radix primitives, fully themeable for whitelabeling, WCAG 2.1 AA compliant |
| State Management | Svelte Stores | Built-in | Reactive state management | Native to Svelte, zero external dependencies, simple API for global state |
| Backend Language | Go | 1.24 | High-performance backend services | Superior concurrency (goroutines), compiled binary simplifies deployment, strong Graph API SDK |
| Backend Framework | Fiber | v3 | Express-like web framework | Built on fasthttp (fastest Go HTTP), minimal overhead, middleware ecosystem, familiar API |
| API Style | REST | OpenAPI 3.0 | HTTP API specification | Widely understood, simple client integration, OpenAPI enables code generation |
| Database | PostgreSQL | 16+ | Primary relational database | ACID compliance, JSON support, pg_cron for scheduling, battle-tested reliability |
| Cache | Redis | 7+ | Caching and job queue | High-performance in-memory store, native support in asynq, Pub/Sub for real-time updates |
| File Storage | Local Filesystem | N/A | Email body and attachment storage | Zero additional service cost, simple backup/restore, NFS/SAN compatible for scale |
| Authentication | JWT + bcrypt | golang-jwt/jwt v5, bcrypt cost 12 | Secure authentication system | Stateless auth, industry-standard token format, adaptive cost factor for future-proofing |
| Search Engine | Meilisearch | 1.6+ | Full-text search with typo tolerance | Sub-200ms search, instant indexing, 10x simpler than Elasticsearch, low memory footprint |
| Frontend Testing | Vitest + Testing Library | Vitest 1.x, Svelte Testing Library 4.x | Component and unit testing | Vite-native test runner (faster), idiomatic Svelte testing patterns |
| Backend Testing | Go Testing + Testify | Standard library, Testify 1.9+ | Unit and integration testing | Native Go testing, Testify adds assertions and mocking, no framework overhead |
| E2E Testing | Playwright | 1.40+ | Browser automation and E2E tests | Cross-browser support, excellent TypeScript integration, fast parallel execution |
| Build Tool | Vite | 5.x | Frontend build and dev server | Fastest HMR, native ESM, optimized production builds |
| Bundler | Rollup (via Vite) | 4.x | Production bundling | Tree-shaking, code splitting, minimal bundle size |
| IaC Tool | Docker Compose | 2.x | Infrastructure as code | Simple declarative YAML, sufficient for single-server deployments, Docker-native |
| CI/CD | GitHub Actions | N/A | Continuous integration and deployment | Free for open-source, excellent ecosystem, matrix builds for multi-arch Docker images |
| Monitoring | Prometheus + Grafana | Prometheus 2.48+, Grafana 10.x (Optional) | Metrics collection and visualization | Industry standard, optional deployment (not required), extensive integrations |
| Logging | Structured JSON Logging | Zap (Go), Pino (JS) | Centralized structured logging | Machine-parseable, grep-friendly, supports log aggregation tools |
| CSS Framework | TailwindCSS | 4.x | Utility-first CSS framework | Rapid prototyping, theme system for whitelabeling, purged CSS for minimal bundle |

