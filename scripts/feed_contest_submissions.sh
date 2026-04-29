#!/bin/bash

# Feed Contest Submissions Script
# Usage: ./feed_contest_submissions.sh <contest_id>
# 
# This script feeds sample contest submissions into MongoDB via the contest-mongo container.
# The contest_id should be a valid MongoDB ObjectId.

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

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}  Feed Contest Submissions Script${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "Container: ${YELLOW}${CONTAINER_NAME}${NC}"
echo -e "Database:  ${YELLOW}${DATABASE_NAME}${NC}"
echo -e "Collection: ${YELLOW}${COLLECTION_NAME}${NC}"
echo -e "Contest ID: ${YELLOW}${CONTEST_ID}${NC}"
echo ""

# Define sample submissions
# Format: username|problem_id|verdict|points|eval_status|submission_type|minutes_offset
SUBMISSIONS=(
    "alice|1|ACCEPTED|100|FINISHED|RATED|5"
    "alice|2|WRONG_ANSWER|0|FINISHED|RATED|10"
    "alice|2|WRONG_ANSWER|0|FINISHED|RATED|15"
    "alice|2|ACCEPTED|100|FINISHED|RATED|20"
    "alice|3|PARTIAL_RESULT|50|FINISHED|RATED|25"
    "bob|1|WRONG_ANSWER|0|FINISHED|RATED|6"
    "bob|1|ACCEPTED|100|FINISHED|RATED|12"
    "bob|2|ACCEPTED|100|FINISHED|RATED|18"
    "bob|3|TIME_LIMIT_EXCEEDED|0|FINISHED|RATED|22"
    "charlie|1|COMPILATION_ERROR|0|FINISHED|RATED|3"
    "charlie|1|RUNTIME_ERROR|0|FINISHED|RATED|8"
    "charlie|1|ACCEPTED|100|FINISHED|RATED|15"
    "charlie|2|PARTIAL_RESULT|75|FINISHED|RATED|30"
    "charlie|3|ACCEPTED|100|FINISHED|RATED|45"
    "david|1|PENDING|0|PENDING|RATED|2"
    "david|2|ACCEPTED|100|FINISHED|RATED|40"
    "eve|1|ACCEPTED|100|FINISHED|VIRTUAL|10"
    "eve|2|WRONG_ANSWER|0|FINISHED|VIRTUAL|15"
    "eve|2|ACCEPTED|100|FINISHED|VIRTUAL|25"
)

# Function to generate a random ObjectId
generate_objectid() {
    # Generate 24 random hex characters
    cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 24 | head -n 1
}

# Function to get current timestamp in ISO format
get_iso_timestamp() {
    local minutes_offset=$1
    # Get current time and add offset in minutes
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        date -u -v+${minutes_offset}M +"%Y-%m-%dT%H:%M:%SZ"
    else
        # Linux
        date -u -d "+${minutes_offset} minutes" +"%Y-%m-%dT%H:%M:%SZ"
    fi
}

echo -e "${YELLOW}Inserting ${#SUBMISSIONS[@]} sample submissions...${NC}"
echo ""

# Build the insertMany command
DOCUMENTS="["
FIRST=true

for submission in "${SUBMISSIONS[@]}"; do
    IFS='|' read -r username problem_id verdict points eval_status submission_type minutes_offset <<< "$submission"
    
    submission_id=$(generate_objectid)
    submit_at=$(get_iso_timestamp $minutes_offset)
    updated_at=$(get_iso_timestamp $minutes_offset)
    
    # Determine ignored status (false for all in this sample)
    ignored="false"
    
    if [ "$FIRST" = true ]; then
        FIRST=false
    else
        DOCUMENTS+=","
    fi
    
    DOCUMENTS+="{
        \"contest_id\": ObjectId(\"${CONTEST_ID}\"),
        \"submission_id\": \"${submission_id}\",
        \"username\": \"${username}\",
        \"problem_id\": NumberLong(${problem_id}),
        \"submit_at\": ISODate(\"${submit_at}\"),
        \"verdict\": \"${verdict}\",
        \"points\": ${points},
        \"eval_status\": \"${eval_status}\",
        \"ignored\": ${ignored},
        \"updated_at\": ISODate(\"${updated_at}\"),
        \"submission_type\": \"${submission_type}\"
    }"
done

DOCUMENTS+="]"

# Create the MongoDB command
MONGO_COMMAND="db.${COLLECTION_NAME}.insertMany(${DOCUMENTS})"

# Execute the command in the container
echo -e "${YELLOW}Executing MongoDB insertMany command...${NC}"
echo ""

docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "${MONGO_COMMAND}"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Successfully inserted ${#SUBMISSIONS[@]} submissions${NC}"
    echo ""
    
    # Show summary
    echo -e "${YELLOW}Summary of inserted submissions:${NC}"
    docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --quiet --eval "
        var result = db.${COLLECTION_NAME}.aggregate([
            { \$match: { contest_id: ObjectId('${CONTEST_ID}') } },
            { \$group: { 
                _id: '\$username', 
                total_submissions: { \$sum: 1 },
                accepted: { \$sum: { \$cond: [{ \$eq: ['\$verdict', 'ACCEPTED'] }, 1, 0] } },
                total_points: { \$max: '\$points' }
            }},
            { \$sort: { accepted: -1, total_points: -1 } }
        ]).toArray();
        print('Username\t\tSubmissions\tAccepted\tMax Points');
        print('--------\t\t-----------\t--------\t----------');
        result.forEach(function(r) {
            print(r._id + '\t\t\t' + r.total_submissions + '\t\t' + r.accepted + '\t\t' + r.total_points);
        });
    "
else
    echo -e "${RED}✗ Failed to insert submissions${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Done!${NC}"
echo ""
echo -e "To view all submissions for this contest, run:"
echo -e "${YELLOW}docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --eval \"db.${COLLECTION_NAME}.find({contest_id: ObjectId('${CONTEST_ID}')}).pretty()\"${NC}"
echo ""
echo -e "To delete all submissions for this contest, run:"
echo -e "${YELLOW}docker exec -i ${CONTAINER_NAME} mongosh ${DATABASE_NAME} --eval \"db.${COLLECTION_NAME}.deleteMany({contest_id: ObjectId('${CONTEST_ID}')})\"${NC}"
