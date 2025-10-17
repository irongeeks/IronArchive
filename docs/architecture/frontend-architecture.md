## Frontend Architecture

### Component Architecture

**Component Organization:**

```
/frontend/src/
├── lib/
│   ├── components/
│   │   ├── ui/                    # Shadcn-Svelte primitives
│   │   │   ├── button.svelte
│   │   │   ├── card.svelte
│   │   │   ├── input.svelte
│   │   │   ├── modal.svelte
│   │   │   └── table.svelte
│   │   ├── layout/                # Layout components
│   │   │   ├── Header.svelte
│   │   │   ├── Sidebar.svelte
│   │   │   ├── Nav.svelte
│   │   │   └── PageContainer.svelte
│   │   ├── features/              # Feature-specific components
│   │   │   ├── TenantWizard.svelte
│   │   │   ├── SearchBar.svelte
│   │   │   ├── EmailList.svelte
│   │   │   ├── EmailDetailPanel.svelte
│   │   │   └── StorageWidget.svelte
│   │   └── shared/                # Shared components
│   │       ├── LoadingSpinner.svelte
│   │       ├── ErrorBoundary.svelte
│   │       └── ThemeSelector.svelte
│   ├── stores/                    # Svelte stores
│   │   ├── auth.ts
│   │   ├── theme.ts
│   │   ├── tenants.ts
│   │   └── jobs.ts
│   ├── services/                  # API client services
│   │   ├── api-client.ts
│   │   ├── auth-service.ts
│   │   ├── tenant-service.ts
│   │   ├── search-service.ts
│   │   └── export-service.ts
│   ├── utils/                     # Utility functions
│   │   ├── formatters.ts
│   │   ├── validators.ts
│   │   └── date-helpers.ts
│   └── types/                     # TypeScript types
│       ├── api.ts
│       ├── models.ts
│       └── stores.ts
└── routes/                        # SvelteKit routes
    ├── +layout.svelte
    ├── +page.svelte
    ├── setup/
    │   └── +page.svelte
    ├── login/
    │   └── +page.svelte
    ├── dashboard/
    │   └── +page.svelte
    ├── tenants/
    │   ├── +page.svelte
    │   └── [id]/
    │       └── +page.svelte
    ├── search/
    │   └── +page.svelte
    ├── exports/
    │   └── +page.svelte
    ├── settings/
    │   └── +page.svelte
    └── profile/
        └── +page.svelte
```

**Component Template Example:**

```typescript
<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { Button, Card } from '$lib/components/ui';

  export let data; // From +page.ts load function

  let loading = false;
  let error: string | null = null;

  onMount(() => {
    // Component initialization
  });

  async function handleAction() {
    loading = true;
    error = null;

    try {
      // API call via service layer
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>

<!-- Template markup -->
<Card>
  <!-- Component content -->
</Card>

<style>
  /* Scoped component styles (complementing Tailwind) */
</style>
```

### State Management Architecture

**State Structure:**

```typescript
// lib/stores/auth.ts
import { writable } from 'svelte/store';
import type { User } from '$lib/types/models';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  loading: boolean;
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>({
    user: null,
    accessToken: null,
    isAuthenticated: false,
    loading: true
  });

  return {
    subscribe,
    login: (user: User, accessToken: string) => {
      update(state => ({
        ...state,
        user,
        accessToken,
        isAuthenticated: true,
        loading: false
      }));
      // Persist to localStorage
      localStorage.setItem('access_token', accessToken);
    },
    logout: () => {
      set({ user: null, accessToken: null, isAuthenticated: false, loading: false });
      localStorage.removeItem('access_token');
    },
    setLoading: (loading: boolean) => {
      update(state => ({ ...state, loading }));
    }
  };
}

export const authStore = createAuthStore();
```

**State Management Patterns:**
- Use Svelte stores for global state (auth, theme, active jobs)
- Use component-local reactive statements for UI-specific state
- Persist critical state (auth tokens, theme preference) to localStorage
- Use derived stores for computed values

### Routing Architecture

**Route Organization:**

```
/routes
├── +layout.svelte           # Root layout with nav/header
├── +layout.ts               # Root load function (setup check, auth)
├── +page.svelte             # Home/redirect to dashboard
├── setup/+page.svelte       # Setup wizard (no auth required)
├── login/+page.svelte       # Login page (no auth required)
└── (authenticated)/         # Route group requiring auth
    ├── +layout.svelte       # Authenticated layout
    ├── +layout.server.ts    # Server-side auth check
    ├── dashboard/+page.svelte
    ├── tenants/
    │   ├── +page.svelte
    │   ├── +page.ts         # Load tenants list
    │   └── [id]/
    │       ├── +page.svelte
    │       └── +page.ts     # Load tenant details
    ├── search/+page.svelte
    └── settings/+page.svelte
```

**Protected Route Pattern:**

```typescript
// routes/(authenticated)/+layout.server.ts
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, url }) => {
  const token = cookies.get('access_token');

  if (!token) {
    throw redirect(303, `/login?redirect=${url.pathname}`);
  }

  // Validate token and fetch user
  try {
    const user = await validateToken(token);
    return { user };
  } catch (err) {
    throw redirect(303, '/login');
  }
};
```

### Frontend Services Layer

**API Client Setup:**

```typescript
// lib/services/api-client.ts
import { authStore } from '$lib/stores/auth';
import { get } from 'svelte/store';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

interface ApiError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const auth = get(authStore);
    const headers = new Headers(options.headers);

    if (auth.accessToken) {
      headers.set('Authorization', `Bearer ${auth.accessToken}`);
    }

    headers.set('Content-Type', 'application/json');

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers
    });

    if (response.status === 401) {
      // Token expired, attempt refresh or logout
      authStore.logout();
      throw new Error('Session expired');
    }

    if (!response.ok) {
      const error: ApiError = await response.json();
      throw new Error(error.message);
    }

    return response.json();
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  async patch<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: JSON.stringify(data)
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const apiClient = new ApiClient();
```

**Service Example:**

```typescript
// lib/services/tenant-service.ts
import { apiClient } from './api-client';
import type { Tenant, Mailbox } from '$lib/types/models';

export const tenantService = {
  async list(): Promise<Tenant[]> {
    const response = await apiClient.get<{ tenants: Tenant[] }>('/tenants');
    return response.tenants;
  },

  async get(id: string): Promise<Tenant> {
    return apiClient.get<Tenant>(`/tenants/${id}`);
  },

  async create(data: Partial<Tenant>): Promise<Tenant> {
    return apiClient.post<Tenant>('/tenants', data);
  },

  async discoverMailboxes(tenantId: string): Promise<Mailbox[]> {
    return apiClient.post<Mailbox[]>(`/tenants/${tenantId}/discover-mailboxes`, {});
  },

  async triggerSync(tenantId: string): Promise<{ jobId: string }> {
    return apiClient.post(`/tenants/${tenantId}/sync`, {});
  }
};
```

