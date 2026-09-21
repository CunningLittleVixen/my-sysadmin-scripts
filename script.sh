#!/usr/bin/env bash

set -euo pipefail

INTERVAL=5


LOG_FILE="monitor.log"

cleanup() {
    echo ""
    echo "Monitoring stopped. Log file: $LOG_FILE"
    exit 0
}

trap cleanup INT

if ! touch "$LOG_FILE" 2>/dev/null; then
    echo "Error: cannot write to log file '$LOG_FILE'" >&2
    exit 1
fi

while true; do
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

    {
        echo "--- $TIMESTAMP ---"
        echo "=== free -h ==="
        free -h
        echo "=== df -h ==="
        df -h
        echo "=== uptime ==="
        uptime
        echo ""
    } >> "$LOG_FILE"

    sleep "$INTERVAL"
done