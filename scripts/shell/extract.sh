#!/usr/bin/env bash
# Extract any archive format
# Usage: extract <file>

if [ -z "$1" ]; then
    echo "Usage: extract <file>" >&2
    exit 1
fi

if [ ! -f "$1" ]; then
    echo "Error: '$1' is not a file" >&2
    exit 1
fi

case "$1" in
    *.tar.bz2)  tar xjf "$1"    ;;
    *.tar.gz)   tar xzf "$1"    ;;
    *.tar.xz)   tar xJf "$1"    ;;
    *.tar.zst)  tar --zstd -xf "$1" ;;
    *.bz2)      bunzip2 "$1"    ;;
    *.rar)      unrar x "$1"    ;;
    *.gz)       gunzip "$1"     ;;
    *.tar)      tar xf "$1"     ;;
    *.tbz2)     tar xjf "$1"    ;;
    *.tgz)      tar xzf "$1"    ;;
    *.zip)      unzip "$1"      ;;
    *.Z)        uncompress "$1" ;;
    *.7z)       7z x "$1"       ;;
    *)
        echo "Unknown archive format: $1" >&2
        exit 1
        ;;
esac

echo "Extracted: $1"
