# Online-judge

## Quick Start

### Option 1: Each service in its own container (docker-compose.yml)

```bash
docker compose up --build
```

### Option 2: All services in one dev container (dev-compose.yml)

```bash
# Build and start the dev container + infrastructure
./dev.sh

# Or use a pre-built image instead of building locally
docker compose -f dev-compose.pull.yml up -d

# Enter the container
docker exec -it online-judge-dev bash

# Install dependencies (first time only)
install-deps.sh

# Start all services
start.sh

# Stop all services
stop.sh

# View logs
log.sh <service>
# e.g. log.sh auth

# Restart a single service
stop.sh && start.sh
```

Stop everything:

```bash
docker compose -f dev-compose.yml down
```

## Configuration

- `.env` — production endpoints (Docker service names)
- `.env.dev` — dev endpoints (localhost), used by `dev-compose.yml`

## Services

| Service | Port | Language | Log |
|---|---|---|---|
| gateway | 81 | Go | `/code/logs/gateway.log` |
| submission-judge | 8000 | Go | `/code/logs/submission-judge.log` |
| contest | 8001 | Go | `/code/logs/contest.log` |
| problem | 3000 | Go | `/code/logs/problem.log` |
| auth-v2 | 50051 | Node.js | `/code/logs/auth.log` |
