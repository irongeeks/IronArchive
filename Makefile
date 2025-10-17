.PHONY: help dev build test clean docker-up docker-down migrate-up migrate-down backend-dev frontend-dev

# Default target - show help
help:
	@echo "IronArchive - Available Make Targets"
	@echo "===================================="
	@echo "  dev            - Start full development environment (Docker + Backend + Frontend)"
	@echo "  build          - Build backend binary and frontend assets"
	@echo "  test           - Run all tests (backend + frontend)"
	@echo "  clean          - Clean build artifacts and dependencies"
	@echo ""
	@echo "  docker-up      - Start Docker services (PostgreSQL, Redis, Meilisearch)"
	@echo "  docker-down    - Stop Docker services"
	@echo ""
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Roll back database migrations"
	@echo ""
	@echo "  backend-dev    - Start backend development server"
	@echo "  frontend-dev   - Start frontend development server"

# Start full development environment
dev: docker-up
	@echo "Starting IronArchive development environment..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Frontend will be available at http://localhost:5173"
	@echo ""
	@echo "Run 'make backend-dev' and 'make frontend-dev' in separate terminals"

# Build backend and frontend
build:
	@echo "Building backend..."
	cd backend && go build -o bin/server ./cmd/server
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Build complete!"

# Run all tests
test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running frontend tests..."
	cd frontend && npm run test
	@echo "All tests passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf backend/tmp
	rm -rf frontend/.svelte-kit
	rm -rf frontend/build
	@echo "Clean complete!"

# Docker services
docker-up:
	@echo "Starting Docker services..."
	cd docker && docker-compose up -d
	@echo "Docker services started!"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"
	@echo "  Meilisearch: localhost:7700"

docker-down:
	@echo "Stopping Docker services..."
	cd docker && docker-compose down
	@echo "Docker services stopped!"

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	@echo "TODO: Implement migration command"

migrate-down:
	@echo "Rolling back database migrations..."
	@echo "TODO: Implement rollback command"

# Development servers
backend-dev:
	@echo "Starting backend development server..."
	cd backend && go run ./cmd/server

frontend-dev:
	@echo "Starting frontend development server..."
	cd frontend && npm run dev
