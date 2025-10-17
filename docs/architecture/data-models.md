## Data Models

IronArchive's core domain entities are designed to support multi-tenant email archiving with compliance features and whitelabeling.

### User

**Purpose:** Represents system users with role-based access (MSP Admin, Tenant Admin, User)

**Key Attributes:**
- `id`: UUID - Primary key
- `email`: string - Unique email address (login credential)
- `password_hash`: string - bcrypt hash (cost factor 12+)
- `display_name`: string - User's display name
- `role`: enum - MSP Admin | Tenant Admin | User
- `tenant_id`: UUID (nullable) - Foreign key to tenant (null for MSP Admins)
- `mfa_secret`: string (nullable) - TOTP secret for MFA
- `mfa_enabled`: boolean - MFA enrollment status
- `created_at`: timestamp
- `updated_at`: timestamp

**TypeScript Interface:**

```typescript
interface User {
  id: string;
  email: string;
  displayName: string;
  role: 'MSP_ADMIN' | 'TENANT_ADMIN' | 'USER';
  tenantId?: string;
  mfaEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}
```

**Relationships:**
- Belongs to Tenant (if role is Tenant Admin or User)
- Has many AuditLogs

### Tenant

**Purpose:** Represents an MSP customer with M365 tenant configuration

**Key Attributes:**
- `id`: UUID - Primary key
- `name`: string - Customer/tenant display name
- `azure_tenant_id`: UUID - Microsoft 365 tenant ID
- `azure_app_credentials`: string - Encrypted JSON blob with app_id + app_secret via pgcrypto
- `retention_policy_days`: integer - Email retention period (default from global settings)
- `legal_hold`: boolean - Legal hold flag (prevents deletion)
- `whitelabel_config`: JSONB - Custom branding (logo URL, colors, etc.)
- `created_at`: timestamp
- `storage_bytes`: bigint - Total storage used (computed)

**TypeScript Interface:**

```typescript
interface Tenant {
  id: string;
  name: string;
  azureTenantId: string;
  azureAppId: string;
  retentionPolicyDays: number;
  legalHold: boolean;
  whitelabelConfig?: WhitelabelConfig;
  createdAt: string;
  storageBytes: number;
}

interface WhitelabelConfig {
  logoUrl?: string;
  faviconUrl?: string;
  primaryColor?: string;
  secondaryColor?: string;
  companyName?: string;
}
```

**Relationships:**
- Has many Mailboxes
- Has many Users (Tenant Admins and Users)
- Has many Jobs

### Mailbox

**Purpose:** Represents an M365 mailbox configured for backup

**Key Attributes:**
- `id`: UUID - Primary key
- `tenant_id`: UUID - Foreign key to tenant
- `email_address`: string - Mailbox email (unique within tenant)
- `display_name`: string - Mailbox display name
- `mailbox_type`: enum - User | Shared | Room | Equipment
- `sync_enabled`: boolean - Whether mailbox is actively synced
- `last_sync_at`: timestamp (nullable) - Last successful sync time
- `last_delta_token`: string (nullable) - Graph API delta token for incremental sync
- `email_count`: integer - Total emails archived (computed)
- `storage_bytes`: bigint - Mailbox storage usage (computed)
- `created_at`: timestamp

**TypeScript Interface:**

```typescript
interface Mailbox {
  id: string;
  tenantId: string;
  emailAddress: string;
  displayName: string;
  mailboxType: 'USER' | 'SHARED' | 'ROOM' | 'EQUIPMENT';
  syncEnabled: boolean;
  lastSyncAt?: string;
  emailCount: number;
  storageBytes: number;
  createdAt: string;
}
```

**Relationships:**
- Belongs to Tenant
- Has many Emails
- Has many Jobs

### Email

**Purpose:** Represents an archived email message with metadata

