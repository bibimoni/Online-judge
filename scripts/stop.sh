#!/bin/bash

LOG_DIR="/code/logs"

for pidfile in "$LOG_DIR"/*.pid; do
    [ -f "$pidfile" ] || continue
    name=$(basename "$pidfile" .pid)
    pid=$(cat "$pidfile")
    if kill -0 "$pid" 2>/dev/null; then
        kill -- -"$pid" 2>/dev/null
        echo "Stopped $name (pid $pid)"
    else
        echo "$name (pid $pid) already stopped"
    fi
    rm "$pidfile"
done
