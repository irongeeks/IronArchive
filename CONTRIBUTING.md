# Contributing to IronArchive

Thank you for your interest in contributing to IronArchive! We welcome contributions from the community and are grateful for your support.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and collaborative environment for all contributors.

## How to Contribute

### Reporting Issues

If you encounter a bug or have a feature request:

1. Check if the issue already exists in [GitHub Issues](https://github.com/irongeeks/IronArchive/issues)
2. If not, create a new issue with a clear title and description
3. For bugs, include:
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment details (OS, Go version, Node version)
   - Relevant logs or error messages

### Submitting Pull Requests

1. **Fork the repository** and create a new branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards (see below)

3. **Test your changes**:
   ```bash
   make test
   ```

4. **Commit your changes** with clear, descriptive commit messages:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** against the `main` branch with:
   - Clear description of changes
   - Reference to related issues (if applicable)
   - Screenshots (for UI changes)

## Development Workflow

### Setting Up Your Development Environment

1. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/IronArchive.git
   cd IronArchive
   ```

2. Install dependencies:
   ```bash
   # Backend dependencies are managed by Go modules
   cd backend && go mod download

   # Frontend dependencies
   cd ../frontend && npm install
   ```

3. Start Docker services:
   ```bash
   make docker-up
   ```

4. Run the application:
   ```bash
   # Terminal 1: Backend
   cd backend && go run ./cmd/server

   # Terminal 2: Frontend
   cd frontend && npm run dev
   ```

## Coding Standards

### Backend (Go)

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for code formatting
- Exported functions must have doc comments
- Use meaningful variable and function names
- Keep functions small and focused
- Write tests for new functionality

### Frontend (TypeScript/Svelte)

- Follow TypeScript best practices
- Use ESLint and Prettier for code formatting
- Component names use PascalCase (e.g., `UserProfile.svelte`)
- Functions use camelCase (e.g., `formatDate()`)
- Use Svelte stores for global state management
- Write component tests using Vitest

### Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Go Functions (exported) | PascalCase | `CreateTenant()` |
| Go Functions (private) | camelCase | `validateEmail()` |
| TypeScript Functions | camelCase | `formatDate()` |
| Svelte Components | PascalCase | `UserProfile.svelte` |
| API Routes | kebab-case | `/api/user-profile` |
| Database Tables | snake_case | `user_profiles` |
| Constants | SCREAMING_SNAKE_CASE | `MAX_FILE_SIZE` |

## Testing

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm run test
```

### End-to-End Tests

```bash
npm run test:e2e
```

## Documentation

- Update documentation for any new features or API changes
- Add inline code comments for complex logic
- Update README.md if prerequisites or setup steps change

## Git Commit Messages

Write clear, descriptive commit messages:

```
Add user authentication with JWT

- Implement JWT token generation
- Add middleware for protected routes
- Create login/logout endpoints
- Add tests for auth service
```

**Format:**
- First line: Brief summary (50 chars or less)
- Blank line
- Detailed description (if needed)

## Need Help?

If you have questions or need assistance:

- Open a [Discussion](https://github.com/irongeeks/IronArchive/discussions)
- Join our community chat (coming soon)
- Email: support@irongeeks.de

Thank you for contributing to IronArchive!
