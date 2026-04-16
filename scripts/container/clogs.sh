#!/usr/bin/env bash
# clogs — tail container logs interactively
CT=$(command -v podman &>/dev/null && echo podman || echo docker)
container=$($CT ps --format '{{.Names}}' 2>/dev/null |
    fzf --height=40% --reverse --border --header="Tail logs ($CT)")
[[ -z "$container" ]] && exit 0
$CT logs -f --tail=100 "$container"
