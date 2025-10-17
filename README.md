# IronArchive

**IronArchive** is a self-hosted Microsoft 365 email archiving platform designed for long-term email retention, compliance, and legal discovery. Built with scalability and security in mind, IronArchive provides robust email archiving capabilities for organizations of all sizes.

## Features

- **M365 Integration**: Seamless integration with Microsoft 365 via Graph API
- **Multi-Tenant Support**: Manage multiple organizations from a single platform
- **Full-Text Search**: Fast, typo-tolerant search powered by Meilisearch
- **Legal Hold**: Immutable email retention for compliance and eDiscovery
- **Role-Based Access Control**: Granular permissions for users and administrators
- **Background Processing**: Asynchronous email ingestion and indexing
- **Self-Hosted**: Complete control over your data and infrastructure

## Tech Stack

- **Backend**: Go 1.24 with Fiber v3
- **Frontend**: SvelteKit 2.x with TypeScript and TailwindCSS 4.x
- **Database**: PostgreSQL 16
- **Cache & Queue**: Redis 7
- **Search Engine**: Meilisearch 1.6
- **Infrastructure**: Docker Compose for development and deployment

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: 1.24 or later ([Download](https://go.dev/dl/))
- **Node.js**: 20.x or later ([Download](https://nodejs.org/))
- **Docker**: 24.x or later ([Download](https://www.docker.com/products/docker-desktop/))
- **Docker Compose**: 2.x or later (included with Docker Desktop)
- **Git**: For version control

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/irongeeks/IronArchive.git
cd IronArchive
```

### 2. Set Up Environment Variables

Copy the example environment file and configure your local settings:

```bash
cp .env.example .env
# Edit .env with your preferred text editor
```

### 3. Start Docker Services

Start PostgreSQL, Redis, and Meilisearch:

```bash
make docker-up
```

### 4. Run Backend Server

In a new terminal, start the backend:

```bash
cd backend
go run ./cmd/server
```

### 5. Run Frontend Development Server

In another terminal, start the frontend:

```bash
cd frontend
npm install
npm run dev
```

The application will be available at:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080

## Development

### Available Make Targets

```bash
make dev            # Start full development environment
make build          # Build backend binary and frontend assets
make test           # Run all tests (backend + frontend)
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make migrate-up     # Run database migrations
make migrate-down   # Roll back database migrations
make migrate-create # Create a new migration file
make migrate-status # Check current migration version
make clean          # Clean build artifacts
```

### Database Tests

The backend schema tests require the following prerequisites:

- **PostgreSQL**: Running at `postgres://ironarchive:ironarchive_password@localhost:5432/ironarchive?sslmode=disable` (started via `make docker-up`)
- **migrate CLI**: The `golang-migrate` command-line tool must be available on your PATH
  - **macOS**: `brew install golang-migrate`
  - **Linux**: Download from [golang-migrate releases](https://github.com/golang-migrate/migrate/releases)
  - **Verify installation**: `migrate -version`

Once prerequisites are ready, run the database tests:

```bash
cd backend
go test ./internal/database/...
```

**Note**: The test suite automatically runs migrations up and down for each test, ensuring a clean database state.

### Project Structure

```
ironarchive/
├── backend/         # Go backend application
│   ├── cmd/         # Application entrypoints
│   ├── internal/    # Internal packages
│   └── pkg/         # Public packages
├── frontend/        # SvelteKit frontend application
│   ├── src/         # Source code
│   └── static/      # Static assets
├── docker/          # Docker configurations
├── migrations/      # Database migrations
├── docs/            # Project documentation
└── scripts/         # Utility scripts
```

## Architecture

IronArchive follows a monorepo structure with a clear separation between backend and frontend concerns. The backend is a monolithic Go application using the repository pattern for data access, while the frontend is a SvelteKit application with server-side rendering capabilities.

For detailed architecture documentation, see:
- [Architecture Overview](docs/architecture/)
- [Backend Architecture](docs/architecture/backend-architecture.md)
- [Frontend Architecture](docs/architecture/frontend-architecture.md)

## Contributing

We welcome contributions! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

For questions, issues, or feature requests, please open an issue on [GitHub](https://github.com/irongeeks/IronArchive/issues).

---

Built with ❤️ by IronGeeks
