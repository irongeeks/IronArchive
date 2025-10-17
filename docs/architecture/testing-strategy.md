## Testing Strategy

### Testing Pyramid

```
         E2E Tests (10%)
    /                    \
   Integration Tests (30%)
 /                          \
Frontend Unit  Backend Unit (60%)
```

### Test Organization

**Frontend Tests:**

```
/frontend/src/
├── lib/
│   ├── components/
│   │   └── SearchBar.test.ts       # Component tests
│   ├── services/
│   │   └── api-client.test.ts      # Service tests
│   └── utils/
│       └── formatters.test.ts      # Utility tests
└── routes/
    └── dashboard/
        └── +page.test.ts             # Page tests
```

**Backend Tests:**

```
/backend/
├── internal/
│   ├── api/handlers/
│   │   └── tenants_test.go          # Handler tests
│   ├── services/
│   │   └── sync_service_test.go     # Service tests
│   └── database/repositories/
│       └── tenant_repository_test.go # Repository tests
```

**E2E Tests:**

```
/tests/e2e/
├── setup-wizard.spec.ts
├── login.spec.ts
├── tenant-onboarding.spec.ts
└── email-search-export.spec.ts
```

### Test Examples

**Frontend Component Test:**

```typescript
// lib/components/SearchBar.test.ts
import { render, screen, fireEvent } from '@testing-library/svelte';
import { vi } from 'vitest';
import SearchBar from './SearchBar.svelte';

describe('SearchBar', () => {
  it('debounces search input', async () => {
    const onSearch = vi.fn();
    render(SearchBar, { props: { onSearch } });

    const input = screen.getByPlaceholderText('Search emails...');
    await fireEvent.input(input, { target: { value: 'test' } });

    // Should not call immediately
    expect(onSearch).not.toHaveBeenCalled();

    // Should call after 300ms
    await new Promise(resolve => setTimeout(resolve, 350));
    expect(onSearch).toHaveBeenCalledWith('test');
  });
});
```

**Backend API Test:**

```go
// internal/api/handlers/tenants_test.go
package handlers_test

import (
    "testing"
    "net/http/httptest"
    "github.com/gofiber/fiber/v3"
    "github.com/stretchr/testify/assert"
)

func TestTenantHandler_List(t *testing.T) {
    app := fiber.New()
    handler := NewTenantHandler(mockTenantService)
    app.Get("/tenants", handler.List)

    req := httptest.NewRequest("GET", "/tenants", nil)
    req.Header.Set("Authorization", "Bearer valid-token")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

**E2E Test:**

```typescript
// tests/e2e/tenant-onboarding.spec.ts
import { test, expect } from '@playwright/test';

test('MSP Admin can onboard new tenant', async ({ page }) => {
  await page.goto('/login');
  await page.fill('[name="email"]', 'admin@example.com');
  await page.fill('[name="password"]', 'password123');
  await page.click('button[type="submit"]');

  await expect(page).toHaveURL('/dashboard');

  await page.click('text=Add Tenant');
  await page.fill('[name="name"]', 'Acme Corp');
  await page.fill('[name="azureTenantId"]', 'tenant-uuid-here');
  await page.click('button:has-text("Next")');

  // Continue wizard steps...

  await expect(page.locator('text=Tenant created successfully')).toBeVisible();
});
```

