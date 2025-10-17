# Source Tree Reference

This document provides a detailed reference of the IronArchive source code organization, explaining the purpose and contents of each directory and key files within the monorepo.

## Root Directory Structure

```
ironarchive/
├── backend/              # Go backend application
├── frontend/             # SvelteKit frontend application
├── docker/               # Docker configurations and compose files
├── migrations/           # Database migration files
├── docs/                 # Project documentation
├── scripts/              # Utility and automation scripts
├── .bmad-core/          # BMAD agent configurations and workflows
├── .claude/             # Claude Code configurations
├── .env.example         # Environment variables template
├── .gitignore           # Git ignore patterns
├── Makefile             # Top-level build and dev commands
├── LICENSE              # MIT License
├── README.md            # Project overview and quick start
├── CONTRIBUTING.md      # Contribution guidelines
└── AGENTS.md            # AI agent documentation
```

## Backend (`/backend`)

The Go backend follows a standard Go project layout with clear separation between application entry points, internal packages, and public APIs.

```
backend/
├── cmd/                 # Application entry points
│   └── server/
│       └── main.go      # HTTP server entry point
├── internal/            # Private application code
│   ├── config/
│   │   └── config.go    # Configuration loading and validation
│   ├── database/
│   │   ├── postgres.go  # PostgreSQL connection and queries
│   │   ├── redis.go     # Redis client and cache operations
│   │   └── meilisearch.go # Meilisearch client and indexing
│   └── utils/
│       └── logger.go    # Structured logging utilities
├── backend/             # Legacy/alternative backend structure (to be consolidated)
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── bin/                 # Compiled binaries (gitignored)
├── go.mod               # Go module definition
├── go.sum               # Go dependency checksums
└── server               # Compiled server binary (gitignored)
```

### Backend Architecture Patterns

- **Entry Points** (`cmd/`): Each subdirectory represents a deployable binary
- **Internal Packages** (`internal/`): Private code not importable by external projects
  - `config/`: Environment and configuration management
  - `database/`: Database clients and connection management
  - `utils/`: Shared utilities (logging, validation, etc.)
- **Public APIs** (`pkg/`): Reusable packages that can be imported externally

### Key Backend Files

- `cmd/server/main.go`: HTTP server initialization, middleware setup, route registration
- `internal/config/config.go`: Environment variable parsing, config validation
- `internal/database/*.go`: Database connection pools, query builders
- `go.mod`: Dependency management (Fiber v3, PostgreSQL drivers, Redis client)

## Frontend (`/frontend`)

SvelteKit application with TypeScript, TailwindCSS 4.x, and component-based architecture.

```
frontend/
├── src/                 # Source code
│   ├── lib/             # Shared libraries and components
│   │   ├── assets/      # Images, icons, fonts
│   │   │   └── favicon.svg
│   │   └── index.ts     # Library exports
│   ├── routes/          # File-based routing
│   │   ├── +page.svelte       # Home page
│   │   └── +layout.svelte     # Root layout
│   ├── app.html         # HTML template
│   ├── app.css          # Global styles and Tailwind imports
│   └── app.d.ts         # TypeScript declarations
├── static/              # Static assets served as-is
├── node_modules/        # NPM dependencies (gitignored)
├── package.json         # NPM dependencies and scripts
├── package-lock.json    # NPM dependency lock
├── svelte.config.js     # SvelteKit configuration
├── tailwind.config.js   # TailwindCSS configuration
├── vite.config.ts       # Vite bundler configuration
├── tsconfig.json        # TypeScript configuration
└── vitest.config.ts     # Vitest test runner configuration
```

### Frontend Architecture Patterns

- **File-Based Routing** (`routes/`): SvelteKit automatically creates routes based on file structure
  - `+page.svelte`: Page component
  - `+layout.svelte`: Layout wrapper
  - `+server.ts`: API endpoint
  - `+page.server.ts`: Server-side data loading
- **Component Library** (`lib/`): Reusable Svelte components, utilities, and stores
- **Static Assets** (`static/`): Served directly without processing

