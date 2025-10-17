## Backend Architecture

### Service Architecture (Monolithic Go Binary)

The backend is organized as a single Go binary with distinct internal packages:

```
/backend
├── cmd/
│   └── server/
│       └── main.go              # Application entrypoint
├── internal/
│   ├── api/                     # HTTP handlers
│   │   ├── middleware/          # Auth, CORS, logging middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   └── rbac.go
│   │   ├── handlers/            # Route handlers
│   │   │   ├── auth.go
│   │   │   ├── tenants.go
│   │   │   ├── mailboxes.go
│   │   │   ├── search.go
│   │   │   ├── exports.go
│   │   │   └── jobs.go
│   │   └── routes.go            # Route registration
│   ├── cache/                   # Redis-backed caches and token stores
│   │   └── refreshstore.go
│   ├── config/                  # Configuration
│   │   └── config.go
│   ├── database/                # Database layer
│   │   ├── migrations/          # Migration files
│   │   ├── postgres.go          # Connection setup
│   │   └── repositories/        # Repository pattern
│   │       ├── user_repository.go
│   │       ├── tenant_repository.go
│   │       ├── mailbox_repository.go
│   │       ├── email_repository.go
│   │       └── job_repository.go
│   ├── services/                # Business logic
│   │   ├── auth_service.go
│   │   ├── sync_service.go
│   │   ├── search_service.go
│   │   ├── export_service.go
│   │   └── notification_service.go
│   ├── workers/                 # Background workers
│   │   ├── sync_worker.go
│   │   ├── export_worker.go
│   │   └── retention_worker.go
│   ├── scheduler/               # Cron scheduler
│   │   └── scheduler.go
│   ├── graph/                   # Microsoft Graph API client
│   │   └── client.go
│   ├── models/                  # Domain models
│   │   ├── user.go
│   │   ├── tenant.go
│   │   ├── mailbox.go
│   │   ├── email.go
│   │   └── job.go
│   └── utils/                   # Utilities
│       ├── crypto.go
│       ├── filesystem.go
│       └── logger.go
└── pkg/                         # Public packages
    └── errors/                  # Custom error types
        └── errors.go
```

**Controller/Handler Template:**

```go
// internal/api/handlers/tenants.go
package handlers

import (
    "github.com/gofiber/fiber/v3"
    "ironarchive/internal/services"
    "ironarchive/internal/models"
)

type TenantHandler struct {
    tenantService *services.TenantService
}

func NewTenantHandler(tenantService *services.TenantService) *TenantHandler {
    return &TenantHandler{tenantService: tenantService}
}

func (h *TenantHandler) List(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)

    tenants, err := h.tenantService.List(c.Context(), user)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.JSON(fiber.Map{
        "tenants": tenants,
    })
}

func (h *TenantHandler) Create(c *fiber.Ctx) error {
    var req models.CreateTenantRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }

    tenant, err := h.tenantService.Create(c.Context(), &req)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.Status(fiber.StatusCreated).JSON(tenant)
}
```

### Database Access Layer (Repository Pattern)

```go
// internal/database/repositories/tenant_repository.go
package repositories

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "ironarchive/internal/models"
)

type TenantRepository struct {
    db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) *TenantRepository {
    return &TenantRepository{db: db}
}

func (r *TenantRepository) FindAll(ctx context.Context, userRole string, userTenantID *string) ([]models.Tenant, error) {
    query := `
        SELECT id, name, azure_tenant_id, retention_policy_days, legal_hold, storage_bytes, created_at
        FROM tenants
    `

    var args []interface{}

    // Apply RBAC filtering
    if userRole == "TENANT_ADMIN" && userTenantID != nil {
        query += " WHERE id = $1"
        args = append(args, *userTenantID)
    }

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tenants []models.Tenant
    for rows.Next() {
        var tenant models.Tenant
        err := rows.Scan(
            &tenant.ID,
            &tenant.Name,
            &tenant.AzureTenantID,
            &tenant.RetentionPolicyDays,
            &tenant.LegalHold,
            &tenant.StorageBytes,
            &tenant.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        tenants = append(tenants, tenant)
    }

    return tenants, nil
}

func (r *TenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
    query := `
        INSERT INTO tenants (name, azure_tenant_id, azure_app_credentials)
        VALUES ($1, $2, pgp_sym_encrypt($3, $4))
        RETURNING id, created_at
    `

    err := r.db.QueryRow(
        ctx,
        query,
        tenant.Name,
        tenant.AzureTenantID,
        tenant.AzureAppCredentials, // JSON string with app_id + app_secret
        getEncryptionKey(), // From config
    ).Scan(&tenant.ID, &tenant.CreatedAt)

    return err
}
```

### Authentication and Authorization Middleware

```go
// internal/api/middleware/auth.go
package middleware

import (
    "github.com/gofiber/fiber/v3"
    "github.com/golang-jwt/jwt/v5"
    "ironarchive/internal/models"
)

func JWTMiddleware(jwtSecret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")
        if tokenString == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
        }

        // Remove "Bearer " prefix
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
        }

        claims := token.Claims.(jwt.MapClaims)
        user := &models.User{
            ID:    claims["user_id"].(string),
            Email: claims["email"].(string),
            Role:  claims["role"].(string),
        }

        c.Locals("user", user)
        return c.Next()
    }
}

func RequireMSPAdmin() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*models.User)
        if user.Role != "MSP_ADMIN" {
            return fiber.NewError(fiber.StatusForbidden, "MSP Admin role required")
        }
        return c.Next()
    }
}
```

### Configuration Loading

The configuration package exposes `config.Load()` which populates a typed `Config` struct from environment variables (loading `.env` in development). All services and handlers should receive configuration via dependency injection, avoiding global environment lookups.

### Refresh Token Session Store

Refresh token sessions are stored in Redis via `internal/cache/refreshstore.go`. The store persists hashed refresh token identifiers, enforces the seven-day TTL defined in the security requirements, and exposes helpers used by `AuthService` to revoke sessions on logout or refresh.
