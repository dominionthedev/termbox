#!/usr/bin/env bash
# Dev session template - for development work
# Usage: dev.sh <session-name>

session="$1"

# Get the first window of the session
first_window=$(tmux list-windows -t "$session" -F "#{window_id}" | head -n 1)

# Layout:
#  ┌─────────────────┬──────────┐
#  │                 │          │
#  │     Editor      │   Git    │
#  │                 │          │
#  ├─────────────────┼──────────┤
#  │    Terminal     │  Logs    │
#  └─────────────────┴──────────┘

# Rename first window
tmux rename-window -t "$first_window" "dev"

# Split right for git (pane 2)
tmux split-window -h -p 30 -t "$first_window"

# Get pane IDs for the window
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))

# Split bottom left for terminal (from the first pane)
tmux split-window -v -p 30 -t "${panes[0]}"

# Re-get panes
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))

# Split bottom right for logs (from the second pane)
tmux split-window -v -p 30 -t "${panes[1]}"

# Re-get all 4 panes
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))

# Set up panes
tmux send-keys -t "${panes[0]}" "# Editor pane - nvim ." C-m
tmux send-keys -t "${panes[1]}" "git status" C-m
tmux send-keys -t "${panes[2]}" "clear" C-m
tmux send-keys -t "${panes[3]}" "# Logs pane" C-m

# Select editor pane
tmux select-pane -t "${panes[0]}"

# Create additional windows
tmux new-window -t "$session" -n "tests"
tmux new-window -t "$session" -n "server"

# Select first window
tmux select-window -t "$first_window"
