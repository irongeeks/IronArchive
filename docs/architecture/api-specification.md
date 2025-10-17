## API Specification

IronArchive uses a REST API with OpenAPI 3.0 specification. The API is versioned under `/api/v1/` prefix.

### REST API Specification

```yaml
openapi: 3.0.0
info:
  title: IronArchive API
  version: 1.0.0
  description: REST API for IronArchive M365 email archiving platform
servers:
  - url: https://ironarchive.example.com/api/v1
    description: Production server
  - url: http://localhost:3000/api/v1
    description: Local development server

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    Error:
      type: object
      properties:
        error:
          type: object
          properties:
            code:
              type: string
            message:
              type: string
            details:
              type: object
            timestamp:
              type: string
              format: date-time
            requestId:
              type: string

security:
  - bearerAuth: []

paths:
  # Authentication Endpoints
  /auth/login:
    post:
      summary: User login
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                  format: email
                password:
                  type: string
              required: [email, password]
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  accessToken:
                    type: string
                  refreshToken:
                    type: string
                  user:
                    $ref: '#/components/schemas/User'

  /auth/refresh:
    post:
      summary: Refresh access token
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                refreshToken:
                  type: string
      responses:
        '200':
          description: Token refreshed
          content:
            application/json:
              schema:
                type: object
                properties:
                  accessToken:
                    type: string

  # Setup Endpoints
  /setup/status:
    get:
      summary: Check setup status
      security: []
      responses:
        '200':
          description: Setup status
          content:
            application/json:
              schema:
                type: object
                properties:
                  setupComplete:
                    type: boolean

  /setup/initialize:
    post:
      summary: Initialize first MSP Admin
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                  format: email
                password:
                  type: string
                displayName:
                  type: string
      responses:
        '201':
          description: Setup completed
          content:
            application/json:
              schema:
                type: object
                properties:
                  accessToken:
                    type: string
                  refreshToken:
                    type: string
                  user:
                    $ref: '#/components/schemas/User'

  # Tenant Endpoints
  /tenants:
    get:
      summary: List tenants
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
      responses:
        '200':
          description: List of tenants
          content:
            application/json:
              schema:
                type: object
                properties:
                  tenants:
                    type: array
                    items:
                      $ref: '#/components/schemas/Tenant'
                  total:
                    type: integer
                  page:
                    type: integer
                  limit:
                    type: integer

    post:
      summary: Create new tenant
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                azureTenantId:
                  type: string
                azureAppId:
                  type: string
                azureAppSecret:
                  type: string
      responses:
        '201':
          description: Tenant created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Tenant'

  /tenants/{id}:
    get:
      summary: Get tenant details
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Tenant details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Tenant'

    patch:
      summary: Update tenant
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                retentionPolicyDays:
                  type: integer
                legalHold:
                  type: boolean
      responses:
        '200':
          description: Tenant updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Tenant'

  /tenants/{id}/sync:
    post:
      summary: Trigger tenant-wide sync
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '202':
          description: Sync job enqueued
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Job'

  /tenants/{id}/mailboxes:
    get:
      summary: List mailboxes for tenant
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: List of mailboxes
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Mailbox'

  /tenants/{id}/discover-mailboxes:
    post:
      summary: Discover mailboxes from M365 tenant
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Mailboxes discovered
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Mailbox'

  # Mailbox Endpoints
  /mailboxes/{id}:
    get:
      summary: Get mailbox details
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Mailbox details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Mailbox'

    patch:
      summary: Update mailbox settings
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                syncEnabled:
                  type: boolean
      responses:
        '200':
          description: Mailbox updated

  /mailboxes/{id}/sync:
    post:
      summary: Trigger mailbox sync
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '202':
          description: Sync job enqueued
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Job'

  # Search Endpoints
  /search:
    get:
      summary: Search emails
      parameters:
        - name: q
          in: query
          required: true
          schema:
            type: string
          description: Search query
        - name: tenantId
          in: query
          schema:
            type: string
            format: uuid
        - name: mailboxId
          in: query
          schema:
            type: string
            format: uuid
        - name: from
          in: query
          schema:
            type: string
            format: date-time
        - name: to
          in: query
          schema:
            type: string
            format: date-time
        - name: hasAttachments
          in: query
          schema:
            type: boolean
        - name: limit
          in: query
          schema:
            type: integer
            default: 50
        - name: offset
          in: query
          schema:
            type: integer
            default: 0
      responses:
        '200':
          description: Search results
          content:
            application/json:
              schema:
                type: object
                properties:
                  hits:
                    type: array
                    items:
                      $ref: '#/components/schemas/Email'
                  total:
                    type: integer
                  processingTimeMs:
                    type: integer

  # Export Endpoints
  /exports:
    post:
      summary: Create export job
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                emailIds:
                  type: array
                  items:
                    type: string
                    format: uuid
                format:
                  type: string
                  enum: [EML_ZIP, PST]
      responses:
        '202':
          description: Export job created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Job'

  # Job Endpoints
  /jobs/{id}:
    get:
      summary: Get job status
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Job details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Job'

  /jobs:
    get:
      summary: List jobs
      parameters:
        - name: status
          in: query
          schema:
            type: string
            enum: [QUEUED, RUNNING, COMPLETED, FAILED]
        - name: type
          in: query
          schema:
            type: string
        - name: limit
          in: query
          schema:
            type: integer
            default: 50
      responses:
        '200':
          description: List of jobs
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Job'
```
