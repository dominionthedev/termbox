#!/usr/bin/env bash
# cpods — list running containers (podman-first, docker fallback)
CT=$(command -v podman &>/dev/null && echo podman || echo docker)
echo "  Running containers ($CT):"
$CT ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null
