#!/usr/bin/env bash
# Default session template
# Usage: default.sh <session-name>

session="$1"

# Get the first window of the session
first_window=$(tmux list-windows -t "$session" -F "#{window_id}" | head -n 1)

# Create a 3-pane layout
# Layout: 
#  ┌────────┬────────┐
#  │        │        │
#  │  Main  │  Side  │
#  │        │        │
#  ├────────┴────────┤
#  │    Bottom       │
#  └─────────────────┘

# Split bottom half
tmux split-window -v -t "$first_window"
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))

# Split top half horizontally
tmux split-window -h -t "${panes[0]}"

# Re-get all 3 panes
panes=($(tmux list-panes -t "$first_window" -F "#{pane_id}"))

# Select main pane
tmux select-pane -t "${panes[0]}"
