#!/usr/bin/env bash
# Docker status for Starship

if ! command -v docker >/dev/null 2>&1; then
    exit 0
fi

# Get running container count
running=$(docker ps -q 2>/dev/null | wc -l | tr -d ' ')

if [ "$running" -gt 0 ]; then
    echo " $running"
fi
