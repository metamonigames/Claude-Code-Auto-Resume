#!/bin/bash
# claude_auto_ping.sh

# Get current hour (0-23)
HOUR=$(date +%H)

# Define time-based logic
if [ "$HOUR" -ge 2 ] && [ "$HOUR" -lt 8 ]; then
    # Night Shift (02:00 - 07:59): Sonnet mode, wait for response
    echo "[$(date)] Night Shift: Using Sonnet for task progression."
    export CLAUDE_CODE_MODEL=claude-sonnet-4-6
    claude "Proceed with the next task. Commit after each task is finished."
else
    # Day Shift: Haiku mode, fire-and-forget (no need to wait for response)
    echo "[$(date)] Day Shift: Using Haiku for session keep-alive."
    export CLAUDE_CODE_MODEL=claude-haiku-4-5-20251001
    nohup claude -p "ping" > /dev/null 2>&1 &
    echo "[$(date)] Ping sent. Exiting without waiting for response."
fi