# Next Steps

## UX Expert Prompt

I need your expertise to design the user experience and interface architecture for IronArchive.

Please review the attached PRD (`docs/prd.md`) and Project Brief (`docs/brief.md`), then create detailed UX/UI specifications including:

1. **Information Architecture:** Complete sitemap with page hierarchy and navigation flows
2. **Wireframes:** Low-fidelity wireframes for all core screens identified in the "User Interface Design Goals" section
3. **Component Library:** Design system specification (typography, spacing, color tokens, reusable components)
4. **Interaction Patterns:** Detailed interaction specs for complex workflows (tenant onboarding wizard, search interface, task monitoring)
5. **Responsive Breakpoints:** Specific layouts for mobile (320px), tablet (768px), desktop (1920px+)
6. **Accessibility Guidelines:** Implementation guidance for WCAG 2.1 AA compliance
7. **Theme System Design:** CSS variable structure for the 6 themes (Catppuccin Mocha/Latte, Nord, Cyberpunk, Dracula, Tokyo Night)

**Key Priorities:**
- Progressive disclosure: empty states → gradually revealing complexity
- Real-time feedback: live updates without page refresh
- Instant visual confirmation for all user actions

**Deliverables:**
- Wireframes (Figma, Sketch, or similar)
- Design system documentation
- Component specifications with props/variants
- Accessibility annotation on wireframes

**Timeline:** Please complete initial wireframes and design system within 2 weeks to unblock frontend development (Epic 1, Story 1.6-1.9).

## Architect Prompt

I need your expertise to design the technical architecture for IronArchive.

Please review the attached PRD (`docs/prd.md`) and Project Brief (`docs/brief.md`), then create a comprehensive architecture document including:

1. **System Architecture Diagram:** High-level component diagram showing backend, frontend, databases, external services (Microsoft Graph API), and data flows
2. **Database Schema:** Detailed ERD with all tables, relationships, indexes, constraints based on Story 1.2 plus any optimizations you recommend
3. **API Architecture:** RESTful API design with endpoint specifications, request/response schemas, authentication flows
4. **Service Layer Design:** Go backend package structure (`/cmd`, `/internal`, `/pkg`) with clear separation of concerns (handlers, services, repositories)
5. **Sync Engine Architecture:** Detailed design for email sync workflow including delta query handling, job queue processing, filesystem storage, Meilisearch indexing
6. **Frontend Architecture:** SvelteKit project structure, routing strategy, state management approach, API client design
7. **Security Architecture:** Multi-tenant data isolation strategy (application-enforced vs. PostgreSQL RLS), authentication/authorization flow, CSRF protection
8. **Deployment Architecture:** Docker Compose service definitions, networking, volume mounts, environment configuration
9. **Scalability Considerations:** Horizontal scaling strategy (if needed post-v1.0), database partitioning recommendations, caching strategy
10. **Performance Optimization:** Database query optimization, Meilisearch index configuration, async job processing patterns

**Critical Architecture Decisions Needed:**

1. **Multi-Tenant Isolation:** Application-enforced tenant filtering on all queries vs. PostgreSQL Row-Level Security (RLS) - recommend approach with pros/cons
2. **Attachment Storage:** Filesystem structure optimization for millions of files - validate `/archive/tenants/{uuid}/attachments/{sha256}.{ext}` approach
3. **Meilisearch Index Design:** Searchable attributes, ranking rules, typo tolerance config for optimal <200ms query performance
4. **Job Queue Concurrency:** Redis lock pattern for preventing simultaneous mailbox syncs, concurrency limits (max 5 concurrent jobs)
5. **Database Indexing:** B-tree vs. GIN indexes for specific queries, partitioning strategy for emails table if >10M rows

**Technical Stack (from PRD):**
- **Backend:** Go 1.24, Fiber v3, pgx/GORM, asynq (Redis job queue), msgraph-sdk-go
- **Frontend:** SvelteKit 2.x, TailwindCSS 4.x, Shadcn-Svelte, Chart.js
- **Infrastructure:** PostgreSQL 16, Redis 7, Meilisearch 1.6, Docker Compose, Traefik

**Deliverables:**
- Architecture document (`docs/architecture.md`)
- Database schema SQL with migrations
- API specification (OpenAPI/Swagger format preferred)
- Package structure with placeholder files/comments
- Deployment configuration (docker-compose.yml, .env template)

**Timeline:** Please complete architecture document within 2 weeks to enable Epic 1 implementation (Stories 1.1-1.3 depend on architecture decisions).

---

**🎉 PRD COMPLETE - Ready for Implementation Planning!**

