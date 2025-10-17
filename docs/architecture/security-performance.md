## Security and Performance

### Security Requirements

**Frontend Security:**
- **CSP Headers:** `Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';`
- **XSS Prevention:** HTML sanitization on email body rendering, escape user input
- **Secure Storage:** JWT tokens in HTTP-only cookies, no sensitive data in localStorage

**Backend Security:**
- **Input Validation:** Strict validation using Go validator library, parameterized SQL queries
- **Rate Limiting:** 100 requests/minute per IP for API endpoints, 10 requests/minute for auth endpoints
- **CORS Policy:** Whitelist frontend origin only

**Authentication Security:**
- **Token Storage:** Access tokens in HTTP-only cookies (SameSite=Strict), refresh tokens in secure cookies
- **Session Management:** 15-minute access token expiry, 7-day refresh token expiry
- **Password Policy:** Minimum 12 characters, uppercase, lowercase, number, special character

### Performance Optimization

**Frontend Performance:**
- **Bundle Size Target:** < 500KB gzipped JavaScript
- **Loading Strategy:** Code splitting per route, lazy load heavy components
- **Caching Strategy:** Service worker caching for static assets, ETags for API responses

**Backend Performance:**
- **Response Time Target:** 95th percentile < 500ms for API endpoints
- **Database Optimization:** Connection pooling (max 25 connections), prepared statements, indexed queries
- **Caching Strategy:** Redis cache for frequently accessed data (tenants, mailboxes), 5-minute TTL

