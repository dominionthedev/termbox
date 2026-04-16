#!/usr/bin/env bash
# Conventional commit helper — interactive type + scope + message
# Usage: gci (no args — fully interactive)

if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "Error: not inside a git repository" >&2
    exit 1
fi

# Pick commit type
commit_type=$(printf 'feat\nfix\ndocs\nstyle\nrefactor\nperf\ntest\nchore\nbuild\nci\nrevert' |
    fzf --prompt="Commit type: " \
        --height=14 \
        --reverse \
        --border \
        --no-multi)

[ -z "$commit_type" ] && exit 0

# Optional scope
printf "Scope (optional, press Enter to skip): "
read -r scope

# Required message
while true; do
    printf "Message: "
    read -r message
    [ -n "$message" ] && break
    echo "Message cannot be empty."
done

# Optional breaking change
printf "Breaking change? [y/N]: "
read -r breaking
breaking_marker=""
[[ "$breaking" =~ ^[Yy]$ ]] && breaking_marker="!"

# Build commit string
if [ -n "$scope" ]; then
    commit_msg="${commit_type}(${scope})${breaking_marker}: ${message}"
else
    commit_msg="${commit_type}${breaking_marker}: ${message}"
fi

echo ""
echo "  → $commit_msg"
printf "Confirm commit? [Y/n]: "
read -r confirm
[[ "$confirm" =~ ^[Nn]$ ]] && echo "Aborted." && exit 0

git commit -m "$commit_msg"
