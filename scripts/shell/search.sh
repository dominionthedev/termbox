#!/usr/bin/env bash
# Interactive ripgrep search — opens result in nvim at the matched line
# Usage: search [pattern]

if ! command -v rg >/dev/null 2>&1; then
    echo "Error: ripgrep (rg) is required" >&2
    exit 1
fi
if ! command -v fzf >/dev/null 2>&1; then
    echo "Error: fzf is required" >&2
    exit 1
fi

rg --color=always --line-number --no-heading --smart-case "${*:-}" \
  | fzf --ansi \
        --height=80% \
        --reverse \
        --border \
        --preview "bat --color=always {1} --highlight-line {2}" \
        --preview-window 'right:60%:wrap' \
        --bind 'enter:become(nvim {1} +{2})'
