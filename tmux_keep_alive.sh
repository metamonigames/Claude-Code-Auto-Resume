#!/bin/bash
TARGET_TEXT="You've hit your limit"

tmux list-panes -a -F '#{session_name}:#{window_index}.#{pane_index}' | while read -r pane_id; do
    content=$(tmux capture-pane -pt "$pane_id" -S -15)
    
    if echo "$content" | grep -iq "$TARGET_TEXT"; then
        tmux send-keys -t "$pane_id" Enter Enter
        echo "[$(date)] Sent 'Double Enter' to $pane_id"
    fi
done