**Key Attributes:**
- `id`: UUID - Primary key
- `mailbox_id`: UUID - Foreign key to mailbox
- `message_id`: string - Microsoft Graph message ID (unique)
- `subject`: string - Email subject
- `sender`: string - Sender email address
- `recipients`: string[] - Array of recipient emails (To, CC, BCC combined)
- `sent_at`: timestamp - Email send time
- `has_attachments`: boolean - Attachment presence flag
- `size_bytes`: integer - Total email size including attachments
- `file_path`: string - Filesystem path to email body JSON
- `indexed_at`: timestamp (nullable) - Meilisearch indexing time
- `deleted_at`: timestamp (nullable) - Soft delete timestamp
- `created_at`: timestamp

**TypeScript Interface:**

```typescript
interface Email {
  id: string;
  mailboxId: string;
  messageId: string;
  subject: string;
  sender: string;
  recipients: string[];
  sentAt: string;
  hasAttachments: boolean;
  sizeBytes: number;
  filePath: string;
  indexedAt?: string;
  deletedAt?: string;
  createdAt: string;
}
```

**Relationships:**
- Belongs to Mailbox
- Has many Attachments

### Attachment

**Purpose:** Represents email attachment with deduplication

**Key Attributes:**
- `id`: UUID - Primary key
- `email_id`: UUID - Foreign key to email
- `filename`: string - Original attachment filename
- `content_type`: string - MIME type
- `size_bytes`: integer - Attachment size
- `sha256_hash`: string - SHA-256 hash for deduplication
- `file_path`: string - Filesystem path (deduplicated by hash)
- `created_at`: timestamp

**TypeScript Interface:**

```typescript
interface Attachment {
  id: string;
  emailId: string;
  filename: string;
  contentType: string;
  sizeBytes: number;
  sha256Hash: string;
  filePath: string;
  createdAt: string;
}
```

**Relationships:**
- Belongs to Email

### Job

**Purpose:** Tracks background job execution (sync, export, retention cleanup)

**Key Attributes:**
- `id`: UUID - Primary key
- `type`: enum - SYNC_MAILBOX | SYNC_TENANT | SYNC_ALL | EXPORT | RETENTION_CLEANUP
- `status`: enum - QUEUED | RUNNING | COMPLETED | FAILED
- `tenant_id`: UUID (nullable) - Associated tenant
- `mailbox_id`: UUID (nullable) - Associated mailbox (for SYNC_MAILBOX)
- `user_id`: UUID (nullable) - User who initiated job (for exports)
- `progress`: integer - Progress percentage (0-100)
- `metadata`: JSONB - Job-specific data (export format, error details, etc.)
- `started_at`: timestamp (nullable)
- `completed_at`: timestamp (nullable)
- `created_at`: timestamp

**TypeScript Interface:**

```typescript
interface Job {
  id: string;
  type: 'SYNC_MAILBOX' | 'SYNC_TENANT' | 'SYNC_ALL' | 'EXPORT' | 'RETENTION_CLEANUP';
  status: 'QUEUED' | 'RUNNING' | 'COMPLETED' | 'FAILED';
  tenantId?: string;
  mailboxId?: string;
  userId?: string;
  progress: number;
  metadata?: Record<string, any>;
  startedAt?: string;
  completedAt?: string;
  createdAt: string;
}
```

**Relationships:**
- Belongs to Tenant (optional)
- Belongs to Mailbox (optional)
- Belongs to User (optional)

### AuditLog

**Purpose:** Immutable audit trail for compliance (DSGVO/GoBD)

**Key Attributes:**
- `id`: UUID - Primary key
- `user_id`: UUID (nullable) - User who performed action
- `action`: string - Action performed (LOGIN, SEARCH, EXPORT, CONFIG_CHANGE, etc.)
- `ip_address`: string - Client IP address
- `details`: JSONB - Action-specific details (search query, export params, etc.)
- `timestamp`: timestamp - Action timestamp (immutable)

**TypeScript Interface:**

```typescript
interface AuditLog {
  id: string;
  userId?: string;
  action: string;
  ipAddress: string;
  details: Record<string, any>;
  timestamp: string;
}
```

**Relationships:**
- Belongs to User (optional)

