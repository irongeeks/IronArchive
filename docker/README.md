# Docker Configuration

## Resource Limits

Resource limits prevent runaway processes from consuming all system resources.

### Current Limits (Development)

| Service | CPU Limit | Memory Limit | CPU Reserved | Memory Reserved |
|---------|-----------|--------------|--------------|-----------------|
| PostgreSQL | 1.0 cores | 1GB | 0.25 cores | 256MB |
| Redis | 0.5 cores | 512MB | 0.1 cores | 128MB |
| Meilisearch | 1.0 cores | 1GB | 0.25 cores | 256MB |

### Production Tuning

For production deployments, adjust limits based on workload:

**PostgreSQL**:
- Light: 2 cores, 2GB
- Medium: 4 cores, 4GB
- Heavy: 8 cores, 8GB

**Redis**:
- Light: 0.5 cores, 512MB
- Medium: 1 core, 1GB
- Heavy: 2 cores, 2GB

**Meilisearch**:
- Light: 1 core, 1GB
- Medium: 2 cores, 2GB
- Heavy: 4 cores, 4GB

## Usage

Start services:
```bash
cd docker
docker-compose up -d
```

Monitor resource usage:
```bash
docker stats
```

Stop services:
```bash
docker-compose down
```
