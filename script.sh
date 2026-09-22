#!/usr/bin/env bash

set -euo pipefail

INTERVAL=5

LOG_DIR="${LOG_DIR:-/var/log/script}"
LOG_FILE="${LOG_FILE:-$LOG_DIR/monitor.log}"

mkdir -p "$LOG_DIR"

cleanup() {
    echo ""
    echo "Monitoring stopped. Log file: $LOG_FILE"
    exit 0
}

trap cleanup INT TERM

if ! touch "$LOG_FILE" 2>/dev/null; then
    echo "Error: cannot write to log file '$LOG_FILE'" >&2
    exit 1
fi

echo "Starting system monitoring. Logging to: $LOG_FILE"
echo "Press Ctrl+C to stop."

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