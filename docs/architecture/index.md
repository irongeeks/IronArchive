# IronArchive Fullstack Architecture Document

## Table of Contents

- [IronArchive Fullstack Architecture Document](#ironarchive-fullstack-architecture-document)
  - [Introduction](./introduction.md)
    - [Starter Template or Existing Project](./introduction.md#starter-template-or-existing-project)
    - [Change Log](./introduction.md#change-log)
  - [High Level Architecture](./high-level-architecture.md)
    - [Technical Summary](./high-level-architecture.md#technical-summary)
    - [Platform and Infrastructure Choice](./high-level-architecture.md#platform-and-infrastructure-choice)
    - [Repository Structure](./high-level-architecture.md#repository-structure)
    - [High Level Architecture Diagram](./high-level-architecture.md#high-level-architecture-diagram)
    - [Architectural Patterns](./high-level-architecture.md#architectural-patterns)
  - [Tech Stack](./tech-stack.md)
    - [Technology Stack Table](./tech-stack.md#technology-stack-table)
  - [Data Models](./data-models.md)
    - [User](./data-models.md#user)
    - [Tenant](./data-models.md#tenant)
    - [Mailbox](./data-models.md#mailbox)
    - [Email](./data-models.md#email)
    - [Attachment](./data-models.md#attachment)
    - [Job](./data-models.md#job)
    - [AuditLog](./data-models.md#auditlog)
  - [API Specification](./api-specification.md)
    - [REST API Specification](./api-specification.md#rest-api-specification)
  - [Components](./components.md)
    - [API Server](./components.md#api-server)
    - [Sync Worker](./components.md#sync-worker)
    - [Job Queue Processor](./components.md#job-queue-processor)
    - [Scheduler](./components.md#scheduler)
    - [Microsoft Graph API Client](./components.md#microsoft-graph-api-client)
    - [Search Service](./components.md#search-service)
    - [Repository Layer](./components.md#repository-layer)
    - [Frontend Application](./components.md#frontend-application)
    - [Authentication Service](./components.md#authentication-service)
    - [Notification Service](./components.md#notification-service)
    - [Component Interaction Diagram](./components.md#component-interaction-diagram)
  - [External APIs](./external-apis.md)
    - [Microsoft Graph API](./external-apis.md#microsoft-graph-api)
    - [SMTP Server](./external-apis.md#smtp-server-email-notifications)
    - [Microsoft Teams Webhooks](./external-apis.md#microsoft-teams-webhooks)
    - [Discord Webhooks](./external-apis.md#discord-webhooks)
  - [Core Workflows](./core-workflows.md)
    - [Email Sync Workflow (Initial)](./core-workflows.md#email-sync-workflow-initial)
    - [Email Search Workflow](./core-workflows.md#email-search-workflow)
    - [Export Workflow (EML/ZIP)](./core-workflows.md#export-workflow-emlzip)
    - [Tenant Onboarding Workflow](./core-workflows.md#tenant-onboarding-workflow)
  - [Database Schema](./database-schema.md)
    - [Complete SQL Schema (PostgreSQL 16+)](./database-schema.md#complete-sql-schema-postgresql-16)
  - [Frontend Architecture](./frontend-architecture.md)
    - [Component Architecture](./frontend-architecture.md#component-architecture)
    - [State Management Architecture](./frontend-architecture.md#state-management-architecture)
    - [Routing Architecture](./frontend-architecture.md#routing-architecture)
    - [Frontend Services Layer](./frontend-architecture.md#frontend-services-layer)
  - [Backend Architecture](./backend-architecture.md)
    - [Service Architecture (Monolithic Go Binary)](./backend-architecture.md#service-architecture-monolithic-go-binary)
    - [Database Access Layer (Repository Pattern)](./backend-architecture.md#database-access-layer-repository-pattern)
    - [Authentication and Authorization Middleware](./backend-architecture.md#authentication-and-authorization-middleware)
  - [Unified Project Structure](./project-structure.md)
  - [Development Workflow](./development-workflow.md)
    - [Local Development Setup](./development-workflow.md#local-development-setup)
    - [Environment Configuration](./development-workflow.md#environment-configuration)
  - [Deployment Architecture](./deployment-architecture.md)
    - [Deployment Strategy](./deployment-architecture.md#deployment-strategy)
    - [CI/CD Pipeline](./deployment-architecture.md#cicd-pipeline)
    - [Environments](./deployment-architecture.md#environments)
  - [Security and Performance](./security-performance.md)
    - [Security Requirements](./security-performance.md#security-requirements)
    - [Performance Optimization](./security-performance.md#performance-optimization)
  - [Testing Strategy](./testing-strategy.md)
    - [Testing Pyramid](./testing-strategy.md#testing-pyramid)
    - [Test Organization](./testing-strategy.md#test-organization)
    - [Test Examples](./testing-strategy.md#test-examples)
  - [Coding Standards](./coding-standards.md)
    - [Critical Fullstack Rules](./coding-standards.md#critical-fullstack-rules)
    - [Naming Conventions](./coding-standards.md#naming-conventions)
  - [Error Handling Strategy](./error-handling.md)
    - [Error Flow](./error-handling.md#error-flow)
    - [Error Response Format](./error-handling.md#error-response-format)
    - [Frontend Error Handling](./error-handling.md#frontend-error-handling)
    - [Backend Error Handling](./error-handling.md#backend-error-handling)
  - [Monitoring and Observability](./monitoring.md)
    - [Monitoring Stack](./monitoring.md#monitoring-stack)
    - [Key Metrics](./monitoring.md#key-metrics)

---

## Overview

This document outlines the complete fullstack architecture for **IronArchive**, a modern, open-source, self-hosted M365 email archiving platform built specifically for MSPs (Managed Service Providers) serving multiple tenants.

The architecture prioritizes:
- **Cost-efficiency** (75% reduction vs commercial solutions)
- **Developer productivity** (AI-driven development workflows)
- **MSP-specific requirements** (whitelabeling, multi-tenant, compliance)

### Key Architectural Decisions

✅ **Monolithic Go backend** for operational simplicity with clear path to microservices
✅ **SvelteKit frontend** for best-in-class developer experience and performance
✅ **Docker Compose deployment** targeting self-hosted single-server environments
✅ **Hybrid data storage** optimizing for cost (filesystem), performance (Meilisearch), and integrity (PostgreSQL)
✅ **Multi-tenant SaaS patterns** with application-enforced isolation and RBAC
✅ **Compliance-first design** with immutable audit logs, retention policies, and legal holds

All development must follow this architecture specification to ensure consistency and enable successful AI-assisted development.
