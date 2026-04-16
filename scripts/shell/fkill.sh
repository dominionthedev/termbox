#!/usr/bin/env bash
# Interactive process killer using fzf
# Usage: fkill [signal]  (default signal: 9)

signal="${1:-9}"

pid=$(ps -ef | sed 1d |
    fzf --multi \
        --height=60% \
        --reverse \
        --border \
        --header="Select process(es) to kill with signal $signal" \
        --preview 'echo {}' \
    | awk '{print $2}')

if [ -n "$pid" ]; then
    echo "$pid" | xargs kill "-${signal}" && echo "Sent signal $signal to PID(s): $pid"
fi
