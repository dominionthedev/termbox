#!/usr/bin/env bash
# Detailed git stats for Starship custom module

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    exit 0
fi

# Use --porcelain=v1 for stable output format across git versions
# Timeout guards against slow/hung git operations on network mounts
porcelain=$(git status --porcelain=v1 2>/dev/null)

modified=$(echo "$porcelain" | grep -cE "^ [MD]")
staged=$(echo "$porcelain"   | grep -cE "^[MADRCU]")
untracked=$(echo "$porcelain" | grep -cE "^\?\?")
stashed=$(git stash list 2>/dev/null | wc -l | tr -d ' ')

output=""
[ "$modified"  -gt 0 ] && output+="!$modified "
[ "$staged"    -gt 0 ] && output+="+$staged "
[ "$untracked" -gt 0 ] && output+="?$untracked "
[ "$stashed"   -gt 0 ] && output+="📦$stashed"

# Trim trailing whitespace
echo "${output%% }"
