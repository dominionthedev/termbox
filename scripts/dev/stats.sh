#!/usr/bin/env bash
# Code statistics for the current or specified directory
# Usage: stats [directory]

dir="${1:-.}"

if [ ! -d "$dir" ]; then
    echo "Error: '$dir' is not a directory" >&2
    exit 1
fi

echo "Code statistics: $dir"
echo ""

if command -v tokei >/dev/null 2>&1; then
    tokei "$dir"
elif command -v cloc >/dev/null 2>&1; then
    cloc "$dir"
else
    # Pure shell fallback: file counts per extension
    echo "Install tokei (cargo install tokei) for detailed stats."
    echo ""
    find "$dir" -type f -not -path '*/.git/*' |
        sed 's/.*\.//' | sort | uniq -c | sort -rn |
        head -20 |
        awk '{printf "  %-20s %s files\n", $2, $1}'
fi
