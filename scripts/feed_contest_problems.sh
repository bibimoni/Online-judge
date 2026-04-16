#!/bin/bash

# Feed Contest Problems Script
# Usage: ./feed_contest_problems.sh <contest_id>
# 
# This script adds sample problems to a contest in MongoDB via the contest-mongo container.
# The contest_id should be a valid MongoDB ObjectId.

set -e

# Configuration
CONTAINER_NAME="contest-mongo"
DATABASE_NAME="contestdb"
COLLECTION_NAME="Contest"

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

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}   Feed Contest Problems Script${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "Container: ${YELLOW}${CONTAINER_NAME}${NC}"
echo -e "Database:  ${YELLOW}${DATABASE_NAME}${NC}"
echo -e "Collection: ${YELLOW}${COLLECTION_NAME}${NC}"
echo -e "Contest ID: ${YELLOW}${CONTEST_ID}${NC}"
echo ""

# Check if contest exists
CONTEST_EXISTS=$(docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "db.${COLLECTION_NAME}.countDocuments({_id: ObjectId('${CONTEST_ID}')})")

if [ "$CONTEST_EXISTS" = "0" ]; then
    echo -e "${RED}Error: Contest with ID '${CONTEST_ID}' not found${NC}"
    echo "Use ./scripts/list_contests.sh to see available contests"
    exit 1
fi

# Define sample problems
# Format: problem_id|label|max_points
PROBLEMS=(
    "1|A|100"
    "2|B|100"
    "3|C|100"
    "4|D|200"
    "5|E|300"
)

echo -e "${YELLOW}Adding ${#PROBLEMS[@]} problems to contest...${NC}"
echo ""

# Build the problems array
PROBLEMS_ARRAY="["
FIRST=true

for problem in "${PROBLEMS[@]}"; do
    IFS='|' read -r problem_id label max_points <<< "$problem"
    
    if [ "$FIRST" = true ]; then
        FIRST=false
    else
        PROBLEMS_ARRAY+=","
    fi
    
    PROBLEMS_ARRAY+="{
        \"problem_id\": NumberLong(${problem_id}),
        \"label\": \"${label}\",
        \"max_points\": ${max_points}
    }"
    
    echo -e "  Problem ${YELLOW}${label}${NC} (ID: ${problem_id}, Max Points: ${max_points})"
done

PROBLEMS_ARRAY+="]"

echo ""

# Create the MongoDB update command
MONGO_COMMAND="db.${COLLECTION_NAME}.updateOne(
    { _id: ObjectId('${CONTEST_ID}') },
    { \$set: { problems: ${PROBLEMS_ARRAY} } }
)"

# Execute the command in the container
echo -e "${YELLOW}Executing MongoDB update command...${NC}"
echo ""

docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "${MONGO_COMMAND}"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Successfully added ${#PROBLEMS[@]} problems to contest${NC}"
    echo ""
    
    # Show updated contest
    echo -e "${YELLOW}Updated contest problems:${NC}"
    docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "
        var contest = db.${COLLECTION_NAME}.findOne({ _id: ObjectId('${CONTEST_ID}') }, { problems: 1, name: 1 });
        print('Contest: ' + (contest.name || 'Unnamed'));
        print('');
        print('Label\tProblem ID\tMax Points');
        print('-----\t----------\t----------');
        if (contest.problems) {
            contest.problems.forEach(function(p) {
                print(p.label + '\t' + p.problem_id + '\t\t' + p.max_points);
            });
        }
    "
else
    echo -e "${RED}✗ Failed to update contest problems${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Done!${NC}"
echo ""
echo -e "To view the contest, run:"
echo -e "${YELLOW}docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --eval \"db.${COLLECTION_NAME}.findOne({_id: ObjectId('${CONTEST_ID}')})\"${NC}"
