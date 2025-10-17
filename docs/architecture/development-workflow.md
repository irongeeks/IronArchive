## Development Workflow

### Local Development Setup

**Prerequisites:**

```bash
# Install required tools
brew install go@1.24              # Go 1.24
brew install node@20              # Node.js 20 LTS
brew install docker               # Docker Desktop
brew install golang-migrate       # Database migrations CLI
```

**Initial Setup:**

```bash
# Clone repository
git clone https://github.com/yourusername/ironarchive.git
cd ironarchive

# Run setup script
./scripts/setup-dev.sh

# Start services
docker-compose -f docker/docker-compose.dev.yml up -d
```

**Development Commands:**

```bash
# Start all services (backend, frontend, PostgreSQL, Redis, Meilisearch)
make dev

# Start frontend only (with HMR)
cd frontend && npm run dev

# Start backend only (with hot-reload via air)
cd backend && air

# Run tests
make test                         # All tests
make test-backend                 # Backend tests only
make test-frontend                # Frontend tests only
make test-e2e                     # E2E tests

# Database operations
make migrate-up                   # Apply migrations
make migrate-down                 # Rollback migrations
make migrate-create NAME=<name>   # Create new migration

# Linting
make lint                         # Run all linters
make lint-backend                 # Go linters (golangci-lint)
make lint-frontend                # ESLint + Prettier
```

### Environment Configuration

**Frontend (.env.local):**

```bash
VITE_API_URL=http://localhost:3000/api/v1
VITE_APP_NAME=IronArchive
VITE_ENABLE_DEBUG=true
```

**Backend (.env):**

```bash
# Server
PORT=3000
JWT_SECRET=your-secret-key-here

# Database
DATABASE_URL=postgres://ironarchive:password@localhost:5432/ironarchive?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# Meilisearch
MEILISEARCH_URL=http://localhost:7700
MEILISEARCH_API_KEY=masterKey

# File Storage
STORAGE_PATH=/var/lib/ironarchive/archive

# Microsoft Graph API (for development testing)
DEV_AZURE_TENANT_ID=
DEV_AZURE_APP_ID=
DEV_AZURE_APP_SECRET=

# SMTP (optional for local dev)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASS=

# Logging
LOG_LEVEL=debug
```

**Shared (.env for docker-compose):**

```bash
POSTGRES_USER=ironarchive
POSTGRES_PASSWORD=password
POSTGRES_DB=ironarchive
```

