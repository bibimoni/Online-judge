---
description: How to add a new Judge Node on a separate machine
---

# Adding a Horizontal Judge Node

This workflow describes how to scale your judging capacity by adding a new `submission-judge` instance on a **different computer**.

## Prerequisites

1.  **Docker & Docker Compose** installed on the new machine.
2.  **Network Access**: The new machine must be able to reach:
    *   **Redis**: The main Redis instance (e.g., `redis-submission-judge` on the main server).
    *   **Problem Service**: The Problem Service API (e.g., `problem` on the main server).
    *   **MongoDB**: (Optional, if the judge writes directly to DB, otherwise it just needs Redis). *Note: Current implementation writes to DB.*

## Step 1: Prepare Environment Variables

Create a `.env` file on the new machine. You need to point the services to the **IP Address** of your main server (where Redis and Problem Service are running).

```bash
# Main Server IP (Replace with actual IP)
MAIN_SERVER_IP=192.168.1.100

# Service Configuration
SUBMISSION_HOST=0.0.0.0
SUBMISSION_PORT=8000
SUBMISSION_LOG_LEVEL=debug

# Redis Connection (Critical for Job Queue)
SUBMISSION_REDIS_URI=${MAIN_SERVER_IP}:6379
SUBMISSION_REDIS_PASSWORD=root

# MongoDB Connection (For writing results)
SUBMISSION_MONGODB_URI=mongodb://${MAIN_SERVER_IP}:27017/submissionjudgedb
SUBMISSION_MONGODB_DATABASE_NAME=submissionjudgedb

# Problem Service (For downloading test cases)
# Note: Ensure this port is exposed on the main server
PROBLEM_HOST=${MAIN_SERVER_IP}
PROBLEM_PORT=3000

# Judge Configuration
SUBMISSION_NUMBER_OF_JUDGE=5
SUBMISSION_JUDGE_ID_OFFSET=0 
# Tip: If you run multiple nodes, you might want to adjust ID_OFFSET to avoid ID collisions if logs/metrics use it, 
# though for stateless workers it matters less.

# Paths
# The judge will use this local cache for downloaded test cases
SUBMISSION_JUDGE_PROBLEM_DIR=/app/cache/problems
```

## Step 2: Create Docker Compose

Create a `docker-compose.yml` for the judge node:

```yaml
services:
  submission-judge-node:
    image: your-docker-registry/submission-judge:latest # Or build locally
    container_name: submission-judge-node
    restart: always
    env_file:
      - .env
    ports:
      - "8000:8000"
    volumes:
      # Local cache for test cases (Hybrid Approach)
      - ./judge-cache:/app/cache/problems
    privileged: true # Required for Isolate sandbox
```

## Step 3: Run the Judge

1.  Start the container:
    ```bash
    docker-compose up -d
    ```

2.  **Verify**:
    *   Check logs: `docker logs -f submission-judge-node`
    *   You should see it connecting to Redis and waiting for jobs.
    *   When a submission comes in, if it doesn't have the test case locally, it will log that it's downloading from `http://${MAIN_SERVER_IP}:3000/...`.

## Troubleshooting

*   **Connection Refused**: Ensure ports `6379` (Redis), `27017` (Mongo), and `3000` (Problem) are exposed on the Main Server and not blocked by firewall.
*   **Download Failures**: Check if `PROBLEM_HOST` and `PROBLEM_PORT` are correct and reachable from inside the container.