### Key Frontend Files

- `src/app.html`: Root HTML template with `%sveltekit.head%` and `%sveltekit.body%` placeholders
- `src/app.css`: Global styles, TailwindCSS directives
- `src/routes/+page.svelte`: Landing page component
- `svelte.config.js`: Adapter configuration, preprocessors
- `vite.config.ts`: Development server, build optimizations, plugins

## Docker (`/docker`)

Docker configurations for containerized deployment and local development.

```
docker/
├── Dockerfile.backend       # Multi-stage Go build
├── Dockerfile.frontend      # Node build + static serve
├── docker-compose.yml       # Production-like deployment
└── docker-compose.dev.yml   # Dev environment with volume mounts
```

### Docker Services

Development environment includes:
- **PostgreSQL 16**: Primary database on port 5432
- **Redis 7**: Cache and queue on port 6379
- **Meilisearch 1.6**: Search engine on port 7700

## Migrations (`/migrations`)

Database schema migrations using golang-migrate or similar tool.

```
migrations/
├── 000001_initial_schema.up.sql      # Initial schema creation
├── 000001_initial_schema.down.sql    # Initial schema rollback
├── 000002_add_legal_hold.up.sql      # Add legal hold tables
└── 000002_add_legal_hold.down.sql    # Rollback legal hold
```

### Migration Naming Convention

- Format: `{version}_{description}.{up|down}.sql`
- Versions are sequential, zero-padded (000001, 000002, etc.)
- Each migration has an `up` (apply) and `down` (rollback) version

## Documentation (`/docs`)

Comprehensive project documentation organized by type.

```
docs/
├── architecture/        # Sharded architecture documentation
│   ├── README.md                      # Architecture overview
│   ├── index.md                       # Architecture index
│   ├── introduction.md                # Introduction to IronArchive
│   ├── high-level-architecture.md     # System overview and diagrams
│   ├── components.md                  # Component descriptions
│   ├── backend-architecture.md        # Backend design and patterns
│   ├── frontend-architecture.md       # Frontend design and patterns
│   ├── data-models.md                 # Domain models and entities
│   ├── database-schema.md             # Database table structures
│   ├── api-specification.md           # REST API endpoints
│   ├── external-apis.md               # Third-party integrations
│   ├── core-workflows.md              # Key user flows
│   ├── deployment-architecture.md     # Infrastructure and deployment
│   ├── tech-stack.md                  # Technology decisions
│   ├── coding-standards.md            # Code style and conventions
│   ├── testing-strategy.md            # Test approach and coverage
│   ├── error-handling.md              # Error handling patterns
│   ├── security-performance.md        # Security and performance
│   ├── monitoring.md                  # Observability and logging
│   ├── development-workflow.md        # Dev environment setup
│   ├── project-structure.md           # High-level structure overview
│   └── source-tree.md                 # This document
├── prd/                 # Product requirements (sharded)
│   └── (PRD sections)
├── qa/                  # Quality assurance documentation
│   ├── gates/           # Quality gates and checklists
│   └── assessments/     # QA assessment reports
└── stories/             # User stories and epics
    └── *.story.md       # Individual story files
```

## Scripts (`/scripts`)

Automation scripts for development, testing, and deployment.

```
scripts/
├── setup-dev.sh         # Initial development environment setup
├── migrate.sh           # Database migration runner
├── build-docker.sh      # Docker image builder
└── (additional utility scripts)
```

## BMAD Core (`.bmad-core`)

BMAD agent system configurations and reusable workflows.

```
.bmad-core/
├── core-config.yaml     # Project-wide BMAD configuration
├── tasks/               # Reusable task workflows
├── templates/           # Document templates
├── checklists/          # Quality and architecture checklists
└── utils/               # Utility scripts and helpers
```

## Claude Code (`.claude`)

Claude Code agent configurations and custom commands.

```
.claude/
├── commands/            # Custom slash commands
└── config.json          # Claude Code settings
```

## Environment Configuration

### `.env.example`

Template for environment variables required by the application:

