# Epic 1: Foundation & Core Infrastructure

**Goal:** Establish the foundational project structure, development environment, core services, database schema, and authentication system. Deliver a functional setup wizard that creates the first MSP Admin account, demonstrating end-to-end capability from database to frontend. This epic enables all subsequent development while delivering a deployable, testable application with basic authentication.

## Story 1.1: Project Scaffolding and Repository Structure

As a developer,
I want the project repository initialized with proper structure and tooling,
so that I have a productive development environment from day one.

**Acceptance Criteria:**

1. Repository created with monorepo structure (`/backend`, `/frontend`, `/docker`, `/docs`, `/migrations`, `/scripts`)
2. Backend Go module initialized with `go.mod`, folder structure (`/cmd`, `/internal`, `/pkg`)
3. Frontend SvelteKit project initialized with TypeScript, Tailwind CSS configured
4. `.gitignore` configured for Go, Node.js, and IDE files
5. `README.md` includes project description, setup instructions, and architecture overview
6. `LICENSE` file (MIT) and `CONTRIBUTING.md` created
7. Docker Compose file includes PostgreSQL 16, Redis 7, Meilisearch 1.6 service definitions
8. All services start successfully with `docker-compose up`
9. Backend can connect to PostgreSQL, Redis, and Meilisearch with connection validation logging

## Story 1.2: Database Schema Design and Migration System

As a developer,
I want the core database schema designed and migration system configured,
so that I can evolve the schema safely throughout development.

**Acceptance Criteria:**

1. Migration tool (golang-migrate) configured with up/down migrations
2. Initial migration creates core tables: `users`, `tenants`, `mailboxes`, `emails`, `attachments`, `jobs`, `audit_logs`
3. Users table includes: id (UUID), email, password_hash, role (MSP Admin/Tenant Admin/User), mfa_secret, created_at, updated_at
4. Tenants table includes: id (UUID), name, azure_tenant_id, azure_app_credentials (encrypted), retention_policy_days, legal_hold (boolean), created_at
5. Mailboxes table includes: id (UUID), tenant_id (FK), email_address, display_name, mailbox_type, sync_enabled, last_sync_at, created_at
6. Emails table includes: id (UUID), mailbox_id (FK), message_id, subject, sender, recipients, sent_at, body_text, body_html, has_attachments, size_bytes, indexed_at
7. Attachments table includes: id (UUID), email_id (FK), filename, content_type, size_bytes, sha256_hash, file_path
8. Jobs table includes: id (UUID), type, status, tenant_id (FK), mailbox_id (FK nullable), started_at, completed_at, error_message, metadata (JSONB)
9. Audit_logs table includes: id (UUID), user_id (FK), action, ip_address, timestamp, details (JSONB), immutable (enforced by trigger)
10. Foreign key constraints properly defined with ON DELETE CASCADE where appropriate
11. Indexes created on frequently queried columns (tenant_id, mailbox_id, email message_id, sent_at)
12. Migration runs successfully on fresh database with `make migrate-up` command
13. Down migration successfully rolls back all changes

## Story 1.3: Backend API Foundation with Fiber Framework

As a developer,
I want the backend API server scaffolded with routing, middleware, and error handling,
so that I can build endpoints efficiently with consistent patterns.

**Acceptance Criteria:**

1. Fiber v3 application initialized with structured routing
2. Environment configuration loaded from `.env` file (database URL, Redis URL, JWT secret, server port)
3. Middleware configured: CORS, request logging, panic recovery, request ID generation
4. Database connection pool established using pgx with health check endpoint (`GET /health`)
5. Redis client initialized with connection validation
6. Meilisearch client initialized with connection validation
7. Structured logging configured (JSON format to stdout) with log levels (DEBUG, INFO, WARN, ERROR)
8. Error handling middleware returns consistent JSON error responses with status codes
9. API versioning structure (`/api/v1/...`) established
10. Server starts successfully and responds to health check with status 200 and service connectivity status
11. Unit tests for health check endpoint pass
12. Server gracefully shuts down on SIGTERM/SIGINT, closing database connections

## Story 1.4: Authentication System (JWT, bcrypt, RBAC Middleware)

As a developer,
I want authentication and authorization implemented,
so that API endpoints can be secured with role-based access control.

**Acceptance Criteria:**

1. Password hashing implemented using bcrypt with cost factor 12
2. JWT token generation implemented with 15-minute access token expiration and 7-day refresh token expiration
3. JWT middleware validates tokens and extracts user claims (user_id, role, email)
4. RBAC middleware created with role checks (RequireMSPAdmin, RequireTenantAdmin, RequireUser)
5. Authentication endpoints implemented: `POST /api/v1/auth/login` (returns access + refresh tokens), `POST /api/v1/auth/refresh` (returns new access token), `POST /api/v1/auth/logout` (invalidates refresh token)
6. Login endpoint validates email/password, returns tokens and user profile on success, returns 401 on failure
7. Refresh endpoint validates refresh token, returns new access token, returns 401 on invalid token
8. JWT tokens include claims: user_id, email, role, issued_at, expires_at
9. Tokens stored as HTTP-only, secure cookies (SameSite=Strict) in addition to JSON response
10. Unit tests cover: password hashing/verification, JWT generation/validation, role-based access checks
11. Integration tests cover: successful login, failed login (wrong password), token refresh, protected endpoint access with/without valid token

## Story 1.5: Setup Wizard Backend API

As a developer,
I want the setup wizard API endpoint created,
so that the first MSP Admin can be initialized through the frontend.

