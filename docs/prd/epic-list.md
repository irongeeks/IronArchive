# Epic List

The IronArchive v1.0 development is structured into **5 sequential epics**, each delivering significant, end-to-end functionality that builds upon previous work. Each epic is fully deployable and testable, following agile best practices.

**Epic 1: Foundation & Core Infrastructure**
Establish project scaffolding, database schema, authentication system, and basic API infrastructure. Delivers a functional setup wizard that creates the first MSP Admin account and initializes the system. This epic sets up the development environment and core services (PostgreSQL, Redis, Meilisearch) while delivering a deployable "hello world" level application with authentication.

**Epic 2: Tenant & Mailbox Management**
Build the tenant onboarding wizard, Microsoft Graph API integration, mailbox discovery, and sync scheduling engine. Delivers the ability to add M365 tenants, discover mailboxes, and initiate manual backups. By the end of this epic, the system can successfully connect to M365, retrieve emails, and store them in the archive.

**Epic 3: Search, Retrieval & Export**
Implement the Meilisearch-powered search interface, email detail views, and export functionality (EML/ZIP and PST formats). Delivers the core user value proposition: fast, intuitive search across archived emails with the ability to export results. This epic makes the archive practically useful to end users and tenant admins.

**Epic 4: Dashboard, Monitoring & Observability**
Build the adaptive dashboard with storage usage widgets, task monitoring interface, and notification system (Email, Teams, Discord). Delivers visibility into system operations, backup job status, and proactive error alerting. This epic transforms IronArchive from a functional tool into a production-ready, monitorable system.

**Epic 5: Compliance, Settings & Polish**
Implement retention policies, legal hold, audit logging, whitelabeling, theme system, and profile settings. Delivers the compliance features required for DSGVO/GoBD adherence and the customization capabilities that differentiate IronArchive from competitors. This epic completes the v1.0 feature set and prepares for public release.

---

**Rationale for Epic Structure:**

- **Epic 1 is foundational yet functional**: Delivers working authentication and setup wizard (not just scaffolding)
- **Sequential value delivery**: Each epic builds on previous work without dependencies on future epics
- **Testable milestones**: Each epic conclusion represents a demonstrable feature set
- **Balanced scope**: No epic is disproportionately large or small; each represents ~2-4 weeks of AI-assisted development
- **Cross-cutting concerns integrated**: Testing, logging, error handling flow through stories rather than being "tacked on" at the end
