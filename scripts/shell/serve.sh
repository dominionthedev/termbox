#!/usr/bin/env bash
# Serve current directory over HTTP
# Usage: serve [port]

port="${1:-8000}"
echo "Serving $(pwd) at http://localhost:$port"
echo "Press Ctrl+C to stop."
python3 -m http.server "$port"
