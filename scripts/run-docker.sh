#!/bin/bash
set -e

docker compose -f dev-compose.yml down --remove-orphans 2>/dev/null || true
docker compose -f dev-compose.yml up -d

