## Unified Project Structure

```
ironarchive/
├── .github/
│   └── workflows/
│       ├── ci.yaml                  # Run tests, linting
│       ├── build.yaml               # Build multi-arch Docker images
│       └── deploy.yaml              # Deploy to test environment
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/                    # (see Backend Architecture above)
│   ├── pkg/
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── lib/                     # (see Frontend Architecture above)
│   │   └── routes/
│   ├── static/
│   ├── svelte.config.js
│   ├── tailwind.config.js
│   ├── vite.config.ts
│   ├── package.json
│   └── tsconfig.json
├── docker/
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   ├── docker-compose.yml
│   └── docker-compose.dev.yml      # Dev environment with hot-reload
├── migrations/
│   ├── 000001_initial_schema.up.sql
│   ├── 000001_initial_schema.down.sql
│   └── ...
├── docs/
│   ├── prd/                         # Sharded PRD files
│   ├── architecture/                # Sharded architecture files (this doc)
│   ├── CONTRIBUTING.md
│   └── DEPLOYMENT.md
├── scripts/
│   ├── setup-dev.sh                 # Local dev environment setup
│   ├── migrate.sh                   # Run database migrations
│   └── build-docker.sh              # Build Docker images
├── .env.example
├── .gitignore
├── LICENSE
├── Makefile                         # Top-level build commands
└── README.md
```

