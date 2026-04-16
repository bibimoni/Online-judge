# Online-judge

## Quick Start

### Option 1: Each service in its own container (docker-compose.yml)

```bash
docker compose up --build
```

### Option 2: All services in one dev container (dev-compose.yml)

Build once, then run the script inside:

```bash
# Build the dev image (only needed once, or when Dockerfile.dev changes)
./scripts/build-docker.sh

# Start infrastructure (databases, redis) + dev container
./scripts/run-docker.sh

# Attach a shell
docker exec -it online-judge-dev bash

# Start all services
./scripts/start.sh

# Or invidually
./scripts/start-<service_name>.sh

# View logs
./scripts/log.sh 
```

Stop everything:

```bash
docker compose -f dev-compose.yml down
```

## Services

| Service | Port | Language |
|---|---|---|
| gateway | 81 | Go |
| submission-judge | 8000 | Go |
| contest | 8001 | Go |
| problem | 3000 | Go |
| auth-v2 | 50051 | Node.js |
