
INTERVAL=5


SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"


LOG_FILE="$SCRIPT_DIR/monitor.log"

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