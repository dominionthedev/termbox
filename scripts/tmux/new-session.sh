#!/usr/bin/env bash
# Triggered when new session is created
session_name="$1"

# Display welcome message
tmux display-message "🔥 Session '$session_name' created! Press Ctrl+Space Space for dashboard"
