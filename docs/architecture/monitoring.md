## Monitoring and Observability

### Monitoring Stack

- **Frontend Monitoring:** Browser console errors logged to backend, optional Sentry integration
- **Backend Monitoring:** Prometheus metrics exported on `/metrics` endpoint (optional)
- **Error Tracking:** Structured JSON logging to stdout, aggregated by log management tool
- **Performance Monitoring:** API response time metrics, database query performance tracking

### Key Metrics

**Frontend Metrics:**
- Core Web Vitals (LCP, FID, CLS)
- JavaScript errors count and stack traces
- API response times (client-side measurement)
- User interactions (search queries, exports)

**Backend Metrics:**
- Request rate (requests/second)
- Error rate (errors/total requests)
- Response time (p50, p95, p99)
- Database query performance (slow queries > 1s)
- Email sync throughput (emails/hour)
- Job queue depth (pending jobs count)

---

## Summary

This architecture document defines IronArchive's complete fullstack architecture for v1.0 release. Key decisions:

✅ **Monolithic Go backend** for operational simplicity with clear path to microservices
✅ **SvelteKit frontend** for best-in-class developer experience and performance
✅ **Docker Compose deployment** targeting self-hosted single-server environments
✅ **Hybrid data storage** optimizing for cost (filesystem), performance (Meilisearch), and integrity (PostgreSQL)
✅ **Multi-tenant SaaS patterns** with application-enforced isolation and RBAC
✅ **Compliance-first design** with immutable audit logs, retention policies, and legal holds

The architecture prioritizes **cost-efficiency** (75% reduction vs commercial solutions), **developer productivity** (AI-driven development workflows), and **MSP-specific requirements** (whitelabeling, multi-tenant, compliance).

All development must follow this architecture specification to ensure consistency and enable successful AI-assisted development.
