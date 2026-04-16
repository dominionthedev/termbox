#!/usr/bin/env bash
# cclean — prune stopped containers and dangling images
CT=$(command -v podman &>/dev/null && echo podman || echo docker)
echo "  Pruning stopped containers ($CT)..."
$CT container prune -f
echo "  Pruning dangling images..."
$CT image prune -f
echo "  Done."