**Acceptance Criteria:**

1. Setup status endpoint implemented: `GET /api/v1/setup/status` returns `{"setup_complete": boolean}`
2. Setup complete check: Returns true if any user exists in database with role=MSP Admin
3. Setup endpoint implemented: `POST /api/v1/setup/initialize` accepts `{email, password, display_name}`
4. Setup endpoint validates: password meets complexity requirements (12+ chars, uppercase, lowercase, number, special char), email is valid format, display_name is non-empty
5. Setup endpoint creates first MSP Admin user with hashed password, role=MSP Admin
6. Setup endpoint returns JWT tokens and user profile on success
7. Setup endpoint returns 409 Conflict if setup already completed
8. Setup endpoint returns 400 Bad Request with validation errors for invalid input
9. Transaction ensures atomicity (user creation succeeds or fails completely)
10. Unit tests cover validation logic, password complexity checks
11. Integration tests cover: successful setup, duplicate setup attempt (fails), invalid input handling

## Story 1.6: Frontend Foundation (SvelteKit, TailwindCSS, Routing)

As a developer,
I want the frontend application scaffolded with routing, layout, and styling foundation,
so that I can build UI features efficiently with consistent design.

**Acceptance Criteria:**

1. SvelteKit 2.x project initialized with TypeScript, Vite build tool configured
2. TailwindCSS 4.x configured with custom theme setup (placeholder for theme system)
3. Shadcn-Svelte component library installed and configured
4. Root layout component created with navigation structure (header, sidebar placeholder, main content area)
5. Routing structure established: `/setup`, `/login`, `/dashboard`, `/tenants`, `/search`, `/settings`, `/profile`
6. Authentication store created (Svelte store) tracking: user profile, access token, isAuthenticated flag
7. API client utility created with fetch wrapper handling: JWT token injection, error responses, token refresh on 401
8. Protected route logic implemented: redirects to `/login` if unauthenticated, redirects to `/dashboard` if authenticated user accesses `/login`
9. Setup check implemented on app load: if setup incomplete, redirect all routes to `/setup`
10. Responsive layout works on mobile (320px), tablet (768px), desktop (1920px+)
11. Navigation menu shows/hides based on authentication state
12. Development server starts with `npm run dev`, hot module replacement works

## Story 1.7: Setup Wizard Frontend UI

As an installer,
I want a clean, guided setup wizard,
so that I can create the first MSP Admin account and initialize IronArchive.

**Acceptance Criteria:**

1. Setup wizard page (`/setup`) displays: application logo, welcome message, single-step form
2. Form includes fields: Email (text input), Password (password input with show/hide toggle), Confirm Password (password input), Display Name (text input)
3. Password strength indicator displays (weak/medium/strong) based on complexity
4. Password complexity requirements shown: 12+ characters, uppercase, lowercase, number, special character
5. Form validation: email format, passwords match, password complexity, display name non-empty
6. Submit button disabled until all validation passes
7. Loading state displayed during API call (button shows spinner, form disabled)
8. Success: Stores JWT tokens in auth store, redirects to `/dashboard`
9. Error: Displays error message above form (e.g., "Setup already completed", "Invalid input")
10. Visual design: clean, minimal, centered card layout with IronArchive branding
11. Accessibility: keyboard navigation works, form labels properly associated, ARIA attributes set
12. Responsive design: works on mobile, tablet, desktop

## Story 1.8: Login Page Frontend UI

As a user,
I want a login page,
so that I can authenticate and access the IronArchive system.

**Acceptance Criteria:**

1. Login page (`/login`) displays: application logo, login form, "Forgot Password?" link (placeholder for future)
2. Form includes fields: Email (text input), Password (password input with show/hide toggle), "Remember Me" checkbox (optional)
3. Form validation: email format, password non-empty
4. Submit button disabled until validation passes
5. Loading state displayed during API call
6. Success: Stores JWT tokens in auth store, redirects to `/dashboard` (or originally requested protected route)
7. Error: Displays error message ("Invalid email or password", network errors)
8. "Remember Me" extends refresh token expiration to 30 days (if checked)
9. Logout functionality implemented: clears auth store, calls `/api/v1/auth/logout`, redirects to `/login`
10. Accessibility: keyboard navigation, form labels, ARIA attributes
11. Responsive design: mobile, tablet, desktop
12. Visual design: consistent with setup wizard styling

## Story 1.9: Basic Dashboard Page (Empty State)

As a developer,
I want a basic dashboard page with empty state,
so that authenticated users have a landing page and I can build upon it in future epics.

**Acceptance Criteria:**

1. Dashboard page (`/dashboard`) displays for authenticated users
2. Empty state shows: "Welcome to IronArchive" heading, "Add Your First Tenant" call-to-action button (placeholder, not functional yet)
3. Dashboard includes top navigation bar with: IronArchive logo, user profile dropdown (shows display name, email, Logout button)
4. User profile dropdown Logout button calls logout API, clears tokens, redirects to `/login`
5. Dashboard includes sidebar navigation (placeholder links for: Dashboard, Tenants, Search, Settings, Profile)
6. Main content area displays empty state illustration/icon and instructional text
7. Page title set correctly ("Dashboard - IronArchive")
8. Loading state displays while fetching initial data (if applicable)
9. Error handling: if API calls fail, show error message with retry button
10. Accessibility: navigation keyboard-accessible, focus management, ARIA landmarks
11. Responsive design: mobile, tablet, desktop
12. Visual design: establishes design system patterns for future pages

---
