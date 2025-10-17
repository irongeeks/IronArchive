## Deployment Architecture

### Deployment Strategy

**Frontend Deployment:**
- **Platform:** Served by backend (static files compiled into Go binary) OR separate CDN/static host
- **Build Command:** `npm run build` (generates static files in `/build`)
- **Output Directory:** `frontend/build`
- **CDN/Edge:** Optional Cloudflare CDN for static assets

**Backend Deployment:**
- **Platform:** Docker container on Linux server (Ubuntu 22.04+)
- **Build Command:** `go build -o ironarchive cmd/server/main.go`
- **Deployment Method:** Docker Compose OR single Docker container with volume mounts

### CI/CD Pipeline

```yaml
# .github/workflows/ci.yaml
name: CI

on: [push, pull_request]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - name: Run backend tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...

  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - name: Run frontend tests
        run: |
          cd frontend
          npm ci
          npm run test

  build-docker:
    needs: [test-backend, test-frontend]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - uses: docker/setup-buildx-action@v2
      - uses: docker/build-push-action@v4
        with:
          context: .
          file: docker/Dockerfile.backend
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ironarchive/backend:latest
```

### Environments

| Environment | Frontend URL | Backend URL | Purpose |
|-------------|-------------|-------------|---------|
| Development | http://localhost:5173 | http://localhost:3000 | Local development |
| Staging | https://staging.ironarchive.example.com | https://api-staging.ironarchive.example.com | Pre-production testing |
| Production | https://ironarchive.example.com | https://api.ironarchive.example.com | Live environment |

