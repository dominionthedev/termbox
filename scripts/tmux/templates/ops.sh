#!/usr/bin/env bash
# Ops session template - for operations/monitoring
# Usage: ops.sh <session-name>

session="$1"

# Get the first window of the session
first_window=$(tmux list-windows -t "$session" -F "#{window_id}" | head -n 1)

# Layout:
#  ┌──────────┬──────────┐
#  │  Monitor │   Logs   │
#  ├──────────┼──────────┤
#  │  Docker  │   Kubectl│
#  └──────────┴──────────┘

tmux rename-window -t "$first_window" "ops"

# Create 2x2 grid
tmux split-window -h -t "$first_window"
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))
tmux split-window -v -t "${panes[0]}"
tmux split-window -v -t "${panes[1]}"

# Re-get all 4 panes
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))

# Monitor pane (top-left)
tmux send-keys -t "${panes[0]}" "# System monitor" C-m
if command -v htop >/dev/null 2>&1; then
    tmux send-keys -t "${panes[0]}" "htop" C-m
else
    tmux send-keys -t "${panes[0]}" "top" C-m
fi

# Logs pane (top-right)
tmux send-keys -t "${panes[1]}" "# Application logs" C-m

# Docker pane (bottom-left)
tmux send-keys -t "${panes[2]}" "# Docker management" C-m
if command -v docker >/dev/null 2>&1; then
    tmux send-keys -t "${panes[2]}" "docker ps" C-m
fi

# Kubectl pane (bottom-right)
tmux send-keys -t "${panes[3]}" "# Kubernetes management" C-m
if command -v kubectl >/dev/null 2>&1; then
    tmux send-keys -t "${panes[3]}" "kubectl get pods" C-m
fi

# Create additional windows
tmux new-window -t "$session" -n "ssh"
tmux new-window -t "$session" -n "scripts"

# Select first window, monitor pane
tmux select-window -t "$first_window"
tmux select-pane -t "${panes[0]}"
