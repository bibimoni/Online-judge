#!/bin/bash

# Clear Contest Submissions Script
# Usage: ./clear_contest_submissions.sh <contest_id>
# 
# This script removes all submissions for a specific contest from MongoDB.

set -e

# Configuration
CONTAINER_NAME="contest-mongo"
DATABASE_NAME="contestdb"
COLLECTION_NAME="ContestSubmission"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if contest_id is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Contest ID is required${NC}"
    echo "Usage: $0 <contest_id>"
    echo "Example: $0 683c5b4e1234567890abcdef"
    exit 1
fi

CONTEST_ID="$1"

# Validate contest_id format (24 character hex string)
if ! [[ "$CONTEST_ID" =~ ^[0-9a-fA-F]{24}$ ]]; then
    echo -e "${RED}Error: Invalid contest ID format. Must be a 24-character hex string (MongoDB ObjectId)${NC}"
    echo "Example: 683c5b4e1234567890abcdef"
    exit 1
fi

# Check if Docker container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${RED}Error: Container '${CONTAINER_NAME}' is not running${NC}"
    echo "Please start it with: docker compose up -d contestmongodb"
    exit 1
fi

echo -e "${YELLOW}Clearing submissions for contest: ${CONTEST_ID}${NC}"
echo ""

# Count existing submissions
COUNT=$(docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "db.${COLLECTION_NAME}.countDocuments({contest_id: ObjectId('${CONTEST_ID}')})")

if [ "$COUNT" = "0" ]; then
    echo -e "${YELLOW}No submissions found for this contest.${NC}"
    exit 0
fi

echo -e "Found ${YELLOW}${COUNT}${NC} submissions to delete."
read -p "Are you sure you want to delete them? (y/N): " confirm

if [[ "$confirm" =~ ^[Yy]$ ]]; then
    docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "db.${COLLECTION_NAME}.deleteMany({contest_id: ObjectId('${CONTEST_ID}')})"
    echo -e "${GREEN}✓ Successfully deleted ${COUNT} submissions${NC}"
else
    echo -e "${YELLOW}Operation cancelled.${NC}"
fi
