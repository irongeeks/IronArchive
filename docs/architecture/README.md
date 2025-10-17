# IronArchive Architecture Documentation

This directory contains the sharded architecture documentation for IronArchive v1.0.

## 📚 Quick Navigation

Start with **[index.md](./index.md)** for the complete table of contents.

## 📁 Document Structure

| File | Description | Lines |
|------|-------------|-------|
| [index.md](./index.md) | Main table of contents and overview | 107 |
| [introduction.md](./introduction.md) | Project intro, greenfield decision, changelog | 25 |
| [high-level-architecture.md](./high-level-architecture.md) | Technical summary, platform choice, diagrams, patterns | 149 |
| [tech-stack.md](./tech-stack.md) | Complete technology stack with versions | 31 |
| [data-models.md](./data-models.md) | 7 core domain models with TypeScript interfaces | 266 |
| [api-specification.md](./api-specification.md) | Complete REST API (OpenAPI 3.0 spec) | 512 |
| [components.md](./components.md) | 10+ logical components with responsibilities | 211 |
| [external-apis.md](./external-apis.md) | Microsoft Graph, SMTP, Teams, Discord integrations | 71 |
| [core-workflows.md](./core-workflows.md) | 4 sequence diagrams for key workflows | 181 |
| [database-schema.md](./database-schema.md) | Complete PostgreSQL schema with migrations | 175 |
| [frontend-architecture.md](./frontend-architecture.md) | SvelteKit structure, routing, state management | 326 |
| [backend-architecture.md](./backend-architecture.md) | Go backend structure, repositories, middleware | 251 |
| [project-structure.md](./project-structure.md) | Complete monorepo directory structure | 53 |
| [development-workflow.md](./development-workflow.md) | Local setup, dev commands, env config | 110 |
| [deployment-architecture.md](./deployment-architecture.md) | Docker Compose deployment, CI/CD | 73 |
| [security-performance.md](./security-performance.md) | Security requirements and optimization targets | 31 |
| [testing-strategy.md](./testing-strategy.md) | Testing pyramid with frontend/backend/E2E examples | 133 |
| [coding-standards.md](./coding-standards.md) | Critical fullstack rules and naming conventions | 31 |
| [error-handling.md](./error-handling.md) | Error flow, formats, frontend/backend handlers | 126 |
| [monitoring.md](./monitoring.md) | Monitoring stack and key metrics | 41 |

**Total:** 2,903 lines of comprehensive architecture documentation

## 🎯 Key Architectural Decisions

✅ **Monolithic Go backend** (Fiber v3) with clear microservices separation path  
✅ **SvelteKit 2.x frontend** with SSR+SPA hybrid mode  
✅ **Docker Compose deployment** for single-server self-hosted environments  
✅ **Hybrid data storage**: PostgreSQL + Filesystem + Meilisearch  
✅ **Multi-tenant SaaS** with application-enforced isolation and RBAC  
✅ **Compliance-first** with immutable audit logs and retention policies  

## 🚀 For Developers

All AI agents and developers MUST follow this architecture specification for consistency in AI-assisted development.

**Core references for development:**
1. [tech-stack.md](./tech-stack.md) - Exact versions and technologies
2. [data-models.md](./data-models.md) - Domain models and TypeScript interfaces
3. [api-specification.md](./api-specification.md) - REST API contracts
4. [coding-standards.md](./coding-standards.md) - Critical fullstack rules
5. [database-schema.md](./database-schema.md) - Complete PostgreSQL DDL

---

**Version:** 1.0  
**Last Updated:** 2025-10-17  
**Author:** Winston (Architect Agent)
