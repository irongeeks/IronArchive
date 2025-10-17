# Introduction

This document outlines the complete fullstack architecture for **IronArchive**, including backend systems, frontend implementation, and their integration. It serves as the single source of truth for AI-driven development, ensuring consistency across the entire technology stack.

This unified approach combines what would traditionally be separate backend and frontend architecture documents, streamlining the development process for modern fullstack applications where these concerns are increasingly intertwined.

IronArchive is a modern, open-source, self-hosted M365 email archiving platform built specifically for MSPs (Managed Service Providers) serving multiple tenants. The architecture prioritizes cost-efficiency, compliance excellence (DSGVO/GoBD, GDPR), instant search capabilities, and comprehensive multi-tenant management with whitelabeling support.

## Starter Template or Existing Project

**Decision: N/A - Greenfield Project**

IronArchive is a greenfield project built from scratch without relying on pre-existing starter templates. This approach provides maximum flexibility to implement the specific technical stack defined in the PRD:

- **Backend:** Go 1.24 with Fiber v3 framework
- **Frontend:** SvelteKit 2.x with TailwindCSS 4.x
- **Architecture:** Monorepo structure with Docker Compose deployment

While there are fullstack monorepo starters available (T3 Stack, Turborepo templates), none align precisely with the Go + Svelte combination and MSP-specific multi-tenant requirements. Building from scratch ensures optimal architecture for the unique compliance, performance, and whitelabeling requirements.

## Change Log

| Date | Version | Description | Author |
|------|---------|-------------|---------|
| 2025-10-17 | 1.0 | Initial architecture document creation | Winston (Architect Agent) |
