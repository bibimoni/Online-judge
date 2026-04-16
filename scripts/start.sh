#!/bin/bash
set -e

export PATH="/usr/local/go/bin:/root/go/bin:$PATH"

ENV_FILE="/code/.env.dev"

if [ ! -f "$ENV_FILE" ]; then
    echo "ERROR: $ENV_FILE not found"
    exit 1
fi

set -a; source "$ENV_FILE"; set +a

LOG_DIR="/code/logs"
mkdir -p "$LOG_DIR"

check_port() {
    if nc -z 127.0.0.1 "$1" 2>/dev/null; then
        echo "ERROR: Port $1 ($2) is already in use"
        return 1
    fi
}

check_port "$GATEWAY_PORT"       "gateway"          || exit 1
check_port "$SUBMISSION_PORT"   "submission-judge" || exit 1
check_port "$CONTEST_PORT"      "contest"          || exit 1
check_port "$PROBLEM_PORT"      "problem"          || exit 1
check_port "$AUTH_PORT"         "auth"             || exit 1

wait_for() {
    local host="$1" port="$2" name="$3" max=30 i=0
    while ! nc -z "$host" "$port" 2>/dev/null; do
        i=$((i + 1))
        if [ "$i" -ge "$max" ]; then
            echo "ERROR: $name not available at $host:$port" && exit 1
        fi
        sleep 1
    done
    echo "$name is ready"
}

wait_for mongosubmissionjudgedb 27017 "Submission MongoDB"
wait_for problemmongodb 27017 "Problem MongoDB"
wait_for contestmongodb 27017 "Contest MongoDB"
wait_for redissubmissionjudge 6379 "Submission Redis"
wait_for scoreboard_redis 6379 "Scoreboard Redis"
wait_for mypostgres 5432 "PostgreSQL"

start_service() {
    local name="$1"
    shift
    setsid "$@" > "$LOG_DIR/${name}.log" 2>&1 &
    local pid=$!
    echo "$pid" > "$LOG_DIR/${name}.pid"
    echo "$name started (pid $pid) — log: $LOG_DIR/${name}.log"
}

start_service auth         /code/scripts/start-auth.sh
start_service gateway      /code/scripts/start-gateway.sh
start_service problem      /code/scripts/start-problem.sh
start_service contest      /code/scripts/start-contest.sh
start_service submission-judge /code/scripts/start-submission-judge.sh

echo ""
echo "All services started. Use 'tail -f $LOG_DIR/<service>.log' to view logs."