```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/ironarchive
REDIS_URL=redis://localhost:6379
MEILISEARCH_URL=http://localhost:7700
MEILISEARCH_MASTER_KEY=masterKey

# Microsoft Graph API
MS_GRAPH_CLIENT_ID=your-client-id
MS_GRAPH_CLIENT_SECRET=your-client-secret
MS_GRAPH_TENANT_ID=your-tenant-id

# Application
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug
```

## Build System

### Root `Makefile`

Top-level commands for common development tasks:

```makefile
dev              # Start full development environment
build            # Build backend binary and frontend assets
test             # Run all tests (backend + frontend)
docker-up        # Start Docker services
docker-down      # Stop Docker services
migrate-up       # Apply database migrations
migrate-down     # Rollback database migrations
clean            # Clean build artifacts
```

## File Conventions

### Backend (Go)

- **Package naming**: lowercase, single word (e.g., `config`, `database`)
- **File naming**: lowercase with underscores (e.g., `postgres.go`, `meilisearch.go`)
- **Test files**: `*_test.go` alongside implementation
- **Internal packages**: Not importable outside the module

### Frontend (SvelteKit)

- **Component naming**: PascalCase for reusable components (e.g., `Button.svelte`)
- **Route files**: Prefix with `+` (e.g., `+page.svelte`, `+layout.svelte`)
- **TypeScript**: Strict mode enabled, no implicit `any`
- **Styles**: TailwindCSS utility classes, scoped `<style>` blocks

### Documentation

- **Markdown**: GitHub-flavored markdown (GFM)
- **File naming**: kebab-case (e.g., `high-level-architecture.md`)
- **Sharding**: Large docs split into focused, single-topic files

## Git Conventions

### Ignored Files (`.gitignore`)

- Build artifacts: `backend/bin/`, `backend/server`, `frontend/build/`
- Dependencies: `node_modules/`, `vendor/`
- Environment: `.env`, `.env.local`
- IDE files: `.vscode/`, `.idea/`, `*.swp`
- Logs: `*.log`, `.ai/debug-log.md`

## Navigation Tips

### Finding Backend Code

1. **Entry point**: Start at `backend/cmd/server/main.go`
2. **Configuration**: Check `backend/internal/config/config.go`
3. **Database logic**: Look in `backend/internal/database/`
4. **Utilities**: Browse `backend/internal/utils/`

### Finding Frontend Code

1. **Entry point**: Start at `frontend/src/routes/+page.svelte`
2. **Root layout**: Check `frontend/src/routes/+layout.svelte`
3. **Components**: Browse `frontend/src/lib/`
4. **Styles**: See `frontend/src/app.css`

### Finding Documentation

1. **Architecture overview**: `docs/architecture/README.md` or `index.md`
2. **Tech decisions**: `docs/architecture/tech-stack.md`
3. **API reference**: `docs/architecture/api-specification.md`
4. **Development setup**: `docs/architecture/development-workflow.md`

## Planned Expansions

As the project grows, expect these additions:

```
backend/internal/
├── handlers/            # HTTP request handlers
├── services/            # Business logic services
├── models/              # Domain models and DTOs
├── repositories/        # Data access layer
└── middleware/          # HTTP middleware

frontend/src/
├── lib/
│   ├── components/      # Reusable UI components
│   ├── stores/          # Svelte stores (state management)
│   ├── api/             # API client functions
│   └── utils/           # Utility functions
└── routes/
    ├── auth/            # Authentication pages
    ├── dashboard/       # Dashboard pages
    ├── search/          # Search interface
    └── admin/           # Admin panel
```

## Cross-Reference

For more information, see:
- [Project Structure Overview](./project-structure.md) - High-level structure
- [Backend Architecture](./backend-architecture.md) - Backend design patterns
- [Frontend Architecture](./frontend-architecture.md) - Frontend design patterns
- [Development Workflow](./development-workflow.md) - Setup and dev practices
- [Tech Stack](./tech-stack.md) - Technology decisions and rationale

---

**Last Updated**: 2025-10-17
**Maintained By**: Architecture Team
**Status**: Active Development
