## Coding Standards

### Critical Fullstack Rules

- **Type Sharing:** Always define shared types in a common location. Backend uses Go structs, frontend uses TypeScript interfaces. Keep them in sync manually or generate TypeScript types from Go structs using tools like `typescriptify`.

- **API Calls:** Never make direct HTTP calls from frontend components. Always use the service layer (`lib/services/*`) which handles auth, error handling, and retries.

- **Environment Variables:** Access only through config objects, never `process.env` or `os.Getenv()` directly in business logic. Use dedicated config modules.

- **Error Handling:** All API routes must use the standard error handler. Backend returns consistent JSON error format. Frontend displays user-friendly error messages.

- **State Updates:** Never mutate state directly in Svelte components. Use `$store.update()` or reactive assignments (`$: value = ...`).

- **Database Queries:** Always use repository pattern. Never write raw SQL in handlers or services. Apply tenant filtering in repository layer.

- **RBAC Enforcement:** Check user role at both API gateway (middleware) AND repository layer (defense in depth).

### Naming Conventions

| Element | Frontend | Backend | Example |
|---------|----------|---------|---------|
| Components | PascalCase | - | `UserProfile.svelte` |
| Hooks | camelCase with 'use' | - | `useAuth.ts` |
| API Routes | - | kebab-case | `/api/user-profile` |
| Database Tables | - | snake_case | `user_profiles` |
| Go Functions | - | PascalCase (exported) / camelCase (private) | `CreateTenant()`, `validateEmail()` |
| TypeScript Functions | camelCase | - | `formatDate()` |
| Constants | SCREAMING_SNAKE_CASE | SCREAMING_SNAKE_CASE | `MAX_FILE_SIZE` |

## Error Handling Strategy
