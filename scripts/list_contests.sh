#!/bin/bash

# List Contests Script
# Usage: ./list_contests.sh
# 
# This script lists all contests from MongoDB via the contest-mongo container.

# Configuration
CONTAINER_NAME="contest-mongo"
DATABASE_NAME="contestdb"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${RED}Error: Container '${CONTAINER_NAME}' is not running${NC}"
    echo "Please start it with: docker compose up -d contestmongodb"
    exit 1
fi

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}        Available Contests${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""

# Query contests and display them
docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "
    var contests = db.Contest.find({}, {
        _id: 1,
        name: 1,
        status: 1,
        'contest_rule.scoring_type': 1,
        start_time: 1,
        end_time: 1
    }).toArray();
    
    if (contests.length === 0) {
        print('No contests found in the database.');
        print('');
        print('You may need to create a contest first through the API.');
    } else {
        print('ID\t\t\t\t\tName\t\t\tType\tStatus');
        print('--\t\t\t\t\t----\t\t\t----\t------');
        contests.forEach(function(c) {
            var scoringType = c.contest_rule ? c.contest_rule.scoring_type : 'N/A';
            print(c._id + '\t' + (c.name || 'Unnamed') + '\t\t' + scoringType + '\t' + (c.status || 'N/A'));
        });
    }
"

echo ""
echo -e "${YELLOW}To feed submissions into a contest, copy the Contest ID and run:${NC}"
echo -e "${GREEN}./scripts/feed_contest_submissions.sh <contest_id>${NC}"
