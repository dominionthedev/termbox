#!/usr/bin/env bash
# Triggered when pane gets focus

# Find project root (where this script lives)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Update git status using the dynamic path
"$SCRIPT_DIR/git-status.sh" > /dev/null

# Could trigger other updates here
# - Project detection
# - Environment activation
# - Custom notifications
