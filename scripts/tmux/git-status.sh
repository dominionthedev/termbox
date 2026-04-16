#!/usr/bin/env bash
# Git status for tmux status bar — called on each status-interval tick

# Get the path of the active pane
if [ -n "$TMUX" ]; then
    current_path=$(tmux display-message -p -F "#{pane_current_path}" 2>/dev/null)
else
    current_path="$(pwd)"
fi

[ -z "$current_path" ] && exit 0

# Check if in git repo (timeout: network mounts can hang)
if ! timeout 2 git -C "$current_path" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo ""
    exit 0
fi

branch=$(timeout 2 git -C "$current_path" rev-parse --abbrev-ref HEAD 2>/dev/null)
[ -z "$branch" ] && exit 0

dirty=$(timeout 2 git -C "$current_path" status --porcelain 2>/dev/null | wc -l | tr -d ' ')

ahead=0
behind=0
ahead_behind=$(timeout 2 git -C "$current_path" rev-list --left-right --count "HEAD...@{upstream}" 2>/dev/null)
if [ -n "$ahead_behind" ]; then
    ahead=$(echo "$ahead_behind" | awk '{print $1}')
    behind=$(echo "$ahead_behind" | awk '{print $2}')
fi

output=" $branch"
[ "$dirty"  -gt 0 ] && output+=" ✗$dirty"
[ "$ahead"  -gt 0 ] && output+=" ↑$ahead"
[ "$behind" -gt 0 ] && output+=" ↓$behind"

echo "$output "
