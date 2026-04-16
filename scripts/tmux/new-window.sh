#!/usr/bin/env bash
# Triggered when new window is created

# Set window name based on current directory
current_dir=$(basename "$(pwd)")
tmux rename-window "$current_dir" 2>/dev/null || true
