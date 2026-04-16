#!/usr/bin/env bash
# dex — exec into a running container (podman-first, docker fallback)
# Usage: dex [shell]
CT=$(command -v podman &>/dev/null && echo podman || echo docker)
shell="${1:-sh}"
container=$($CT ps --format '{{.Names}}\t{{.Image}}\t{{.Status}}' 2>/dev/null |
    fzf --height=50% --reverse --border \
        --header="Container exec ($CT) — shell: $shell" \
        --preview "$CT logs --tail 30 {1}" \
        --preview-window 'right:50%' \
    | awk '{print $1}')
[[ -z "$container" ]] && exit 0
echo "  → $container ($CT)"
$CT exec -it "$container" "$shell"
